package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

type imageService struct{}

func NewImageService() ImageService {
	return &imageService{}
}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, err
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := is.imagePath(galleryID)
	images, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	for i := range images {
		images[i] = "/" + images[i]
	}
	return images, nil
}
