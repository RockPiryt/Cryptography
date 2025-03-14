package cryptofunc

import (
	"caesaraffineciphers/helpers"
	"fmt"
	"log"
	"strings"
)

//Author: Paulina Kimak

// ------------------------------------------------------------------------Caesar Cipher------------------------------------------------------------------------
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

func CaesarExecute(operation string) {
	var inputText, outputText, ipputTextHelper, outputkey string

	switch operation {
	case "e":
		inputText = "files/plain.txt"
		outputText = "files/crypto.txt"
	case "d":
		inputText = "files/crypto.txt"
		outputText = "files/decrypt.txt"
	case "j":
		inputText = "files/crypto.txt"
		ipputTextHelper = "files/extra.txt"
		outputText = "files/decrypt.txt"
		outputkey = "files/key-found.txt"
	default:
		fmt.Println("Nieobsługiwana operacja dla Cezara.")
		return
	}

	// Read text from input file
	textLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")

	// Read and validate the key
	key := helpers.ValidateCaesarKey("files/key.txt")

	// Encrypt or decrypt
	processedText := CaesarCipher(originalText, key, operation)

	// Save the result
	helpers.SaveOutput(processedText, outputFile)
}


func CaesarExplicitCryptAnalysis(operation string) string {
	var result strings.Builder
	return result.String()
}

func CaesarCryptAnalysis(operation string) string {
	var result strings.Builder
	return result.String()
}

// ------------------------------------------------------------------------Affine Cipher------------------------------------------------------------------------
func AffineCipher(text string, key int, operation string) string {
	return "a"
}
func AffineExecute(operation string) {
	var inputFile, outputFile string

	switch operation {
	case "e":
		inputFile = "files/plain.txt"
		outputFile = "files/crypto.txt"
	case "d":
		inputFile = "files/crypto.txt"
		outputFile = "files/decrypt.txt"
	case "j", "k":
		fmt.Println("Kryptoanaliza dla szyfru afinicznego nie jest jeszcze zaimplementowana.")
		return
	default:
		fmt.Println("Nieobsługiwana operacja dla szyfru afinicznego.")
		return
	}

	// Odczytaj tekst wejściowy
	textLines, err := helpers.GetText(inputFile)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")

	// Odczytaj i zweryfikuj klucze (a, b) dla szyfru afinicznego
	a, b := helpers.ValidateAffineKey("files/key.txt")

	// Wykonaj szyfrowanie lub deszyfrowanie
	processedText := AffineCipher(originalText, a, b, operation)

	// Zapisz wynik do pliku
	helpers.SaveOutput(processedText, outputFile)
}

func AffineExplicitCryptAnalysis(operation string) string {
	var result strings.Builder
	return result.String()
}

func AffineCryptAnalysis(operation string) string {
	var result strings.Builder
	return result.String()
}
