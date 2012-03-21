package main

import (
	"impack"
	"fmt"
)

func main() {
	avgFillRate := impack.Stats(1, 100, 30)
	
	fmt.Printf("Avg. fill rate: %f\n", avgFillRate)
}