package cryptofunc

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

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

// ------------------------------------------------------------------------General functions------------------------------------------------------------------------
// Generic function to handle both Caesar and Affine ciphers
func ExecuteCipher(cipherType string, operation string) {
	params := CipherParams{
		Operation:   operation,
		CipherType:  cipherType,
	}

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
		// Program łamiący szyfr z pomocą tekstu jawnego czyta tekst zaszyfrowany, tekst pomocniczy. Następnie zapisuje znaleziony klucz i odszyfrowany tekst. 
		// Jeśli niemożliwe jest znalezienie klucza, zgłasza sygnał błędu.
		params.InputText = "files/crypto.txt"
		params.InputTextHelper = "files/extra.txt"
		params.OutputText = "files/decrypt.txt"
		params.OutputKey = "files/key-found.txt"
		// If Caesar explicit cryptanalysis (-c -j), call specialized function
		if cipherType == "caesar" {
			CaesarExplicitCryptAnalysis(params.InputText, params.InputTextHelper, params.OutputText, params.OutputKey)
			return
		}
	case "k":
		//Program łamiący szyfr bez pomocy tekstu jawnego czyta jedynie tekst zaszyfrowany i zapisuje jako tekst odszyfrowany wszystkie możliwe kandydatury (25 dla szyfru Cezara, 311 dla szyfru afinicznego).
		params.InputText = "files/crypto.txt"
		params.OutputText = "files/decrypt.txt"
		// If Caesar explicit cryptanalysis (-c -k), call specialized function
		if cipherType == "caesar" {
			CaesarCryptAnalysis(params.InputText, params.OutputText)
			return
		}
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

	var c, a int

	switch params.Operation {
	case "e", "d":
		// Validate the key.
		c, a = helpers.ValidateKey(params.InputKey, params.CipherType)
	case "j":
		// Read the extra text.
		extraTextLines, err := helpers.GetText(params.InputTextHelper)
		if err != nil {
			log.Fatalf("Błąd odczytu pliku pomocniczego: %v", err)
		}
		extraText := strings.Join(extraTextLines, "\n")

		// Find the key based on the extra text.
		c, a = params.KeyFinder(originalText, extraText)

		// Save the key to a file.
		helpers.SaveOutput(fmt.Sprintf("%d %d", c, a), params.OutputKey)
	default:
		log.Fatalf("Nieznana operacja: %s", params.Operation)
	}

	// Execute the cipher function.
	resultText := params.CipherFunc(originalText, a, c, params.Operation)

	// Save the result to a file.
	helpers.SaveOutput(resultText, params.OutputText)
}
// ------------------------------------------------------------------------Caesar Cipher------------------------------------------------------------------------
// CaesarCipher encrypts or decrypts text using the Caesar cipher based on the given flags.
func CaesarCipher(text string, c, _ int, operation string) string {
	var result strings.Builder

	for _, char := range text {
		// Handle encryption for lowercase letters and uppercase letters
		if operation == "e" {
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift + c) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' {
				shift := int(char - 'A')
				shift = (shift + c) % 26
				result.WriteRune(rune(shift + 'A'))
			}
		}

		// Handle decryption for lowercase letters and uppercase letters
		if operation == "d" {
			if char >= 'a' && char <= 'z' {
				shift := int(char - 'a')
				shift = (shift - c + 26) % 26
				result.WriteRune(rune(shift + 'a'))
			} else if char >= 'A' && char <= 'Z' {
				shift := int(char - 'A')
				shift = (shift - c + 26) % 26
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
func CaesarExplicitCryptAnalysis(inputText, inputTextHelper, outputText, outputKey string) {
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

	if len(cryptoText) == 0 || len(extraText) == 0 {
		log.Fatal("Błąd: Brak danych w plikach wejściowych.")
	}
	
	// Guess the Caesar key by comparing characters
	guessedKey := (int(cryptoText[0]) - int(extraText[0]) + 26) % 26
	guessedKeyString := strconv.Itoa(guessedKey)

	// Save the guessed key
	helpers.SaveOutput(guessedKeyString, outputKey)

	// Decrypt using the guessed key
	decryptedText := CaesarCipher(cryptoText, guessedKey, -1, "d")

	helpers.SaveOutput(decryptedText, outputText)
}

func CaesarCryptAnalysis(inputText string, outputText string) {
	// Read the entire ciphertext.
	cryptoLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
	}
	cryptoText := strings.Join(cryptoLines, "\n")

	if len(cryptoText) == 0 {
		log.Fatal("Błąd: Brak danych w pliku wejściowym.")
	}

	var result strings.Builder

	// Test all possible keys (0-25)
	for key := 1; key <= 25; key++ {
		decryptedText := CaesarCipher(cryptoText, key, -1, "d")
		result.WriteString(decryptedText + "\n") 
	}

	helpers.SaveOutput(result.String(), outputText)
}

// ------------------------------------------------------------------------Affine Cipher------------------------------------------------------------------------
// Affine cipher function
func AffineCipher(text string, a, c int, operation string) string {
	result := ""
	m := 26

	fmt.Printf("AffineCipher: a=%d, c=%d, operation=%s\n", a, c, operation)
	fmt.Printf("Input text: %s\n", text)

	// If decrypting, calculate the modular inverse of 'a'
	aInv := 0
	if operation == "d" {
		aInv = helpers.ModInverseExtended(a, m)
	}

	for _, char := range text {
		if unicode.IsLetter(char) && unicode.Is(unicode.Latin, char) { // Skip Polish characters and others
			isUpper := unicode.IsUpper(char)
			var base rune
			if isUpper {
				base = 'A'
			} else {
				base = 'a'
			}

			x := int(char - base)

			if operation == "e" { // Encrypt
				y := (a*x + c) % m
				result += string(rune(y) + base)
			} else if operation == "d" { // Decrypt
				y := (aInv * (x - c + m)) % m
				result += string(rune(y) + base)
			}
		} else {
			result += string(char) // Skip other characters
		}
	}

	return result
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
