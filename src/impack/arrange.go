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
package impack

import (
	"image"
	"sort"
	"math"
)

//Makes all placements of r around the pivot. 
func makePlacements(pivot, r image.Rectangle, out []image.Rectangle) {
	out[0] = image.Rect(pivot.Max.X, pivot.Min.Y, pivot.Max.X+r.Dx(), pivot.Min.Y+r.Dy())
	out[1] = image.Rect(pivot.Max.X, pivot.Max.Y-r.Dy(), pivot.Max.X+r.Dx(), pivot.Max.Y)
	out[2] = image.Rect(pivot.Max.X-r.Dx(), pivot.Max.Y, pivot.Max.X, pivot.Max.Y+r.Dy())
	out[3] = image.Rect(pivot.Min.X, pivot.Max.Y, pivot.Min.X+r.Dx(), pivot.Max.Y+r.Dy())
	out[4] = image.Rect(pivot.Min.X-r.Dx(), pivot.Max.Y-r.Dy(), pivot.Min.X, pivot.Max.Y)
	out[5] = image.Rect(pivot.Min.X-r.Dx(), pivot.Min.Y, pivot.Min.X, pivot.Min.Y+r.Dy())
	out[6] = image.Rect(pivot.Min.X, pivot.Min.Y-r.Dy(), pivot.Min.X+r.Dx(), pivot.Min.Y)
	out[7] = image.Rect(pivot.Max.X-r.Dx(), pivot.Min.Y-r.Dy(), pivot.Max.X, pivot.Min.Y)
}

//Given a set of rectangles with top-left corner at (0,0)
//arranges them so that they occupy without intersection a 
//minimal area with minimal wasted space.
//Returns the union of the arranged rectangles.
func Arrange(rects []image.Rectangle) image.Rectangle {
	if len(rects) == 0 {
		return image.ZR
	}
	
	var arranged RectPtrSlice = make([]*image.Rectangle, len(rects))

	//Populating the result set with the input rectangels and sorting them by size.
	for i := 0; i < len(arranged); i++ {
		arranged[i] = &rects[i]
	}
	sort.Sort(arranged)

	//Feasible placements of any rectange around another rectangle.
	placements := make([]image.Rectangle, 8)

	totalAreaOfArranged := area(*arranged[0])
	union := *arranged[0]

	for i := 1; i < len(arranged); i++ {
		minRate := math.MaxFloat64
		var bestPlacement image.Rectangle

		for j := 0; j < i; j++ {
			makePlacements(*arranged[j], *arranged[i], placements)

			//Searching for a best placement among feasible.
			for k := 0; k < len(placements); k++ {
				if intersectsAny(placements[k], arranged, i) {
					continue
				}

				if newRate := rateOf(placements[k], union, totalAreaOfArranged); newRate < minRate {
					minRate = newRate
					bestPlacement = placements[k]
				}
			}

			*arranged[i] = bestPlacement
		}
		
		totalAreaOfArranged += area(*arranged[i])
		union = union.Union(*arranged[i])
	}

	//Aligning the rectangles.
	for i := 0; i < len(rects); i++ {
		dx := rects[i].Dx()
		dy := rects[i].Dy()
	
		rects[i].Min.X = rects[i].Min.X - union.Min.X
		rects[i].Min.Y = rects[i].Min.Y - union.Min.Y
		rects[i].Max.X = rects[i].Min.X + dx 
		rects[i].Max.Y = rects[i].Min.Y + dy 
	}
	
	union = image.Rect(0, 0, union.Max.X - union.Min.X, union.Max.Y - union.Min.Y)
	
	return union
}
