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

func LoadImages(pathChan chan string) []Image {
	loadingIn := make(chan loaderRequest, nParallelLoaders)
	loadingOut := make(chan *Image, nParallelLoaders)
	stopChan := make(chan int)

	images := make([]Image, 0, 100)

	for i := 0; i < nParallelLoaders; i++ {
		go loader(loadingIn, loadingOut, stopChan)
	}

	cnt := 0
	for pth := <-pathChan; ; pth = <-pathChan {
		end := pth == ""

		if !end {
			if fl, err := os.Open(pth); err == nil {
				loadingIn <- loaderRequest{pth, fl}
				cnt++
			} else {
				fmt.Errorf("%s\n", err)
			}
		}

		if cnt == nParallelLoaders || end {
			getImagesFromChan(loadingOut, &images, &cnt)

			if end {
				break
			}
		}
	}

	for i := 0; i < nParallelLoaders; i++ {
		stopChan <- 1
	}

	return images
}

func LoadImagesFromZip(reader *zip.Reader) []Image {
	loadingIn := make(chan loaderRequest, nParallelLoaders)
	loadingOut := make(chan *Image, nParallelLoaders)
	stopChan := make(chan int)

	images := make([]Image, 0, 100)

	for i := 0; i < nParallelLoaders; i++ {
		go loader(loadingIn, loadingOut, stopChan)
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
				loadingIn <- loaderRequest{f.Name, r}
				cnt++
			} else {
				fmt.Printf("%s\n", err)
			}
		}

		if cnt == nParallelLoaders {
			getImagesFromChan(loadingOut, &images, &cnt)
		}
	}

	getImagesFromChan(loadingOut, &images, &cnt)

	return images
}

func getImagesFromChan(loadingOut chan *Image, images *[]Image, cnt *int) {
	for ; *cnt > 0; *cnt-- {
		if img := <-loadingOut; img != nil {
			*images = append(*images, *img)
		}
	}
}

func loader(in chan loaderRequest, out chan *Image, stop chan int) {
	for {
		select {
		case req := <-in:
			switch strings.ToLower(path.Ext(req.name)) {
			case ".jpg":
				if img, err := jpeg.Decode(req.reader); err == nil {
					out <- &Image{ Name: req.name, Data: img }
				} else {
					fmt.Errorf("%s\n", err)
					out <- nil
				}
			case ".png":
				if img, err := png.Decode(req.reader); err == nil {
					out <- &Image{ Name: req.name, Data: img }
				} else {
					fmt.Errorf("%s\n", err)
					out <- nil
				}
			}

			if err := req.reader.Close(); err != nil {
				fmt.Errorf("%s\n", err)
			}
		case <-stop:
			return
		}
	}
}
