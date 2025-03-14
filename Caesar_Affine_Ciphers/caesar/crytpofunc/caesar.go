package crytpofunc

import (
	"strings"
)

//Author: Paulina Kimak

func CaesarCipher(text string, key int, encrypt bool) string {
	var result strings.Builder
	for _, char := range text {
		if char >= 'a' && char <= 'z' {
			shift := int(char) - 'a'
			if encrypt {
				shift = (shift + key) % 26
			} else {
				shift = (shift - key + 26) % 26
			}
			result.WriteByte(byte(shift + 'a'))
		} else if char >= 'A' && char <= 'Z' {
			shift := int(char) - 'A'
			if encrypt {
				shift = (shift + key) % 26
			} else {
				shift = (shift - key + 26) % 26
			}
			result.WriteByte(byte(shift + 'A'))
		} else {
			result.WriteByte(byte(char))
		}
	}
	return result.String()
}