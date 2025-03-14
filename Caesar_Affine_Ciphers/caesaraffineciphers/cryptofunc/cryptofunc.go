package cryptofunc

import (
	"strings"
)

//Author: Paulina Kimak

// CaesarCipher encrypts or decrypts text using the Caesar cipher based on the given flags.
func CaesarCipher(text string, key int, operationFlag bool) string {
	var result strings.Builder

	for _, char := range text {
		// Handle lowercase letters
		if char >= 'a' && char <= 'z' {
			shift := int(char - 'a')
			if operationFlag { // Encrypt
				shift = (shift + key) % 26
			} else { // Decrypt
				shift = (shift - key + 26) % 26
			}
			result.WriteRune(rune(shift + 'a'))

		} else if char >= 'A' && char <= 'Z' {
			// Handle uppercase letters
			shift := int(char - 'A')
			if operationFlag { // Encrypt
				shift = (shift + key) % 26
			} else { // Decrypt
				shift = (shift - key + 26) % 26
			}
			result.WriteRune(rune(shift + 'A'))

		} else {
			// Non-alphabet characters remain unchanged
			result.WriteRune(char)
		}
	}

	return result.String()
}

func AffineCipher(text string, key int, operationFlag bool) string {
	var result strings.Builder
	return result.String()
}