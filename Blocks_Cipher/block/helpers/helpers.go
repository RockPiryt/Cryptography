// Author: Paulina Kimak
package helpers

import (
	"image"
	"image/bmp"
	"os"
)

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := bmp.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}