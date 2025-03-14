package cryptofunc

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"caesaraffineciphers/helpers"
)

//Author: Paulina Kimak

// Generic function to handle both Caesar and Affine ciphers
func ExecuteCipher(cipherType string, operation string) {
	var inputText, inputTextHelper, inputKey, outputText, outputKey string

	// Define file paths based on operation
	switch operation {
	case "e":
		// Program szyfrujący czyta tekst jawny i klucz i zapisuje tekst zaszyfrowany. Jeśli klucz jest nieprawidłowy, zgłasza jedynie błąd.
		inputText = "files/plain.txt"
		inputKey = "files/key.txt"
		outputText = "files/crypto.txt"
	case "d":
		// Program odszyfrowujący czyta tekst zaszyfrowany i klucz i zapisuje tekst jawny. Jeśli klucz jest nieprawidłowy, zgłasza błąd. 
		// Dla szyfru afinicznego częścią zadania jest znalezienie odwrotności dla liczby a podanej jako część klucza – 
		// nie można zakładać, że program odszyfrowujący otrzymuje tę odwrotność.
		inputText = "files/crypto.txt"
		inputKey = "files/key.txt"
		outputText = "files/decrypt.txt"
	case "j":
		// Program łamiący szyfr z pomocą tekstu jawnego czyta tekst zaszyfrowany, tekst pomocniczy i zapisuje znaleziony klucz i odszyfrowany tekst. 
		// Jeśli niemożliwe jest znalezienie klucza, zgłasza sygnał błędu.
		inputText = "files/crypto.txt"
		inputTextHelper = "files/extra.txt"
		outputText = "files/decrypt.txt"
		outputKey = "files/key-found.txt"
	case "k":
		//Program łamiący szyfr bez pomocy tekstu jawnego czyta jedynie tekst zaszyfrowany i zapisuje jako tekst odszyfrowany wszystkie możliwe kandydatury (25 dla szyfru Cezara, 311 dla szyfru afinicznego).
		inputText = "files/crypto.txt"
		outputText = "files/decrypt.txt"
	default:
		fmt.Println("Nieobsługiwana operacja.")
		return
	}

	// Execute corresponding cipher operations
	switch cipherType {
	case "caesar":
		CipherOperations(operation, inputText, inputTextHelper, inputKey, outputText, outputKey, CaesarCipher, FindCaesarKey)
	case "affine":
		CipherOperations(operation, inputText, inputTextHelper, inputKey, outputText, outputKey, AffineCipher, FindAffineKey)
	default:
		fmt.Println("Nieobsługiwany typ szyfru.")
	}
}

// Generic cipher function for both Caesar and Affine ciphers
func CipherOperations(operation, inputText, inputTextHelper, inputKey, outputText, outputKey string,
	cipherFunc func(string, int, int, string) string, keyFinder func(string, string) (int, int)) {

	// Read input text
	textLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}
	originalText := strings.Join(textLines, "\n")

	var key1, key2 int

	if operation == "e" || operation == "d" {
		// Read and validate the key from key.txt
		key1, key2 = helpers.ValidateKey(inputKey)
	} else if operation == "j" {
		// Read extra helper text for key finding
		extraLines, err := helpers.GetText(inputTextHelper)
		if err != nil {
			log.Fatalf("Błąd przy odczycie pliku: %v", err)
		}
		extraText := strings.Join(extraLines, "\n")

		// Find key based on ciphertext and helper text
		key1, key2 = keyFinder(originalText, extraText)

		// Save the found key
		helpers.SaveOutput(fmt.Sprintf("%d %d", key1, key2), outputKey)
	}

	// Process encryption or decryption
	processedText := cipherFunc(originalText, key1, key2, operation)

	// Save output
	helpers.SaveOutput(processedText, outputText)
}

// ------------------------------------------------------------------------Caesar Cipher------------------------------------------------------------------------
// CaesarCipher encrypts or decrypts text using the Caesar cipher based on the given flags.
func CaesarCipher(text string, key1, _ int, operation string) string {
	var result strings.Builder

	for _, char := range text {
		// Handle encryption for lowercase letters and uppercase letters
		if operation == "e" {
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift + key1) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' {
				shift := int(char - 'A')
				shift = (shift + key1) % 26
				result.WriteRune(rune(shift + 'A'))
			}
		}

		// Handle decryption for lowercase letters and uppercase letters
		if operation == "d" {
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift - key1 + 26) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' {
				shift := int(char - 'A')
				shift = (shift - key1 + 26) % 26
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



