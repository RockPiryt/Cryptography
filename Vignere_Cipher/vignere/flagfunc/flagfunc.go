package flagfunc

import (
	"fmt"
	"log"

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
		outputText := "files/crypto.txt"
		inputText, err := VigenereEncode(plainText, inputKey, outputText)
		if err != nil {
			log.Printf("nie udało się odczytać poprawnego tekstu %v", err)
		}
		fmt.Println("Odczytany tekst: ", inputText)
	case "d":
		// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
		inputText := "files/crypto.txt"
		inputKey := "files/key.txt"
		outputText := "files/decrypt.txt"
		DecodeText(inputText, inputKey, outputText) 
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		inputText := "files/crypto.txt"
		outputText := "files/decrypt.txt"
		VignereCryptAnalysis(inputText, outputText)

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
func VigenereEncode(plainText, inputKey, outputText string) (string, error) {
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

func VignereCryptAnalysis(cryptoText, outputFile string) error{
	// Save the decrypted text to decrypt.txt
	err := helpers.SaveOutput(cryptoText, outputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
}
func DecodeText(inputText, inputKey, outputText string) {
	panic("unimplemented")
}