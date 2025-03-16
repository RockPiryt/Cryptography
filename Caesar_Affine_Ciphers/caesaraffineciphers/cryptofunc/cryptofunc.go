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
	InputText       string
	InputTextHelper string
	InputKey        string
	OutputText      string
	OutputKey       string
	CipherType      string
	CipherFunc      func(string, int, int, string) (string, error)
	KeyFinder       func(string, string) (int, int)
}

// ------------------------------------------------------------------------General functions------------------------------------------------------------------------
// Generic function to handle both Caesar and Affine ciphers
func ExecuteCipher(cipherType string, operation string) {
	params := CipherParams{
		Operation:  operation,
		CipherType: cipherType,
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
		} else if cipherType == "affine" {
			AffineExplicitCryptAnalysis(params.InputText, params.InputTextHelper, params.OutputText, params.OutputKey)
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
		} else if cipherType == "affine" {
			AffineCryptAnalysis(params.InputText, params.OutputText)
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
		//fmt.Printf("Klucz: %d %d\n", c, a)
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

	// Execute the cipher function based on the cipher type
	var resultText string
	var cipherErr error

	switch params.CipherType {
	case "caesar":
		// If it's Caesar, use CaesarCipher
		resultText, cipherErr = CaesarCipher(originalText, a, c, params.Operation)
	case "affine":
		// If it's Affine, use AffineCipher
		resultText, cipherErr = AffineCipher(originalText, a, c, params.Operation)
	default:
		log.Fatalf("Nieobsługiwany typ szyfru: %s", params.CipherType)
	}

	if cipherErr != nil {
		log.Fatalf("Błąd szyfrowania: %v", cipherErr)
	}

	// Save the result to a file.
	helpers.SaveOutput(resultText, params.OutputText)

}

// ------------------------------------------------------------------------Caesar Cipher------------------------------------------------------------------------
// CaesarCipher encrypts or decrypts text using the Caesar cipher based on the given flags.
func CaesarCipher(text string, _, c int, operation string) (string, error) {
	var result strings.Builder
	//fmt.Printf("CaesarCipher: c=%d, operation=%s\n", c, operation)

	if operation != "e" && operation != "d" {
		return "", fmt.Errorf("nieprawidłowa operacja: %s", operation)
	}

	if len(text) == 0 {
		return "", fmt.Errorf("tekst wejściowy jest pusty")
	}

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

	return result.String(), nil
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
			return key, 1
		}
	}

	log.Fatal("Nie udało się znaleźć pasujących znaków do odgadnięcia klucza.")
	return -1, 1
}

// CaesarExplicitCryptAnalysis make analysis of Caesar cipher based on the extra text.
func CaesarExplicitCryptAnalysis(inputText, inputTextHelper, outputText, outputKey string) {
	// Read the entire ciphertext.
	cryptoLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
	}
	cryptoText := strings.Join(cryptoLines, "\n")
	//fmt.Printf("Oczekiwany zaszyfrowany tekst: %s\n", cryptoText)

	// Read the extra text.
	extraLines, err := helpers.GetText(inputTextHelper)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputTextHelper, err)
	}
	extraText := strings.Join(extraLines, "\n")
	//fmt.Printf("Oczekiwany jawny tekst pomocniczy: %s\n", extraText)

	if len(cryptoText) == 0 || len(extraText) == 0 {
		log.Fatal("Błąd: Brak danych w plikach wejściowych.")
	}

	// Guess the Caesar key by comparing characters
	guessedKey := (int(cryptoText[0]) - int(extraText[0]) + 26) % 26
	guessedKeyString := strconv.Itoa(guessedKey)

	// Save the guessed key
	helpers.SaveOutput(guessedKeyString + " 1", outputKey)

	// Decrypt using the guessed key
	decryptedText, _ := CaesarCipher(cryptoText, -1, guessedKey, "d")

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
		decryptedText, _ := CaesarCipher(cryptoText, -1, key, "d")
		result.WriteString(decryptedText + "\n")
	}

	helpers.SaveOutput(result.String(), outputText)
}

// ------------------------------------------------------------------------Affine Cipher------------------------------------------------------------------------
// Affine cipher function
func AffineCipher(text string, a, c int, operation string) (string, error) {
	m := 26
	var result string

	//fmt.Printf("AffineCipher: a=%d, c=%d, operation=%s\n", a, c, operation)
	//fmt.Printf("Input text: %s\n", text)

	// Calculate the modular inverse of 'a' if decrypting
	aInv := 0
	if operation == "d" {
		var err error
		aInv, err = helpers.ModInverseExtended(a, m)
		if err != nil {
			return "", fmt.Errorf("nie można odszyfrować: %v", err)
		}
	}

	// Process each character
	for _, char := range text {
		if unicode.IsLetter(char) && unicode.Is(unicode.Latin, char) { // Only process Latin letters
			isUpper := unicode.IsUpper(char)
			base := 'a'
			if isUpper {
				base = 'A'
			}

			x := int(char - base) // Convert to numerical value (0-25)

			if operation == "e" { // Encrypt
				y := (a*x + c) % m
				result += string(rune(y) + base)
			} else if operation == "d" { // Decrypt
				y := (aInv * ((x - c + m) % m)) % m
				result += string(rune(y) + base)
			}
		} else {
			result += string(char) // Leave non-letter characters unchanged
		}
	}

	return result, nil
}

