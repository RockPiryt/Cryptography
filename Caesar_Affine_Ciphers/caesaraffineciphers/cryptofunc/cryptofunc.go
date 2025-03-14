package cryptofunc

import (
	"fmt"
	"log"
	"strings"

	"caesaraffineciphers/helpers"
)

//Author: Paulina Kimak

// Struct to store cipher parameters.
type CipherParams struct {
	Operation       string
	InputText      string
	InputTextHelper string
	InputKey       string
	OutputText     string
	OutputKey      string
	CipherType     string
	CipherFunc     func(string, int, int, string) string
	KeyFinder      func(string, string) (int, int)
}

// Generic function to handle both Caesar and Affine ciphers
func ExecuteCipher(cipherType string, operation string) {
	params := CipherParams{
		Operation:   operation,
		CipherType:  cipherType,
	}

	// Define file paths based on operation
	switch operation {
	case "e":
		// Program szyfrujący czyta tekst jawny i klucz i zapisuje tekst zaszyfrowany. Jeśli klucz jest nieprawidłowy, zgłasza jedynie błąd.
		params.InputText = "files/plain.txt"
		params.InputKey = "files/key.txt"
		params.OutputText = "files/crypto.txt"
	case "d":
		// Program odszyfrowujący czyta tekst zaszyfrowany i klucz i zapisuje tekst jawny. Jeśli klucz jest nieprawidłowy, zgłasza błąd. 
		// Dla szyfru afinicznego częścią zadania jest znalezienie odwrotności dla liczby a podanej jako część klucza – 
		// nie można zakładać, że program odszyfrowujący otrzymuje tę odwrotność.
		params.InputText = "files/crypto.txt"
		params.InputKey = "files/key.txt"
		params.OutputText = "files/decrypt.txt"
	case "j":
		// Program łamiący szyfr z pomocą tekstu jawnego czyta tekst zaszyfrowany, tekst pomocniczy i zapisuje znaleziony klucz i odszyfrowany tekst. 
		// Jeśli niemożliwe jest znalezienie klucza, zgłasza sygnał błędu.
		params.InputText = "files/crypto.txt"
		params.InputTextHelper = "files/extra.txt"
		params.OutputText = "files/decrypt.txt"
		params.OutputKey = "files/key-found.txt"
	case "k":
		//Program łamiący szyfr bez pomocy tekstu jawnego czyta jedynie tekst zaszyfrowany i zapisuje jako tekst odszyfrowany wszystkie możliwe kandydatury (25 dla szyfru Cezara, 311 dla szyfru afinicznego).
		params.InputText = "files/crypto.txt"
		params.OutputText = "files/decrypt.txt"
	default:
		fmt.Println("Nieobsługiwana operacja.")
		return
	}

	// Assign cipher-specific functions
	switch cipherType {
	case "caesar":
		params.CipherFunc = CaesarCipher
		params.KeyFinder = FindCaesarKey
	case "affine":
		params.CipherFunc = AffineCipher
		params.KeyFinder = FindAffineKey
	default:
		fmt.Println("Nieobsługiwany typ szyfru.")
		return
	}

	// Execute the cipher operation
	CipherOperations(params)
}

// Generic cipher function for both Caesar and Affine ciphers
func CipherOperations(params CipherParams) {
	// Read the input text.
	textLines, err := helpers.GetText(params.InputText)
	if err != nil {
		log.Fatalf("Błąd odczytu pliku: %v", err)
	}
	originalText := strings.Join(textLines, "\n")

	var key1, key2 int

	switch params.Operation {
	case "e", "d":
		// Validate the key.
		key1, key2 = helpers.ValidateKey(params.InputKey, params.CipherType)
	case "j":
		// Read the extra text.
		extraTextLines, err := helpers.GetText(params.InputTextHelper)
		if err != nil {
			log.Fatalf("Błąd odczytu pliku pomocniczego: %v", err)
		}
		extraText := strings.Join(extraTextLines, "\n")

		// Find the key based on the extra text.
		key1, key2 = params.KeyFinder(originalText, extraText)

		// Save the key to a file.
		helpers.SaveOutput(fmt.Sprintf("%d %d", key1, key2), params.OutputKey)
	default:
		log.Fatalf("Nieznana operacja: %s", params.Operation)
	}

	// Execute the cipher function.
	resultText := params.CipherFunc(originalText, key1, key2, params.Operation)

	// Save the result to a file.
	helpers.SaveOutput(resultText, params.OutputText)
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


// // CaesarExplicitCryptAnalysis make analysis of Caesar cipher based on the extra text.
// func CaesarExplicitCryptAnalysis(inputText string, inputTextHelper string, outputText string, outputKey string) {
// 	// Read the entire ciphertext.
// 	cryptoLines, err := helpers.GetText(inputText)
// 	if err != nil {
// 		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
// 	}
// 	cryptoText := strings.Join(cryptoLines, "\n")

// 	// Read the extra text.
// 	extraLines, err := helpers.GetText(inputTextHelper)
// 	if err != nil {
// 		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputTextHelper, err)
// 	}
// 	extraText := strings.Join(extraLines, "\n")

// 	// Check if the input files are not empty.
// 	if len(cryptoText) == 0 || len(extraText) == 0 {
// 		log.Fatal("Błąd: Brak danych w plikach wejściowych.")
// 	}

// 	// Guess the key based on the first characters of the ciphertext and extra text.
// 	_, err = helpers.SaveOutput(guessedKeyString, outputKey)
// 	if err != nil {
// 		log.Fatalf("Błąd przy zapisie klucza: %v", err)
// 	}

// 	// Save thedecrypted text to a file.
// 	err = helpers.SaveOutput(decryptedText, outputText)
// 	if err != nil {
// 		log.Fatalf("Błąd przy zapisie odszyfrowanego tekstu: %v", err)
// 	}
// }

// func CaesarCryptAnalysis(inputText string, outputText string) string {
// 	var result strings.Builder
// 	return result.String()
// }

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
