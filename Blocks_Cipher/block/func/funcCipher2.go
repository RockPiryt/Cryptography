package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/bmp"
	"io/ioutil"
	"os"
)

const blockSize = 8 // 8x8 pixel blocks

func readKey(path string) []byte {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte("default-key")
	}
	return key
}

func loadImage(path string) (image.Image, error) {
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

func saveImage(path string, img *image.Gray) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return bmp.Encode(out, img)
}

func getBlock(img *image.Gray, x, y int) []byte {
	block := make([]byte, 0, blockSize*blockSize)
	for dy := 0; dy < blockSize; dy++ {
		for dx := 0; dx < blockSize; dx++ {
			px := img.GrayAt(x+dx, y+dy).Y
			block = append(block, px)
		}
	}
	return block
}

func writeBlock(img *image.Gray, x, y int, data []byte) {
	idx := 0
	for dy := 0; dy < blockSize; dy++ {
		for dx := 0; dx < blockSize; dx++ {
			val := data[idx]
			img.SetGray(x+dx, y+dy, color.Gray{Y: val})
			idx++
		}
	}
}

func fakeEncrypt(block, key []byte) []byte {
	h := sha256.New()
	h.Write(block)
	h.Write(key)
	return h.Sum(nil)[:blockSize*blockSize] // truncate to block size
}

func processECB(img *image.Gray, key []byte) *image.Gray {
	out := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y += blockSize {
		for x := 0; x < img.Bounds().Dx(); x += blockSize {
			block := getBlock(img, x, y)
			encrypted := fakeEncrypt(block, key)
			writeBlock(out, x, y, encrypted)
		}
	}
	return out
}

func xorBlocks(a, b []byte) []byte {
	res := make([]byte, len(a))
	for i := range a {
		res[i] = a[i] ^ b[i]
	}
	return res
}

func processCBC(img *image.Gray, key []byte) *image.Gray {
	out := image.NewGray(img.Bounds())
	iv := make([]byte, blockSize*blockSize) // init vector = 0s
	prev := iv

	for y := 0; y < img.Bounds().Dy(); y += blockSize {
		for x := 0; x < img.Bounds().Dx(); x += blockSize {
			block := getBlock(img, x, y)
			xored := xorBlocks(block, prev)
			encrypted := fakeEncrypt(xored, key)
			writeBlock(out, x, y, encrypted)
			prev = encrypted
		}
	}
	return out
}

func main() {
	img, err := loadImage("plain.bmp")
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	grayImg := image.NewGray(img.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	key := readKey("key.txt")
	ecb := processECB(grayImg, key)
	cbc := processCBC(grayImg, key)

	saveImage("ecb_crypto.bmp", ecb)
	saveImage("cbc_crypto.bmp", cbc)

	fmt.Println("ECB and CBC images saved.")
}