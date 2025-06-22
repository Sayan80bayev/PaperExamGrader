package service

type Cutter interface {
	Crop(batchDir string) error
}
