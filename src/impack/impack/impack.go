// ImPack - CSS sprites maker
// Copyright (C) 2012 Dmitry Bratus
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
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

	var images []loading.Image
	var loadingErrors []loading.Error
	
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
		
		images, loadingErrors = loading.LoadImages(pathChan)
	} else {
		var reader *zip.ReadCloser

		if r, err := zip.OpenReader(*archiveFileName); err == nil {
			reader = r
		} else {
			fmt.Printf("%s\n", err)
			return
		}
		
		images, loadingErrors = loading.LoadImagesFromZip(&reader.Reader)
		reader.Close()
	}

	for _, r := range loadingErrors {
		fmt.Errorf("Loading of %s failed:\n  %s\n", r.Name, r.Message)
	}

	rects := make([]image.Rectangle, len(images))
	
	for i := 0; i < len(rects); i++ {
		rects[i] = images[i].Data.Bounds()
	}
	
	union := impack.Arrange(rects)
	
	dest := image.NewNRGBA(image.Rect(0, 0, union.Dx(), union.Dy()))
	
	for i := 0; i < len(rects); i++ {
		draw.Draw(dest, rects[i], images[i].Data, image.Pt(0, 0), draw.Src)
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
