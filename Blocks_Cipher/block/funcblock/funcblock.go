// Author: Paulina Kimak
package funcblock

import (
	"image"
	"image/color"
)

const blockSize = 8 // 8x8 pixel blocks

// getBlock extracts an 8x8 block of pixels from the grayscale image
func getBlock(img *image.Gray, x, y int) []byte {
	// Preallocate 1D byte slice to hold 64 grayscale values (for 8x8 block)
	block := make([]byte, 0, blockSize*blockSize)

	// Loop over block rows (dy) and columns (dx)
	for dy := 0; dy < blockSize; dy++ {
		for dx := 0; dx < blockSize; dx++ {
			// Get the brightness (Y) value at pixel (x+dx, y+dy)
			px := img.GrayAt(x+dx, y+dy).Y
			block = append(block, px)
		}
	}
	return block
}

// writeBlock takes a 1D byte slice of 64 grayscale values (8x8)
// and writes it into the grayscale image starting at coordinate (x, y).
func writeBlock(img *image.Gray, x, y int, data []byte) {
	idx := 0 

	// Loop over block rows and columns
	for dy := 0; dy < blockSize; dy++ {
		for dx := 0; dx < blockSize; dx++ {
			// Get the grayscale value from data slice
			val := data[idx]

			// Set the pixel at (x+dx, y+dy) with that grayscale value
			img.SetGray(x+dx, y+dy, color.Gray{Y: val})

			idx++ 
		}
	}
}



func ProcessECB(img *image.Gray, key []byte) *image.Gray {
}

func ProcessCBC(img *image.Gray, key []byte) *image.Gray {
}
