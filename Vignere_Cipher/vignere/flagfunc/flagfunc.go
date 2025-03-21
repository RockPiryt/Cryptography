//Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"strings"

	"vignere/helpers"
)

// Generic function to handle Vignere cipher
func ExecuteCipher(operation string) {
	switch operation {
	case "p":
		// Prepare the text for encryption and save it to plain.txt
		orgFile := "files/org.txt"
		plainFile := "files/plain.txt"
		err := CreatePlainFile(orgFile, plainFile)
		if err != nil {
			fmt.Printf("błąd przy przygotowywaniu tekstu: %v", err)
		}
	case "e":
		// Encode the text from plain.txt using the key from key.txt and saves the result to crypto.txt
		plainFile := "files/plain.txt"
		keyFile := "files/key.txt"
		cryptoFile := "files/crypto.txt"
		encodedText, err := EncodeVignere(plainFile, keyFile, cryptoFile)
		if err != nil {
			log.Printf("nie udało się zaszyfrować tekstu %v", err)
		}
		fmt.Println("Zaszyfrowany tekst: ", encodedText)
	case "d":
		// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
		cryptoFile := "files/crypto.txt"
		keyFile := "files/key.txt"
		decryptedFile := "files/decrypt.txt"
		decodedText, err := DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile)
		if err != nil {
			log.Printf("nie udało się odszyfrować tekstu %v", err)
		}
		fmt.Println("Odszyfrowany tekst: ", decodedText)
		// DecodeVignere(cryptoFile, keyFile, decryptedFile) 
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		cryptoFile := "files/crypto.txt"
		decryptedFile := "files/decrypt.txt"
		CryptAnalysisVignere(cryptoFile, decryptedFile)   

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

// Function to encrypt the plainText using the Vigenère cipher with the provided key.
func EncodeVignere(plainFile, keyFile, cryptoFile string) (string, error) {
	plainText, err := helpers.GetPreparedText(plainFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać plain tekstu")
	}

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}

	fmt.Printf("Klucz: %s\n", key)
	fmt.Printf("Plain Tekst: %s\n", plainText)
	
	if len(plainText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, or key cannot be empty")
	}
	
	var result []rune

	for i, char := range plainText {
		index := strings.IndexRune(helpers.Alphabet, char)
		keyIndex := strings.IndexRune(helpers.Alphabet, rune(key[i % len(key)]))

		encryptedIndex := (index + keyIndex) % helpers.AlphabetLen
		result = append(result, rune(helpers.Alphabet[encryptedIndex]))
	}

	// Save the decrypted text to crypto.txt
	err = helpers.SaveOutput(string(result), cryptoFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// decryptVigenereSimple decrypts the given cryptoFile using the Vigenère cipher with the provided key.
func DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile string) (string, error) {
	cryptoText, err := helpers.GetPreparedText(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać crypto tekstu")
	}
	fmt.Printf("Crypto Tekst: %s\n", cryptoText)

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}
	fmt.Printf("Klucz: %s\n", key)
	if len(cryptoText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, key, or helpers.Alphabet cannot be empty")
	}
	keyLength := len(key)
	var result []rune

	for i, char := range cryptoText {
		index := strings.IndexRune(helpers.Alphabet, char)
		if index == -1 {
			result = append(result, char) // Keep non-helpers.Alphabet characters unchanged
			continue
		}

		keyIndex := strings.IndexRune(helpers.Alphabet, rune(key[i%keyLength]))
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		decryptedIndex := (index - keyIndex + helpers.AlphabetLen) % helpers.AlphabetLen
		result = append(result, rune(helpers.Alphabet[decryptedIndex]))
	}

	fmt.Printf("Odszyfrowany tekst: %s\n", string(result))

	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(string(result), decryptedFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// // getReversedKey generates a reversed key for decrypting using the Wikipedia formula: K2(i) = [26 – K(i)] mod 26
// func getReversedKey(key string) (string, error) {
// 	if len(key) == 0  {
// 		return "", fmt.Errorf("key cannot be empty")
// 	}

// 	var reversedKey []rune

// 	for _, char := range key {
// 		keyIndex := strings.IndexRune(helpers.Alphabet, char)
// 		if keyIndex == -1 {
// 			return "", fmt.Errorf("invalid character in key")
// 		}

// 		reversedIndex := (alphabetLength - keyIndex) % alphabetLength
// 		reversedKey = append(reversedKey, rune(helpers.Alphabet[reversedIndex]))
// 	}

// 	return string(reversedKey), nil
// }

// // decryptReversedKey decrypts the plainText using the reversed key.
// func DecodeVignere(cryptoText, reversedKey string) (string, error) {
// 	return EncodeVignere(cryptoText, reversedKey)
// }


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