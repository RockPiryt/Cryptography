package cryptofunc

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"caesaraffineciphers/helpers"
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
	var inputText, ipputTextHelper, inputKey, outputText,outputkey string

	switch operation {
	case "e":
		// Program szyfrujący czyta tekst jawny i klucz i zapisuje tekst zaszyfrowany. Jeśli klucz jest nieprawidłowy, zgłasza jedynie błąd.
		inputText = "files/plain.txt"
		inputKey = "files/key.txt"
		outputText = "files/crypto.txt"
		CaesarOperations("e", inputText, inputKey, outputText) 
	case "d":
		// Program odszyfrowujący czyta tekst zaszyfrowany i klucz i zapisuje tekst jawny. Jeśli klucz jest nieprawidłowy, zgłasza błąd. 
		// Dla szyfru afinicznego częścią zadania jest znalezienie odwrotności dla liczby a podanej jako część klucza – 
		// nie można zakładać, że program odszyfrowujący otrzymuje tę odwrotność.
		inputText = "files/crypto.txt"
		inputKey = "files/key.txt"
		outputText = "files/decrypt.txt"
		CaesarOperations("d",inputText, inputKey, outputText) 
	case "j":
		// Program łamiący szyfr z pomocą tekstu jawnego czyta tekst zaszyfrowany, tekst pomocniczy i zapisuje znaleziony klucz i odszyfrowany tekst. 
		// Jeśli niemożliwe jest znalezienie klucza, zgłasza sygnał błędu.
		inputText = "files/crypto.txt"
		ipputTextHelper = "files/extra.txt"
		outputText = "files/decrypt.txt"
		outputkey = "files/key-found.txt"
		CaesarExplicitCryptAnalysis(inputText, inputTextHelper, outputText, outputKey)
	case "k":
		//Program łamiący szyfr bez pomocy tekstu jawnego czyta jedynie tekst zaszyfrowany i zapisuje jako tekst odszyfrowany wszystkie możliwe kandydatury (25 dla szyfru Cezara, 311 dla szyfru afinicznego).
		inputText = "files/crypto.txt"
		outputText = "files/decrypt.txt"
		CaesarCryptAnalysis(inputText, outputText) 
	default:
		fmt.Println("Nieobsługiwana operacja dla Cezara.")
		return
	}

}

func CaesarOperations(operation string, inputText string, inputKey string, outputText string) {
	// Read text from input file
	textLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")

	// Read and validate the key
	key, _ := helpers.ValidateKey(inputKey)

	// Encrypt or decrypt
	processedText := CaesarCipher(originalText, key, operation)

	// Save the result
	helpers.SaveOutput(processedText, outputText)
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
	guessedKey := FindCaesarKey(cryptoText, extraText)

	// Save the found key to a file.
	err = helpers.SaveOutput(strconv.Itoa(guessedKey), outputKey)
	if err != nil {
		log.Fatalf("Błąd przy zapisie klucza: %v", err)
	}

	// Decrypt the ciphertext using the guessed key.
	decryptedText := CaesarCipher(cryptoText, guessedKey, "d")

	// Save thedecrypted text to a file.
	err = helpers.SaveOutput(decryptedText, outputText)
	if err != nil {
		log.Fatalf("Błąd przy zapisie odszyfrowanego tekstu: %v", err)
	}
}

// FindCaesarKey calculates the Caesar cipher key based on the first matching characters in the ciphertext and extra text.
func FindCaesarKey(cryptoText, extraText string) int {
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
			return key
		}
	}

	log.Fatal("Nie udało się znaleźć pasujących znaków do odgadnięcia klucza.")
	return -1 
}


func CaesarCryptAnalysis(inputText string, outputText string) string {
	var result strings.Builder
	return result.String()
}

// ------------------------------------------------------------------------Affine Cipher------------------------------------------------------------------------
func AffineCipher(text string, a int, b int,  operation string) string {
	return "a"
}
func AffineExecute(operation string) {
	var inputFile, outputFile string

	switch operation {
	case "e":
		inputText = "files/plain.txt"
		inputKey = "files/key.txt"
		outputText = "files/crypto.txt"
		AffineOperations("e", inputText, inputKey, outputText) 
	case "d":
		inputText = "files/crypto.txt"
		inputKey = "files/key.txt"
		outputText = "files/decrypt.txt"
		AffineOperations("d",inputText, inputKey, outputText) 
	case "j":
		inputText = "files/crypto.txt"
		ipputTextHelper = "files/extra.txt"
		outputText = "files/decrypt.txt"
		outputkey = "files/key-found.txt"
		AffineExplicitCryptAnalysis(inputText, inputTextHelper, outputText, outputKey)
	case "k":
		inputText = "files/crypto.txt"
		outputText = "files/decrypt.txt"
		AffineCryptAnalysis(inputText, outputText) 
	default:
		fmt.Println("Nieobsługiwana operacja dla Cezara.")
		return
	}


	// Odczytaj tekst wejściowy
	textLines, err := helpers.GetText(inputFile)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")

	// Odczytaj i zweryfikuj klucze (a, b) dla szyfru afinicznego
	a, b := helpers.ValidateKey(inputKey)

	// Wykonaj szyfrowanie lub deszyfrowanie
	processedText := AffineCipher(originalText, a, b, operation)

	// Zapisz wynik do pliku
	helpers.SaveOutput(processedText, outputFile)
}


func AffineOperations(operation string, inputText string, inputKey string, outputText string) {
	// Read text from input file
	textLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	originalText := strings.Join(textLines, "\n")

	// Read and validate the a and b keys for the Affine cipher
	a, b := helpers.ValidateKey(inputKey)

	// Encrypt or decrypt
	processedText := CaesarCipher(originalText, a, b, operation)

	// Save the result
	helpers.SaveOutput(processedText, outputText)
}
func AffineExplicitCryptAnalysis(inputText string, inputTextHelper string, outputText string, outputKey string) string {
	var result strings.Builder
	return result.String()
}

func AffineCryptAnalysis(inputText string, outputText string)  string {
	var result strings.Builder
	return result.String()
}
