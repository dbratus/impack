package web

import (
	"fmt"
	"image"
	"impack/loading"
	"appengine"
	"io"
)

const cssFmt = 
`.cls%d {
    background: url(/blob/%s) %dpx %dpx;
    width: %dpx;
    height: %dpx;
}
`

func generateCss(w io.Writer, imageKey appengine.BlobKey, rects []image.Rectangle, images []loading.Image) error {
	for i := 0; i < len(rects); i++ {
		if _, err := fmt.Fprintf(w, "\n/*%s*/\n", images[i].Name); err != nil {
			return err
		}
		
		if _, err := fmt.Fprintf(
			w,
			cssFmt,
			i + 1,
			imageKey,
			-rects[i].Min.X, -rects[i].Min.Y,
			rects[i].Dx(), rects[i].Dy()); err != nil {
			
			return err
		}
	}

	return nil
}
