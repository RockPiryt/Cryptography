// Author: Paulina Kimak
package main

import (
	"fmt"
	"block/helpers"
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

}
