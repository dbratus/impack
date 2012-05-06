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

import "image"

//Returns the area of a rectangle.
func area(r image.Rectangle) int {
	return r.Dx() * r.Dy()
}

//Returns the aspect ratio of a rectangle.
func aspectRatio(r image.Rectangle) float64 {
	size := r.Size()
	var min, max int

	if size.X < size.Y {
		min = size.X
		max = size.Y
	} else {
		min = size.Y
		max = size.X
	}

	return float64(max) / float64(min)
}
