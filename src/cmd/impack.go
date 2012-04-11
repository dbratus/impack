package main

import (
	"impack"
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/png"
	"flag"
	"io"
)

var archiveFileName *string = flag.String("a", "", "Archive file name.")
var imageDirName *string = flag.String("d", "", "Image directory name.")
var outputFileName *string = flag.String("o", "", "Output file name.")

func main() {
	flag.Parse()

	var images []image.Image
	
	if len(*archiveFileName) == 0 {
		var err os.Error
		var pathChan chan string
		
		if len(*imageDirName) > 0 {
			if pathChan, err = getImagePaths(*imageDirName); err != nil {
				fmt.Errorf("%s\n", err)
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
		
		images = loadImages(pathChan)
	} else {
		images = loadImagesFromZip(*archiveFileName)
	}

	rects := make([]image.Rectangle, len(images))
	
	for i := 0; i < len(rects); i++ {
		rects[i] = images[i].Bounds()
	}
	
	union := impack.Arrange(rects)
	
	dest := image.NewNRGBA(union.Dx(), union.Dy())
	
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
