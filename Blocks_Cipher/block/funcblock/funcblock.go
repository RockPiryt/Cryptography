// Author: Paulina Kimak
package funcblock

import (
	"crypto/sha256"
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

// shaEncrypt simulates block encryption using SHA-256.
// It combines the input block and a key, hashes them, and returns
// the first N bytes (where N = 8*8 = 64) as the "encrypted" output.
func shaEncrypt(block, key []byte) []byte {
	// Create a new SHA-256 hash object
	hashObj := sha256.New()

	// Feed the block data into the hash
	hashObj.Write(block)

	// Also feed the key into the hash to simulate key-based encryption
	hashObj.Write(key)

	// Finalize the hash and take only the first 64 bytes
	// (because we want an 8x8 grayscale block = 64 pixels)
	return hashObj.Sum(nil)[:blockSize*blockSize] // truncate to block size
}


func ProcessECB(img *image.Gray, key []byte) *image.Gray {
	out := image.NewGray(img.Bounds())

	for y := 0; y < img.Bounds().Dy(); y += blockSize {
		for x := 0; x < img.Bounds().Dx(); x += blockSize {
			// Extract an 8x8 block of pixels from the image
			block := getBlock(img, x, y)
			// Encrypt the block using the shaEncrypt function
			encrypted := shaEncrypt(block, key)
			// Write the encrypted block back to the output image
			writeBlock(out, x, y, encrypted)
		}
	}
	return out
}

func ProcessCBC(img *image.Gray, key []byte) *image.Gray {
}
