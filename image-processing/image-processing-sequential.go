package imageprocessing

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func ProcessSequential(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		contentType, err := getFileContentType(path)
		if err != nil {
			return err
		}
		if contentType != "image/jpeg" {
			return nil
		}

		thumbnailImage, err := processImageSequential(path)
		if err != nil {
			return err
		}

		err = saveThumbnailImage(path, thumbnailImage)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func processImageSequential(path string) (*image.NRGBA, error) {
	srcImage, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}

	thumbnail := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)
	return thumbnail, nil
}

func saveThumbnailImage(srcImagePath string, thumbnailImage *image.NRGBA) error {
	filename := filepath.Base(srcImagePath)
	dstImagePath := "thumbnail/" + filename

	err := imaging.Save(thumbnailImage, dstImagePath)
	if err != nil {
		return err
	}

	fmt.Printf("%s -> %s\n", srcImagePath, dstImagePath)
	return nil
}

func getFileContentType(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	fileContentType := http.DetectContentType(buffer)
	return fileContentType, nil
}