// FindCaesarKey calculates the Caesar cipher key based on the first matching characters in the ciphertext and extra text.
func FindCaesarKey(cryptoText, extraText string) (int, int) {
	// Znajdujemy pierwszy pasujący znak w obu tekstach
	for i := 0; i < len(extraText) && i < len(cryptoText); i++ {
		cipherChar := cryptoText[i]
		plainChar := extraText[i]

		
		if (cipherChar >= 'A' && cipherChar <= 'Z' && plainChar >= 'A' && plainChar <= 'Z') ||
			(cipherChar >= 'a' && cipherChar <= 'z' && plainChar >= 'a' && plainChar <= 'z') {

			// Calculate the key based on the difference between the characters
			key := int(cipherChar - plainChar)

			// Key must be between 0 and 25
			if key < 0 {
				key += 26
			}
			return key, 0
		}
	}

	log.Fatal("Nie udało się znaleźć pasujących znaków do odgadnięcia klucza.")
	return -1, 0
}


// CaesarExplicitCryptAnalysis make analysis of Caesar cipher based on the extra text.
func CaesarExplicitCryptAnalysis(inputText string, inputTextHelper string, outputText string, outputKey string) {
	// Read the entire ciphertext.
	cryptoLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
	}
	cryptoText := strings.Join(cryptoLines, "\n")

	// Read the extra text.
	extraLines, err := helpers.GetText(inputTextHelper)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputTextHelper, err)
	}
	extraText := strings.Join(extraLines, "\n")

	// Check if the input files are not empty.
	if len(cryptoText) == 0 || len(extraText) == 0 {
		log.Fatal("Błąd: Brak danych w plikach wejściowych.")
	}

	// Guess the key based on the first characters of the ciphertext and extra text.
	guessedKey, _ := FindCaesarKey(cryptoText, extraText)

	// Save the found key to a file.
	err = helpers.SaveOutput(strconv.Itoa(guessedKey), outputKey)
	if err != nil {
		log.Fatalf("Błąd przy zapisie klucza: %v", err)
	}

	// Decrypt the ciphertext using the guessed key.
	decryptedText := CaesarCipher(cryptoText, guessedKey, 0, "d")

	// Save thedecrypted text to a file.
	err = helpers.SaveOutput(decryptedText, outputText)
	if err != nil {
		log.Fatalf("Błąd przy zapisie odszyfrowanego tekstu: %v", err)
	}
}

func CaesarCryptAnalysis(inputText string, outputText string) string {
	var result strings.Builder
	return result.String()
}

// ------------------------------------------------------------------------Affine Cipher------------------------------------------------------------------------
// Affine cipher function
func AffineCipher(text string, a, b int, operation string) string {
	var result strings.Builder

	if operation == "e" {
		for _, char := range text {
			if char >= 'a' && char <= 'z' {
				result.WriteRune('a' + rune((a*int(char-'a')+b)%26))
			} else if char >= 'A' && char <= 'Z' {
				result.WriteRune('A' + rune((a*int(char-'A')+b)%26))
			} else {
				result.WriteRune(char)
			}
		}
	} else if operation == "d" {
		invA := modInverse(a, 26)
		for _, char := range text {
			if char >= 'a' && char <= 'z' {
				result.WriteRune('a' + rune((invA*(int(char-'a')-b+26))%26))
			} else if char >= 'A' && char <= 'Z' {
				result.WriteRune('A' + rune((invA*(int(char-'A')-b+26))%26))
			} else {
				result.WriteRune(char)
			}
		}
	}

	return result.String()
}

// Modular inverse function
func modInverse(a, m int) int {
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x
		}
	}
	log.Fatal("Brak modularnej odwrotności dla danego 'a'.")
	return -1
}

// Find key for Affine cipher (placeholder implementation)
func FindAffineKey(cryptoText, extraText string) (int, int) {
	log.Fatal("Funkcja odgadnięcia klucza dla szyfru afinicznego nie jest jeszcze zaimplementowana.")
	return -1, -1
}



func AffineExplicitCryptAnalysis(inputText string, inputTextHelper string, outputText string, outputKey string) string {
	var result strings.Builder
	return result.String()
}

func AffineCryptAnalysis(inputText string, outputText string)  string {
	var result strings.Builder
	return result.String()
}
