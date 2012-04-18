package main

import (
	"impack"
	"impack/loading"
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/png"
	"flag"
	"io"
	"archive/zip"
)

var archiveFileName *string = flag.String("a", "", "Archive file name.")
var imageDirName *string = flag.String("d", "", "Image directory name.")
var outputFileName *string = flag.String("o", "", "Output file name.")

func main() {
	flag.Parse()

	var images []image.Image
	
	if len(*archiveFileName) == 0 {
		var pathChan chan string
		
		if len(*imageDirName) > 0 {
			if pch, err := getImagePaths(*imageDirName); err == nil {
				pathChan = pch
			} else {
				fmt.Errorf("%s\n", err.Error())
				return
			}
		} else {
			pathChan = make(chan string, 10)
			go func() {
				var line string
				for _, err := fmt.Scan(&line); err == nil; _, err = fmt.Scan(&line) {
					pathChan <- line
				}
				
				pathChan <- ""
			}()
		}
		
		images = loading.LoadImages(pathChan)
	} else {
		var reader *zip.ReadCloser

		if r, err := zip.OpenReader(*archiveFileName); err == nil {
			reader = r
		} else {
			fmt.Printf("%s\n", err)
			return
		}
		
		images = loading.LoadImagesFromZip(&reader.Reader)
		reader.Close()
	}

	rects := make([]image.Rectangle, len(images))
	
	for i := 0; i < len(rects); i++ {
		rects[i] = images[i].Bounds()
	}
	
	union := impack.Arrange(rects)
	
	dest := image.NewNRGBA(image.Rect(0, 0, union.Dx(), union.Dy()))
	
	for i := 0; i < len(rects); i++ {
		draw.Draw(dest, rects[i], images[i], image.Pt(0, 0), draw.Src)
	}
	
	var out io.Writer
	
	if len(*outputFileName) > 0 {
		if f, err := os.Create(*outputFileName); err == nil {
			out = f
		} else {
			fmt.Errorf("%s\n", err)
			return
		}
	} else {
		out = os.Stdout
	}
	
	png.Encode(out, dest)
}
