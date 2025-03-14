package cryptofunc

import (
	"fmt"
	"log"
	"strings"
	"caesaraffineciphers/helpers"
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

func CaesarEncrypt(operationFlag bool) {
	// Read plain text.
	inputFileEncrypt := "files/plain.txt"
	textLines, err := helpers.GetText(inputFileEncrypt)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")
	fmt.Println(originalText)
	
	// Read and check the key.
	key := helpers.ValidateCaesarKey("files/key.txt")

	encodedTxt := CaesarCipher(originalText, key, *operationFlag)
	helpers.SaveOutput(encodedTxt, "files/crypto.txt")
}

func CaesarDecrypt(operationFlag bool) {
	// Read plain text.
	inputFileDecrypt := "files/crypto.txt"
	textLines, err := helpers.GetText(inputFileDecrypt)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")
	fmt.Println(originalText)
	
	// Read and check the key.
	key := helpers.ValidateCaesarKey("files/key.txt")

	encodedTxt := CaesarCipher(originalText, key, *operationFlag)
	helpers.SaveOutput(encodedTxt, "files/decrypt.txt")
}
func AffineCipher(text string, key int, operationFlag bool) string {
	var result strings.Builder
	return result.String()
}