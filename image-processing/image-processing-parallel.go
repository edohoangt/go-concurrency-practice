package imageprocessing

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

type processResult struct {
	srcImagePath string
	thumbnailImage *image.NRGBA
	err error
}

// set up pipeline
func ProcessParallel(root string) error {
	doneCh := make(chan struct{})
	defer close(doneCh)

	// first stage
	pathsCh, errsCh := walkFilesParallel(doneCh, root)

	// second stage
	resultsCh := processImageParallel(doneCh, pathsCh)

	// third stage
	for r := range resultsCh {
		if r.err != nil {
			return r.err
		}
		saveThumbnailImage(r.srcImagePath, r.thumbnailImage)
	}

	if err := <-errsCh; err != nil {
		return err
	}

	return nil
}

func processImageParallel(doneCh <-chan struct{}, pathsCh <-chan string) <-chan *processResult {
	resultsCh := make(chan *processResult)

	thumbnailer := func ()  {
		for path := range pathsCh {
			srcImage, err := imaging.Open(path)
			if err != nil {
				select {
				case resultsCh <- &processResult{path, nil, err}:
				case <- doneCh:
					return
				}
			}

			thumbnail := imaging.Thumbnail(srcImage, 100, 100, imaging.Lanczos)
			select {
			case resultsCh <- &processResult{path, thumbnail, nil}:
			case <- doneCh:
				return
			}
		}
	}

	const numThumbnailer = 6
	var wg sync.WaitGroup
	wg.Add(numThumbnailer)

	for i := 0; i < numThumbnailer; i++ {
		go func ()  {
			thumbnailer()
			wg.Done()
		}()
	}

	go func ()  {
		wg.Wait()
		close(resultsCh)
	}()

	return resultsCh
}

func walkFilesParallel(doneCh <-chan struct{}, root string) (<-chan string, <-chan error) {
	pathsCh := make(chan string)
	errsCh := make(chan error, 1)

	go func () {
		defer close(pathsCh)

		errsCh <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

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

			select {
			case pathsCh <- path:
			case <- doneCh:
				return fmt.Errorf("walk was cancelled")
			}
			
			return nil
		})
	}()

	return pathsCh, errsCh
}
