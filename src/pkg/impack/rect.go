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
