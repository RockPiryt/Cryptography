// Author: Paulina Kimak
package main

import (
	"fmt"
	"block/helpers"
)

func main() {
	img, err := helpers.LoadImage("files/plain.bmp")
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	grayImg := helpers.ConvertToGrayscale(img)

}
