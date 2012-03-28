package main

import (
	"impack"
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/png"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: impack <images-dir> <output-png-file>")
		return
	}

	var images []image.Image

	if pathChan, err := getImagePaths(os.Args[1]); err == nil {
		images = loadImages(pathChan)
	} else {
		fmt.Errorf("%s\n", err)
		return
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
	
	if out, err := os.Create(os.Args[2]); err == nil {
		png.Encode(out, dest)
	} else {
		fmt.Errorf("%s\n", err)
		return
	}
}