// FindAffineKey finds the affine cipher key based on the first matching characters in the ciphertext and extra text.
func FindAffineKey(cryptoText, extraText string) (int, int) {

	// Find the first matching pair of characters in both texts
	for i := 0; i < len(cryptoText)-1; i++ {
		// Convert characters to zero-based indices (0-25)
		x1 := int(unicode.ToUpper(rune(extraText[i])) - 'A')
		y1 := int(unicode.ToUpper(rune(cryptoText[i])) - 'A')

		x2 := int(unicode.ToUpper(rune(extraText[i+1])) - 'A')
		y2 := int(unicode.ToUpper(rune(cryptoText[i+1])) - 'A')

		//fmt.Printf("x1: %d, y1: %d, x2: %d, y2: %d\n", x1, y1, x2, y2)

		a, c, err := solveAffineSystemSr(x1, x2, y1, y2)
		if err != nil {
			log.Fatal("Nie udało się znaleźć poprawnych wartości a i c:", err)
		}

		return a, c
	}

	log.Fatal("Nie udało się znaleźć odpowiednich par do rozwiązania układu równań.")
	return -1, -1
}

// solveAffineSystemSr solves the system of two affine equations to find the key for the affine cipher.
func solveAffineSystemSr(x1, x2, y1, y2 int) (int, int, error) {
	m := 26
	//fmt.Printf("solveAffineSystem: x1=%d, x2=%d, y1=%d, y2=%d\n", x1, x2, y1, y2)
	// Calculate the differences.
	deltaY := (y1 - y2 + m) % m // ex.15 - 16 = -1 -> mod 26 -> 25
	deltaX := (x1 - x2 + m) % m //ex 8 - 5 = 3
	//fmt.Printf("deltaY: %d, deltaX: %d\n", deltaY, deltaX)

	// Find the modular inverse of deltaX (3 mod 26)
	invDeltaX, err := helpers.ModInverseExtended(deltaX, m)
	//fmt.Printf("invDeltaX: %d\n", invDeltaX)
	if err != nil {
		return -1, -1, fmt.Errorf("nie można znaleźć odwrotności modularnej: %v", err)
	}

	x := (deltaY * invDeltaX) % m // x = (25 * 9) mod 26
	//fmt.Printf("Obliczone a: %d\n", x)

	// Calculate y from the first equation: y = 15 - x * 8 mod 26
	y := (y1 - x*x1) % m
	if y < 0 {
		y += m
	}
	//.Printf("Obliczone c: %d\n", y)

	return x, y, nil
}

// Function to break the Affine cipher using known plaintext (extra text)
func AffineExplicitCryptAnalysis(inputText, inputTextHelper, outputText, outputKey string) {
	// Read the entire ciphertext.
	cryptoLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
	}
	cryptoText := strings.Join(cryptoLines, "\n")
	//fmt.Printf("Oczekiwany zaszyfrowany cryptoText: %s\n", cryptoText)

	// Read the extra text (known plaintext).
	extraLines, err := helpers.GetText(inputTextHelper)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputTextHelper, err)
	}
	extraText := strings.Join(extraLines, "\n")
	//fmt.Printf("Oczekiwany jawny extraText: %s\n", extraText)

	if len(cryptoText) == 0 || len(extraText) == 0 {
		log.Fatal("Błąd: Brak danych w plikach wejściowych.")
	}

	// Find the affine cipher key based on the known plaintext
	a, c := FindAffineKey(cryptoText, extraText)
	//fmt.Printf("Znaleziony klucz: a=%d, c=%d\n", a, c)

	// Save the key (a, c) to the output key file
	helpers.SaveOutput(fmt.Sprintf("%d %d", c, a), outputKey)

	// Decrypt the ciphertext using the found key
	decryptedText, err := AffineCipher(cryptoText, a, c, "d")
	if err != nil {
		log.Fatalf("Błąd odszyfrowania: %v", err)
	}
	//fmt.Printf("Odszyfrowany tekst: %s\n", decryptedText)

	// Save the decrypted text to the output file
	helpers.SaveOutput(decryptedText, outputText)
}

// AffineCryptAnalysis tests all possible keys (a, c) for the affine cipher and saves all decrypted texts.
func AffineCryptAnalysis(inputText string, outputText string) {
	m := 26
	var result strings.Builder

	// Read the entire ciphertext
	cryptoLines, err := helpers.GetText(inputText)
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku %s: %v", inputText, err)
	}
	cryptoText := strings.Join(cryptoLines, "\n")
	//fmt.Printf("Oczekiwany zaszyfrowany tekst: %s\n", cryptoText)

	// Check if the ciphertext is empty
	if len(cryptoText) == 0 {
		log.Fatal("Błąd: Brak danych w pliku wejściowym.")
	}

	// Check all possible combinations of a (1-25) and c (0-25) for the affine cipher
	for a := 0; a < m; a++ {
		// Check if 'a' is coprime with 26 (i.e., gcd(a, 26) == 1)
		gcd, _, _ := helpers.ExtendedGCD(a, m)
		if gcd == 1 {
			for c := 0; c < m; c++ {
				// Skip the key (1, 0), as it does not alter the text
				if a == 1 && c == 0 {
					continue
				}

				// Decrypt the text using the affine cipher with key (a, c)
				decryptedText, err := AffineCipher(cryptoText, a, c, "d")
				if err != nil {
					log.Fatalf("Błąd odszyfrowania: %v", err)
				}
				//fmt.Printf("Próba: a=%d, c=%d\n odszyfrowany tekst: %s\n\n", a, c, decryptedText)

				// Save the decrypted text
				result.WriteString(fmt.Sprintf("%s\n", decryptedText))
			}
		}
	}

	helpers.SaveOutput(result.String(), outputText)
}
