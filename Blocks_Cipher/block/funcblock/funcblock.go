// Author: Paulina Kimak
package funcblock

import (
	"bytes"
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

// shaEncrypt simulates block encryption using SHA-256. It hashes the block + key, repeats the hash if needed, and returns
// the first N bytes (N = 8*8 = 64) to represent the encrypted block.
func shaEncrypt(block, key []byte) []byte {
	// Combine block and key, then hash them with SHA-256
	hash := sha256.Sum256(append(block, key...)) // 32-byte result

	// Prepare output slice of 64 bytes (8x8 block)
	out := make([]byte, blockSize*blockSize) 

	// Repeat the hash bytes as needed to fill the output
	copy(out, bytes.Repeat(hash[:], (len(out)/len(hash[:]))+1))

	return out
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

// xorBlocks takes two equal-length byte slices and returns a new slice
// where each byte is the result of XOR'ing corresponding bytes from both inputs.
func xorBlocks(a, b []byte) []byte {
	res := make([]byte, len(a))

	for i := range a {
		res[i] = a[i] ^ b[i]
	}

	return res
}


func ProcessCBC(img *image.Gray, key []byte) *image.Gray {
	out := image.NewGray(img.Bounds())
	iv := make([]byte, blockSize*blockSize) // init vector = 0s
	prev := iv

	for y := 0; y < img.Bounds().Dy(); y += blockSize {
		for x := 0; x < img.Bounds().Dx(); x += blockSize {
			block := getBlock(img, x, y)
			xored := xorBlocks(block, prev)
			encrypted := shaEncrypt(xored, key)
			writeBlock(out, x, y, encrypted)
			prev = encrypted
		}
	}
	return out
}
