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
