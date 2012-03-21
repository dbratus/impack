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
	out[6] = image.Rect(pivot.Min.X, pivot.Min.Y-pivot.Dy(), pivot.Min.X+r.Dx(), pivot.Min.Y)
	out[7] = image.Rect(pivot.Max.X-r.Dx(), pivot.Min.Y-r.Dy(), pivot.Max.X, pivot.Min.Y)
}

//Given a set of rectangles with top-left corner at (0,0)
//arranges them so that they occupy without intersection a 
//minimal area with minimal wasted space.
func Arrange(rects []image.Rectangle) {
	if len(rects) == 0 {
		return
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
}
