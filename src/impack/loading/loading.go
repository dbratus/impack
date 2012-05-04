package loading

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"strings"
)

const nParallelLoaders = 4

type loaderRequest struct {
	name   string
	reader io.ReadCloser
}

type Image struct {
	Name string
	Data image.Image
}

type Error struct {
	Name    string
	Message string
}

func LoadImages(pathChan chan string) ([]Image, []Error) {
	inChan := make(chan loaderRequest, nParallelLoaders)
	outChan := make(chan *Image, nParallelLoaders)
	errChan := make(chan *Error, nParallelLoaders)
	stopChan := make(chan int)

	images := make([]Image, 0, 100)
	errors := make([]Error, 0, 100)

	for i := 0; i < nParallelLoaders; i++ {
		go loader(inChan, outChan, errChan, stopChan)
	}

	cnt := 0
	for pth := <-pathChan; ; pth = <-pathChan {
		end := pth == ""

		if !end {
			if fl, err := os.Open(pth); err == nil {
				inChan <- loaderRequest{pth, fl}
				cnt++
			} else {
				fmt.Errorf("%s\n", err)
			}
		}

		if cnt == nParallelLoaders || end {
			getImagesFromChan(outChan, errChan, &images, &errors, &cnt)

			if end {
				break
			}
		}
	}

	for i := 0; i < nParallelLoaders; i++ {
		stopChan <- 1
	}

	return images, errors
}

func LoadImagesFromZip(reader *zip.Reader) ([]Image, []Error) {
	inChan := make(chan loaderRequest, nParallelLoaders)
	outChan := make(chan *Image, nParallelLoaders)
	errChan := make(chan *Error, nParallelLoaders)
	
	stopChan := make(chan int)

	images := make([]Image, 0, 100)
	errors := make([]Error, 0, 100)

	for i := 0; i < nParallelLoaders; i++ {
		go loader(inChan, outChan, errChan, stopChan)
	}

	cnt := 0
	for _, f := range reader.File {
		var isPng, isJpg bool

		if m, err := path.Match("*.png", f.Name); err == nil {
			isPng = m
		}

		if m, err := path.Match("*.jpg", f.Name); err == nil {
			isJpg = m
		}

		if !strings.HasPrefix(f.Name, "__MACOSX") && (isPng || isJpg) {
			if r, err := f.Open(); err == nil {
				inChan <- loaderRequest{f.Name, r}
				cnt++
			} else {
				fmt.Printf("%s\n", err)
			}
		}

		if cnt == nParallelLoaders {
			getImagesFromChan(outChan, errChan, &images, &errors, &cnt)
		}
	}

	getImagesFromChan(outChan, errChan, &images, &errors, &cnt)

	return images, errors
}

func getImagesFromChan(outChan chan *Image, errChan chan *Error, images *[]Image, errors *[]Error, cnt *int) {
	for ; *cnt > 0; *cnt-- {
		select {
		case img := <- outChan:
			*images = append(*images, *img)
		case err := <- errChan:
			*errors = append(*errors, *err)
		}
	}
}

func loader(inChan chan loaderRequest, outChan chan *Image, errChan chan *Error, stopChan chan int) {
	for {
		select {
		case req := <-inChan:
			switch strings.ToLower(path.Ext(req.name)) {
			case ".jpg":
				if img, err := jpeg.Decode(req.reader); err == nil {
					outChan <- &Image{Name: req.name, Data: img}
				} else {
					errChan <- &Error{Name: req.name, Message: err.Error() }
				}
			case ".png":
				if img, err := png.Decode(req.reader); err == nil {
					outChan <- &Image{Name: req.name, Data: img}
				} else {
					errChan <- &Error{Name: req.name, Message: err.Error() }
				}
			}

			if err := req.reader.Close(); err != nil {
				fmt.Errorf("%s\n", err)
			}
		case <-stopChan:
			return
		}
	}
}
