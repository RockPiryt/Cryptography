// Author: Paulina Kimak
package helpers

import (
	"image"
	"golang.org/x/image/bmp" 
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

// convertToGrayscale takes any image and returns its grayscale version.
// It goes pixel by pixel and copies each color into a new grayscale image.
func ConvertToGrayscale(img image.Image) *image.Gray {
	// Create a new grayscale image with the same size as the original
	grayImg := image.NewGray(img.Bounds())

	// Loop over every pixel 
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			// Get the original pixel color at (x, y)
			originalColor := img.At(x, y)
			// The Set() method will automatically convert color to grayscale
			grayImg.Set(x, y, originalColor)
		}
	}
	return grayImg
}

func ReadKey(path string) []byte {
	key, err := os.ReadFile(path)
	if err != nil {
		// If the key file is not found, return a default key
		return []byte("default-key")
	}
	return key
}

// saveImage creates the file, encodes the image in BMP format, and writes it to disk.
func SaveImage(path string, img *image.Gray) error {
	// Try to create (or overwrite) the file at the given path
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Encode the grayscale image as a BMP file and write it to disk
	// bmp.Encode will write binary BMP format into the 'out' file
	return bmp.Encode(out, img)
}
