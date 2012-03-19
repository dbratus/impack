package impack

import (
	"image"
)

type RectSlice []image.Rectangle

func (s RectSlice) Len() int {
	return len(s)
}

func (s RectSlice) Less(i, j int) bool {
	return area(s[i]) < area(s[j])
}

func (s RectSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Returns the area of a rectangle.
func area(r image.Rectangle) int {
	return r.Dx() * r.Dy()
}

//Returns the total area of all rectangles.
func totalArea(rects []image.Rectangle, limit int) int {
	a := 0
	for i := 0; i < limit; i++ {
		a += area(rects[i])
	}
	
	return a
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

//Returns the aspect ration of the union of r and rects.
func aspectRatioOfUnion(r image.Rectangle, rects []image.Rectangle, limit int) float64 {
	union := Union(rects, limit)
	union = union.Union(r)

	return aspectRatio(union)
}

//Returns the rate of the union of r and rects. 
func rateOf(r image.Rectangle, rects []image.Rectangle, limit int) float64 {
	return float64(wastedArea(r, rects, limit)) * aspectRatioOfUnion(r, rects, limit)
}

//Checks if the rectange intersects with any of the rectangles specified.
func intersectsAny(r image.Rectangle, rects []image.Rectangle, limit int) bool {
	for i := 0; i < limit; i++ {
		if !rects[i].Intersect(r).Empty() {
			return true
		}
	}

	return false
}

//Returns the union of rectangles.
func Union(rects []image.Rectangle, limit int) image.Rectangle {
	if len(rects) == 0 {
		return image.ZR
	}

	union := rects[0]
	
	for i := 1; i < limit; i++ {
		union = union.Union(rects[i])
	}

	return union	
}

//Returns the area wasted in the union of the specified rectangle and
//the specified set of rectangles. 
func wastedArea(r image.Rectangle, rects []image.Rectangle, limit int) int {
	union := Union(rects, limit)
	union = union.Union(r)

	return area(union) - totalArea(rects, limit) - area(r)
}
