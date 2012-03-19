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

//Packes the specified image rectangles into a minimal area. Returns the list of 
//rectangles within the area.
func Pack(rects []image.Rectangle) (packed RectSlice) {
	packed = make([]image.Rectangle, 0, len(rects))

	if len(rects) == 0 {
		return
	}

	//Populating the result set with the input rectangels and sorting them by size.
	copy(packed, rects)
	sort.Sort(packed)

	//Feasible placements of any rectange around another rectangle.
	placements := make([]image.Rectangle, 8)

	for i := 1; i < len(packed); i++ {
		minRate := math.MaxFloat64
		var bestPlacement image.Rectangle

		for j := 0; j < i; j++ {
			makePlacements(packed[j], packed[i], placements)

			//Searching for a best placement among feasible.
			for k := 0; k < len(placements); k++ {
				if intersectsAny(placements[k], packed, i) {
					continue
				}

				if newRate := rateOf(placements[k], packed, i); newRate < minRate {
					minRate = newRate
					bestPlacement = placements[k]
				}
			}

			packed[i] = bestPlacement
		}
	}

	return
}
