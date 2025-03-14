package cryptofunc

import (
	"fmt"
	"log"
	"strings"
	"caesaraffineciphers/helpers"
)

//Author: Paulina Kimak

// CaesarCipher encrypts or decrypts text using the Caesar cipher based on the given flags.
func CaesarCipher(text string, key int, operation string) string {
	var result strings.Builder

	for _, char := range text {
		// Handle encryption for lowercase letters and uppercase letters
		if operation == "e" { 
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift + key) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' { 
				shift := int(char - 'A')
				shift = (shift + key) % 26
				result.WriteRune(rune(shift + 'A'))
			}
		}

		// Handle decryption for lowercase letters and uppercase letters
		if operation == "d" { 
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift - key + 26) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' { 
				shift := int(char - 'A')
				shift = (shift - key + 26) % 26
				result.WriteRune(rune(shift + 'A'))
			}
		}

		// Non-alphabet characters remain unchanged
		if !(char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z') {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func CaesarEncrypt(operation string) {
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

	// Encrypt the text.
	encodedTxt := CaesarCipher(originalText, key, operation)
	helpers.SaveOutput(encodedTxt, "files/crypto.txt")
}

func CaesarDecrypt(operation string) {
	// Read cipher text.
	inputFileDecrypt := "files/crypto.txt"
	textLines, err := helpers.GetText(inputFileDecrypt)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	encodedText := strings.Join(textLines, "\n")
	fmt.Println(encodedText)
	
	// Read and check the key.
	key := helpers.ValidateCaesarKey("files/key.txt")

	// Decrypt the text.
	encodedTxt := CaesarCipher(encodedText, key, operation)
	helpers.SaveOutput(encodedTxt, "files/decrypt.txt")
}


func AffineCipher(text string, key int, operation string) string {
	var result strings.Builder
	return result.String()
}