package service

import (
	"PaperExamGrader/internal/model"
	"PaperExamGrader/internal/repository"
	"PaperExamGrader/internal/storage"
	"PaperExamGrader/internal/transport/request"
	"PaperExamGrader/pkg/logging"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gen2brain/go-fitz" // PDF to image
	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"gocv.io/x/gocv" // OpenCV wrapper
	"image"
	"io"
	"mime/multipart"
	"net/http"
)

type BBoxMeta struct {
	Page        int
	BBoxPercent [4]float64
}

type ManualCropper struct {
	bboxMap     []BBoxMeta
	fileStorage storage.FileStorage
	logger      *logrus.Logger
	imgRepo     *repository.ImageRepository
	answerRepo  *repository.AnswerRepository
}

func NewManualCropper(
	fs storage.FileStorage,
	imgRepo *repository.ImageRepository,
	answerRepo *repository.AnswerRepository,
) *ManualCropper {
	return &ManualCropper{
		fileStorage: fs,
		logger:      logging.GetLogger(),
		imgRepo:     imgRepo,
		answerRepo:  answerRepo,
	}
}

func (m *ManualCropper) SetCropParams(bboxes []BBoxMeta) {
	m.bboxMap = bboxes
}

func (m *ManualCropper) CropFromExamPDFURLs(examID uint) ([]string, error) {
	var uploadedURLs []string

	// Step 1: Fetch all answers for the exam
	answers, err := m.answerRepo.GetByExamID(examID)
	if err != nil {
		m.logger.WithError(err).WithField("examID", examID).Error("Failed to fetch answers")
		return nil, fmt.Errorf("failed to get answers for exam ID %d", examID)
	}

	totalSteps := len(answers) * len(m.bboxMap)
	bar := progressbar.NewOptions(totalSteps,
		progressbar.OptionSetDescription("📄 Cropping PDFs"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "#",
			SaucerPadding: "-",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for _, answer := range answers {
		resp, err := http.Get(answer.PdfURL)
		if err != nil {
			m.logger.WithError(err).WithField("url", answer.PdfURL).Error("Failed to download PDF")
			continue
		}
		if resp.StatusCode != http.StatusOK {
			m.logger.WithField("status", resp.StatusCode).WithField("url", answer.PdfURL).Error("Non-200 response downloading PDF")
			resp.Body.Close()
			continue
		}

		pdfData := new(bytes.Buffer)
		_, err = io.Copy(pdfData, resp.Body)
		resp.Body.Close()
		if err != nil {
			m.logger.WithError(err).WithField("url", answer.PdfURL).Error("Failed to read PDF from response")
			continue
		}

		doc, err := fitz.NewFromMemory(pdfData.Bytes())
		if err != nil {
			m.logger.WithError(err).WithField("url", answer.PdfURL).Error("Failed to open PDF")
			continue
		}
		defer doc.Close()

		baseName := fmt.Sprintf("answer_%d", answer.ID)

		for idx, bm := range m.bboxMap {
			if bm.Page >= doc.NumPage() {
				m.logger.WithFields(logrus.Fields{
					"url":  answer.PdfURL,
					"page": bm.Page,
				}).Warn("Skipping out-of-range page")
				bar.Add(1)
				continue
			}

			img, err := doc.Image(bm.Page)
			if err != nil {
				m.logger.WithError(err).WithField("page", bm.Page).Error("Failed to render page to image")
				bar.Add(1)
				continue
			}

			mat, err := gocv.ImageToMatRGB(img)
			if err != nil {
				m.logger.WithError(err).Error("Failed to convert image to Mat")
				bar.Add(1)
				continue
			}
			defer mat.Close()

			cropped := cropFromPercentBBox(mat, bm.BBoxPercent)
			defer cropped.Close()

			imgBuf, err := gocv.IMEncode(".jpg", cropped)
			if err != nil {
				m.logger.WithError(err).Error("Failed to encode cropped image")
				bar.Add(1)
				continue
			}

			filename := fmt.Sprintf("%s_crop_%d.jpg", baseName, idx+1)
			url, err := m.uploadBytesAsFile(imgBuf.GetBytes(), filename)
			if err != nil {
				m.logger.WithError(err).Error("Failed to upload cropped image")
				bar.Add(1)
				continue
			}

			imgModel := model.Image{
				AnswerID: answer.ID,
				URL:      url,
			}
			if err := m.imgRepo.Create(&imgModel); err != nil {
				m.logger.WithError(err).Error("Failed to save image to DB")
				bar.Add(1)
				continue
			}

			m.logger.WithFields(logrus.Fields{
				"url":      url,
				"page":     bm.Page,
				"answerID": answer.ID,
			}).Info("✅ Cropped and saved image")

			uploadedURLs = append(uploadedURLs, url)
			bar.Add(1)
		}
	}

	return uploadedURLs, nil
}

// CropFromMultipartPDFs crops answers from uploaded PDFs and uploads them to file storage
func (m *ManualCropper) CropFromMultipartPDFs(files []*multipart.FileHeader, answerID uint) ([]string, error) {
	var uploadedURLs []string

	// ✅ Step 1: Check if Answer exists
	answer, err := m.answerRepo.GetByID(answerID)
	if err != nil {
		m.logger.WithError(err).WithField("answerID", answerID).Error("Answer not found")
		return nil, fmt.Errorf("answer with ID %d not found", answerID)
	}

	for _, fileHeader := range files {
		pdfFile, err := fileHeader.Open()
		if err != nil {
			m.logger.WithError(err).Error("Failed to open uploaded PDF")
			continue
		}

		pdfData := new(bytes.Buffer)
		_, err = pdfData.ReadFrom(pdfFile)
		pdfFile.Close()
		if err != nil {
			m.logger.WithError(err).Error("Failed to read uploaded PDF")
			continue
		}

		doc, err := fitz.NewFromMemory(pdfData.Bytes())
		if err != nil {
			m.logger.WithError(err).Error("Failed to decode PDF")
			continue
		}
		defer doc.Close()

		baseName := fileHeader.Filename[:len(fileHeader.Filename)-len(".pdf")]

		for idx, bm := range m.bboxMap {
			if bm.Page >= doc.NumPage() {
				m.logger.WithFields(logrus.Fields{
					"file": fileHeader.Filename,
					"page": bm.Page,
				}).Warn("Skipping out-of-range page")
				continue
			}

			img, err := doc.Image(bm.Page)
			if err != nil {
				m.logger.WithError(err).WithField("page", bm.Page).Error("Failed to render page to image")
				continue
			}

			mat, err := gocv.ImageToMatRGB(img)
			if err != nil {
				m.logger.WithError(err).Error("Failed to convert image to Mat")
				continue
			}
			defer mat.Close()

			cropped := cropFromPercentBBox(mat, bm.BBoxPercent)
			defer cropped.Close()

			imgBuf, err := gocv.IMEncode(".jpg", cropped)
			if err != nil {
				m.logger.WithError(err).Error("Failed to encode cropped image")
				continue
			}

			filename := fmt.Sprintf("%s_crop_%d.jpg", baseName, idx+1)
			url, err := m.uploadBytesAsFile(imgBuf.GetBytes(), filename)
			if err != nil {
				m.logger.WithError(err).Error("Failed to upload cropped image")
				continue
			}

			imgModel := model.Image{
				AnswerID: answer.ID,
				URL:      url,
			}
			if err := m.imgRepo.Create(&imgModel); err != nil {
				m.logger.WithError(err).Error("Failed to save image record to DB")
				continue
			}

			uploadedURLs = append(uploadedURLs, url)

			m.logger.WithFields(logrus.Fields{
				"url":      url,
				"page":     bm.Page,
				"answerID": answerID,
			}).Info("✅ Cropped and saved image")
		}
	}

	return uploadedURLs, nil
}

func cropFromPercentBBox(src gocv.Mat, bbox [4]float64) gocv.Mat {
	h, w := src.Rows(), src.Cols()
	x0 := int(bbox[0] * float64(w))
	y0 := int(bbox[1] * float64(h))
	x1 := int(bbox[2] * float64(w))
	y1 := int(bbox[3] * float64(h))
	return src.Region(image.Rect(x0, y0, x1, y1))
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}

func (m *ManualCropper) uploadBytesAsFile(data []byte, filename string) (string, error) {
	reader := &readSeekCloser{bytes.NewReader(data)} // Now implements multipart.File

	header := &multipart.FileHeader{
		Filename: filename,
		Size:     int64(len(data)),
		Header:   make(map[string][]string),
	}
	header.Header.Set("Content-Type", "image/jpeg")

	return m.fileStorage.UploadFile(reader, header)
}

func FromBBoxRequest(dbMeta request.BBoxMetaDB) (BBoxMeta, error) {
	var arr [4]float64
	err := json.Unmarshal(dbMeta.BBoxPercent, &arr)
	if err != nil {
		return BBoxMeta{}, err
	}
	return BBoxMeta{
		Page:        dbMeta.Page,
		BBoxPercent: arr,
	}, nil
}
