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
)

type RectPtrSlice []*image.Rectangle

func (s RectPtrSlice) Len() int {
	return len(s)
}

func (s RectPtrSlice) Less(i, j int) bool {
	return area(*s[i]) > area(*s[j])
}

func (s RectPtrSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Returns the total area of all rectangles.
func totalArea(rects RectPtrSlice, limit int) int {
	a := 0
	for i := 0; i < limit; i++ {
		a += area(*rects[i])
	}
	
	return a
}

//Returns the rate of the union of r and rects. 
func rateOf(r image.Rectangle, union image.Rectangle, totalArea int) float64 {
	unionAndR := union.Union(r)
	
	wastedArea := area(unionAndR) - totalArea - area(r)
	return float64(wastedArea) * aspectRatio(unionAndR)
}

//Checks if the rectange intersects with any of the rectangles specified.
func intersectsAny(r image.Rectangle, rects []*image.Rectangle, limit int) bool {
	for i := 0; i < limit; i++ {
		if !(*rects[i]).Intersect(r).Empty() {
			return true
		}
	}

	return false
}

//Returns the union of rectangles.
func unionOf(rects RectPtrSlice, limit int) image.Rectangle {
	if len(rects) == 0 {
		return image.ZR
	}

	union := *rects[0]
	
	for i := 1; i < limit; i++ {
		union = union.Union(*rects[i])
	}

	return union	
}
