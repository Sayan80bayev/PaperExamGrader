package service

import (
	"PaperExamGrader/internal/model"
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-fitz" // PDF → image
	"gocv.io/x/gocv"               // OpenCV wrapper
)

type BBoxMeta struct {
	Page        int
	BBoxPercent [4]float64
}

// ManualCropper implements Cutter using predefined BBoxMeta entries.
type ManualCropper struct {
	bboxMap []BBoxMeta
}

// SetCropParams sets the cropping parameters for manual cropping.
func (m *ManualCropper) SetCropParams(bboxes []BBoxMeta) {
	m.bboxMap = bboxes
}

// Crop applies the crop parameters to all PDFs in batchDir.
func (m *ManualCropper) Crop(batchDir string, outputDir string) error {
	files, err := os.ReadDir(batchDir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	for _, entry := range files {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".pdf") {
			continue
		}
		pdfPath := filepath.Join(batchDir, entry.Name())
		if err := cropAnswersFromPDF(pdfPath, outputDir, m.bboxMap); err != nil {
			log.Println("Error cropping", entry.Name(), ":", err)
		}
	}
	return nil
}

func cropFromPercentBBox(src gocv.Mat, bbox [4]float64) gocv.Mat {
	h, w := src.Rows(), src.Cols()
	x0 := int(bbox[0] * float64(w))
	y0 := int(bbox[1] * float64(h))
	x1 := int(bbox[2] * float64(w))
	y1 := int(bbox[3] * float64(h))
	return src.Region(image.Rect(x0, y0, x1, y1))
}

func cropAnswersFromPDF(pdfPath, outputDir string, bboxMap []BBoxMeta) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("fitz.New: %w", err)
	}
	defer doc.Close()

	base := filepath.Base(pdfPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	for idx, bm := range bboxMap {
		if bm.Page >= doc.NumPage() {
			continue
		}

		img, err := doc.Image(bm.Page)
		if err != nil {
			return fmt.Errorf("doc.Image(page %d): %w", bm.Page, err)
		}

		mat, err := gocv.ImageToMatRGB(img)
		if err != nil {
			return fmt.Errorf("ImageToMatRGB: %w", err)
		}
		defer mat.Close()

		cropped := cropFromPercentBBox(mat, bm.BBoxPercent)
		defer cropped.Close()

		outPath := filepath.Join(outputDir, fmt.Sprintf("%s_answer_%d.jpg", name, idx+1))
		if !gocv.IMWrite(outPath, cropped) {
			log.Printf("❌ Failed to write %s", outPath)
		} else {
			log.Printf("✅ Saved %s", outPath)
		}
	}
	return nil
}

func ToBBoxMetaDB(meta BBoxMeta) (model.BBoxMetaDB, error) {
	data, err := json.Marshal(meta.BBoxPercent)
	if err != nil {
		return model.BBoxMetaDB{}, err
	}
	return model.BBoxMetaDB{
		Page:        meta.Page,
		BBoxPercent: data,
	}, nil
}

func FromBBoxMetaDB(dbMeta model.BBoxMetaDB) (BBoxMeta, error) {
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
