package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"errors"
)

func isSupportedImageFormat(pth string) bool {
	ext := strings.ToLower(path.Ext(pth))
	return ext == ".png" || ext == ".jpg"
}

func getImagePaths(pth string) (chan string, error) {
	const pathBufSize = 10

	pathsChan := make(chan string, pathBufSize)

	var pathInfo os.FileInfo

	if pi, err := os.Stat(pth); err == nil {
		pathInfo = pi
	} else {
		return nil, err
	}

	if pathInfo.IsDir() {
		go func() {
			var dir *os.File

			if d, err := os.Open(pth); err == nil {
				dir = d
			} else {
				fmt.Errorf(err.Error())
				pathsChan <- ""
				return
			}

			for {
				var entryNames []string
				if entNames, err := dir.Readdirnames(pathBufSize); err == nil {
					entryNames = entNames
				} else {
					fmt.Errorf(err.Error())
					pathsChan <- ""
					return
				}

				for _, entName := range entryNames {
					p := pth + "/" + entName

					if pStat, err := os.Stat(p); err == nil {
						if !pStat.IsDir() && isSupportedImageFormat(p) {
							pathsChan <- p
						}
					} else {
						fmt.Errorf(err.Error())
					}
				}
			}
		}()
	} else if !pathInfo.IsDir() && isSupportedImageFormat(pth) {
		pathsChan <- pth
	} else {
		return nil, errors.New("The specified file is not an image file.")
	}

	return pathsChan, nil
}
