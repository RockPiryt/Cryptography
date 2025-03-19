package cryptofunc

import (
	"fmt"
	"log"
	"os"
	"strings"

	"vignere/helpers"
)

//Author: Paulina Kimak

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


