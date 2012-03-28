package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"path"
	"strings"
	"os"
	"fmt"
)

const nParallelLoaders = 4

func loadImages(pathChan chan string) []image.Image {
	loadingIn := make(chan string, nParallelLoaders)
	loadingOut := make(chan image.Image, nParallelLoaders)
	stopChan := make(chan int)

	images := make([]image.Image, 0, 100)

	for i := 0; i < nParallelLoaders; i++ {
		go loader(loadingIn, loadingOut, stopChan)
	}
	
	cnt := 0
	for pth := <-pathChan; ; pth = <-pathChan {
		end := pth == ""

		if !end {
			loadingIn <- pth
			cnt++
		}

		if cnt == nParallelLoaders || end {
			for ; cnt > 0; cnt-- {
				if img := <-loadingOut; img != nil {
					images = append(images, img)
				}
			}

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

func loader(in chan string, out chan image.Image, stop chan int) {
	for {
		select {
		case pth := <-in:
			if fl, err := os.Open(pth); err == nil {
				switch strings.ToLower(path.Ext(pth)) {
				case ".jpg":
					if img, err := jpeg.Decode(fl); err == nil {
						out <- img
					} else {
						fmt.Errorf("%s\n", err)
						out <- nil
					}
				case ".png":
					if img, err := png.Decode(fl); err == nil {
						out <- img
					} else {
						fmt.Errorf("%s\n", err)
						out <- nil
					}
				}

				if err = fl.Close(); err != nil {
					fmt.Errorf("%s\n", err)
				}
			} else {
				fmt.Errorf("%s\n", err)
				out <- nil
			}
		case <-stop:
			return
		}
	}
}
