package flagfunc

import (
	"fmt"
	"log"
	"strings"

	"vignere/helpers"
)

//Author: Paulina Kimak

// Generic function to handle Vignere cipher
func ExecuteCipher(operation string) {
	switch operation {
	case "p":
		// Prepare the text for encryption and save it to plain.txt
		inputText := "files/org.txt"
		plainText := "files/plain.txt"
		err := CreatePlainFile(inputText, plainText)
		if err != nil {
			fmt.Printf("błąd przy przygotowywaniu tekstu: %v", err)
		}
	case "e":
		// Encode the text from plain.txt using the key from key.txt and saves the result to crypto.txt
		plainText := "files/plain.txt"
		inputKey := "files/key.txt"
		cryptoFile := "files/crypto.txt"
		inputText, err := EncodeVignere(plainText, inputKey, cryptoFile)
		if err != nil {
			log.Printf("nie udało się odczytać poprawnego tekstu %v", err)
		}
		fmt.Println("Odczytany tekst: ", inputText)
	case "d":
		// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
		inputText := "files/crypto.txt"
		inputKey := "files/key.txt"
		outputText := "files/decrypt.txt"
		decryptVigenereSimple(inputText, inputKey, outputText)
		// DecodeVignere(inputText, inputKey, outputText) 
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		inputText := "files/crypto.txt"
		outputText := "files/decrypt.txt"
		CryptAnalysisVignere(inputText, outputText)

	default:
		fmt.Println("Nieobsługiwana operacja.")
		return
	}	

}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile(inputFile string, outputFile string) error {
	plainText, err := helpers.PrepareText(inputFile)
	if err != nil {
		log.Printf("błąd przy czyszczeniu tekstu: %v", err)
		return fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
	}

	err = helpers.SaveOutput(plainText, outputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
}


// Function to encode text using Vigenere cipher with numeric key shifts
func VigenereEncode2(plainText, inputKey, outputText string) (string, error) {
	keyShifts ,err := helpers.GetKey(inputKey)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać poprawnego klucza")
	}
	fmt.Println("Skonwersowany Klucz  do liczb: ", keyShifts)

	var encodedText []rune

	for i, char := range plainText {
		shift := keyShifts[i % len(keyShifts)] 
		encodedChar := 'a' + (char-'a'+rune(shift))%26
		encodedText = append(encodedText, encodedChar)
	}

	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(string(encodedText), outputText)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(encodedText), nil
}

// encryptVigenere encrypts the given text using the Vigenère cipher.
func EncodeVignere(plainFile, keyFile, cryptoFile string) (string, error) {
	
	fmt.Printf("Plik Klucz: %s\n", keyFile)
	fmt.Printf("Plik Tekst: %s\n", plainFile)
	fmt.Printf("cryptoFile: %s\n", cryptoFile)
	
	plainText, err := helpers.GetPlainText(plainFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać plain tekstu")
	}

	key, err := helpers.GetKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}

	fmt.Printf("Klucz: %s\n", key)
	fmt.Printf("Plain Tekst: %s\n", plainText)
	
	if len(plainText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, or key cannot be empty")
	}
	
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	alphabetLen := len(alphabet)
	var result []rune

	for i, char := range plainText {
		index := strings.IndexRune(alphabet, char)
		keyIndex := strings.IndexRune(alphabet, rune(key[i % len(key)]))

		encryptedIndex := (index + keyIndex) % alphabetLen
		result = append(result, rune(alphabet[encryptedIndex]))
	}

	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(string(result), cryptoFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// decryptVigenereSimple decrypts the given cryptoFile using the Vigenère cipher with the provided key.
func decryptVigenereSimple(cryptoFile, key, encodedFile string) (string, error) {
	if len(plainText) == 0 || len(key) == 0 || len(alphabet) == 0 {
		return "", fmt.Errorf("input plainText, key, or alphabet cannot be empty")
	}

	plainText = strings.ToLower(plainText)
	key = strings.ToLower(key)
	keyLength := len(key)
	alphabetLength := len(alphabet)
	var result []rune

	for i, char := range plainText {
		index := strings.IndexRune(alphabet, char)
		if index == -1 {
			result = append(result, char) // Keep non-alphabet characters unchanged
			continue
		}

		keyIndex := strings.IndexRune(alphabet, rune(key[i%keyLength]))
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		decryptedIndex := (index - keyIndex + alphabetLength) % alphabetLength
		result = append(result, rune(alphabet[decryptedIndex]))
	}

	return string(result), nil
}

// getReversedKey generates a reversed key for decrypting using the Wikipedia formula: K2(i) = [26 – K(i)] mod 26
func getReversedKey(key, alphabet string) (string, error) {
	if len(key) == 0 || len(alphabet) == 0 {
		return "", fmt.Errorf("key or alphabet cannot be empty")
	}

	alphabetLength := len(alphabet)
	var reversedKey []rune

	for _, char := range key {
		keyIndex := strings.IndexRune(alphabet, char)
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		reversedIndex := (alphabetLength - keyIndex) % alphabetLength
		reversedKey = append(reversedKey, rune(alphabet[reversedIndex]))
	}

	return string(reversedKey), nil
}

// decryptReversedKey decrypts the plainText using the reversed key.
func DecodeVignere(cryptoText, reversedKey, alphabet string) (string, error) {
	return EncodeVignere(cryptoText, reversedKey, alphabet)
}


//------------------------------------------------------------Kryptoanaliza------------------------------------------------------------
func CryptAnalysisVignere(cryptoText, outputFile string) error{
	// Save the decrypted text to decrypt.txt
	err := helpers.SaveOutput(cryptoText, outputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
}