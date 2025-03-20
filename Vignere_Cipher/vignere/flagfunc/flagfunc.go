package flagfunc

import (
	"fmt"
	"log"
	"os"
	"strings"

	"vignere/helpers"
)

//Author: Paulina Kimak

// Struct to store cipher parameters.
type CipherParams struct {
	Operation       string
	InputText       string
	PreparedText    string
	InputKey        string
	OutputText      string
	OutputKey       string
	CipherFunc      func(string, int, int, string) (string, error)
	KeyFinder       func(string, string) (int, int)
}

// Generic function to handle Vignere cipher
func ExecuteCipher(operation string) {
	params := CipherParams{
		Operation:  operation,
	}

	switch operation {
	case "p":
		// Prepare the text for encryption and save it to plain.txt
		params.InputText = "files/org.txt"
		params.PreparedText = "files/plain.txt"
		err := CreatePlainFile(params.InputText, params.PreparedText)
		if err != nil {
			fmt.Printf("błąd przy przygotowywaniu tekstu: %v", err)
		}
	case "e":
		// Encode the text from plain.txt using the key from key.txt and saves the result to crypto.txt
		params.InputKey = "files/key.txt"
		params.OutputText = "files/crypto.txt"
		EncodeText(params.InputKey, params.OutputText) 
	case "d":
		// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
		params.InputText = "files/crypto.txt"
		params.InputKey = "files/key.txt"
		params.OutputText = "files/decrypt.txt"
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		params.InputText = "files/crypto.txt"
		params.OutputText = "files/decrypt.txt"
		VignereCryptAnalysis(params.InputText, params.OutputText)

	default:
		fmt.Println("Nieobsługiwana operacja.")
		return
	}	
		CipherOperations(params)
}


func CipherOperations(params CipherParams) {
	panic("unimplemented")
}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile(inputFile string, outputFile string) error {
	preparedText, err := helpers.PrepareText(inputFile)
	if err != nil {
		log.Printf("błąd przy czyszczeniu tekstu: %v", err)
		return fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
	}

	err = helpers.SaveOutput(preparedText, outputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
}

// Function to get the key for Vignere cipher.
func GetKey(inputKey string) (string, error) {
	key, err := helpers.ValidateKey(inputKey)
	if err != nil {
		return "", fmt.Errorf("nie udało się zwalidować klucza")
	}

	numKey,err:=helpers.ConverseKey(key)
	if err != nil {
		return "", fmt.Errorf("nie udało się przekonwertować klucza")
	}

	fmt.Println("Przekonwertowany klucz: ", numKey)
	return numKey, nil
}

// Function to encode the text using Vignere cipher.
func EncodeText(inputKey, outputText string) (string, error) {
	// Check if the plaintext file exists
	plainTextFile := "files/plain.txt"
	if _, err := os.Stat(plainTextFile); os.IsNotExist(err) {
		return "", fmt.Errorf("file %s does not exist", plainTextFile)
	} else if err != nil {
		return "", fmt.Errorf("error checking file %s: %v", plainTextFile, err)
	}

	// Read plaintext from the file
	lines, err := helpers.GetText(plainTextFile)
	if err != nil {
		return "", fmt.Errorf("błąd przy odczycie pliku %s: %v",plainTextFile, err)
	}

	if len(lines) == 0 {
		return "", fmt.Errorf("plik %s jest pusty", plainTextFile)
	}

	inputText := strings.Join(lines, "\n")
	fmt.Println("Odczytany tekst: ", inputText)
	
	// Prepare the key for encryption.
	numKey,err := GetKey(inputKey)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać poprawnego klucza")
	}

	return numKey, nil
}


func VignereCryptAnalysis(s1, s2 string) {
	panic("unimplemented")
}