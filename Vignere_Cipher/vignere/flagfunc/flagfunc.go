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
	InputTextHelper string
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
	case "e":
		// Encode the text from plain.txt using the key from key.txt and saves the result to crypto.txt
		params.InputText = "files/plain.txt"
		params.InputKey = "files/key.txt"
		params.OutputText = "files/crypto.txt"
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

func VignereCryptAnalysis(s1, s2 string) {
	panic("unimplemented")
}

func CipherOperations(params CipherParams) {
	panic("unimplemented")
}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile() error {
	// Check if the org.txt file exists.
	_, err := os.Stat("files/org.txt")
	if os.IsNotExist(err) {
		log.Println("plik org.txt nie istnieje")
		return fmt.Errorf("plik org.txt nie istnieje")
	} else if err != nil {
		log.Printf("błąd przy sprawdzaniu istnienia pliku org.txt: %v", err)
		return fmt.Errorf("błąd przy sprawdzaniu istnienia pliku org.txt: %v", err)
	}

	lines, err := helpers.GetText("files/org.txt")
	if err != nil {
		log.Printf("błąd przy odczycie pliku org.txt: %v", err)
		return fmt.Errorf("błąd przy odczycie pliku org.txt: %v", err)
	}

	if len(lines) == 0 {
		log.Println("plik org.txt jest pusty")
		return fmt.Errorf("plik org.txt jest pusty")
	}

	inputText := strings.Join(lines, "\n")
	preparedText, err := helpers.CleanText(inputText)
	if err != nil {
		log.Printf("błąd przy czyszczeniu tekstu: %v", err)
		return fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
	}

	// Save the prepared text to the plain.txt file
	err = os.WriteFile("files/plain.txt", []byte(preparedText), 0644)
	if err != nil {
		log.Printf("błąd przy zapisywaniu do pliku plain.txt: %v", err)
		return fmt.Errorf("błąd przy zapisywaniu do pliku plain.txt: %v", err)
	}

	return nil
}


