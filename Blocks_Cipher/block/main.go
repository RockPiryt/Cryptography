// Author: Paulina Kimak
package main

import (
	"fmt"
	"block/helpers"
	"block/funcblock"
)

func main() {
	// Load the image
	img, err := helpers.LoadImage("files/plain.bmp")
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	// Convert the image to grayscale
	grayImg := helpers.ConvertToGrayscale(img)

	// Read the key from the file
	key := helpers.ReadKey("files/key.txt") 

	// Process the image using ECB and CBC modes
	ecb := funcblock.ProcessECB(grayImg, key)
	cbc := funcblock.ProcessCBC(grayImg, key)

	// Save the processed images
	helpers.SaveImage("files/ecb_crypto.bmp", ecb)
	helpers.SaveImage("files/cbc_crypto.bmp", cbc)

	fmt.Println("ECB and CBC images saved.")

}
