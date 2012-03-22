package main

import (
	//"impack"
	"fmt"
	"os"
	"path"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	var pathChan chan string
	var err os.Error

	if pathChan, err = getImagePaths(os.Args[1]); err != nil {
		fmt.Errorf(err.String())
		return
	}

	for pth := <- pathChan; pth != ""; pth = <- pathChan {
		fmt.Println(pth)
	}
}

func isSupportedImageFormat(pth string) bool {
	ext := strings.ToLower(path.Ext(pth))
	return ext == ".png" || ext == ".jpg"
}

func getImagePaths(pth string) (chan string, os.Error) {
	const pathBufSize = 10

	pathsChan := make(chan string, pathBufSize)

	var err os.Error
	var pathInfo *os.FileInfo

	if pathInfo, err = os.Stat(pth); err != nil {
		return nil, err
	}

	if pathInfo.IsDirectory() {
		go func() {
			var dir *os.File

			if dir, err = os.Open(pth); err != nil {
				fmt.Errorf(err.String())
				pathsChan <- ""
				return
			}

			for {
				var entryNames []string
				if entryNames, err = dir.Readdirnames(pathBufSize); err != nil {
					if err != os.EOF {
						fmt.Errorf(err.String())
					}
					pathsChan <- ""
					return
				}

				for _, entName := range entryNames {
					p := pth + "/" + entName

					var pStat *os.FileInfo

					if pStat, err = os.Stat(p); err == nil {
						if pStat.IsRegular() && isSupportedImageFormat(p) {
							pathsChan <- p
						}
					} else {
						fmt.Errorf(err.String())
					}
				}
			}
		}()
	} else if pathInfo.IsRegular() && isSupportedImageFormat(pth) {
		pathsChan <- pth
	} else {
		return nil, os.NewError("The specified file is not an image file.")
	}

	return pathsChan, nil
}
