// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)


const Alphabet = "abcdefghijklmnopqrstuvwxyz"

var AlphabetLen = len(Alphabet)

// Function to count selected flags
func CountSelectedFlags(flags []*bool) int {
	count := 0
	for _, f := range flags {
		if *f {
			count++
		}
	}
	return count
}

// Function to read text from txt file
func ReadText(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if scanner.Err() != nil{
		return nil, scanner.Err()
	}
	
	return lines, nil
}


func SaveOutput(result string, outputFile string) error {
	// Check if the file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("błąd przy tworzeniu pliku: %v", err)
		}
		file.Close()
	}

	// Write the result to the file
	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("błąd przy zapisywaniu wyniku: %v", err)
	}

	fmt.Println("Zapisano wynik do pliku:", outputFile)
	return nil
}

// Function to prepare text for encryption, cleans non-letter characters and converts to lowercase.
func CleanText(input string) (string, error) {
	var cleanedText []rune

	// Check if the input text is empty.
	if len(input) == 0 {
		return "", fmt.Errorf("input text is empty")
	}
	// Prepare the text for encryption/
	for _, char := range input {
		if unicode.IsLetter(char) {
			char = unicode.ToLower(char)
			cleanedText = append(cleanedText, char)
		}
	}

	if len(cleanedText) == 0 {
		return "", fmt.Errorf("tekst nie zawiera liter do przetworzenia")

	}
	return string(cleanedText), nil
}

func GetText(filePath string) (string, error) {
	_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			fmt.Printf("plik %s nie istnieje", filePath)
			return "", fmt.Errorf("plik %s nie istnieje", filePath)
		} else if err != nil {
			log.Printf("błąd przy sprawdzaniu istnienia pliku %s %v", filePath, err)
			return "", fmt.Errorf("błąd przy sprawdzaniu istnienia pliku %s: %v", filePath, err)
		}
	
		lines, err := ReadText(filePath)
		if err != nil {
			fmt.Printf("błąd przy odczycie pliku %s: %v", filePath, err)
			return "", fmt.Errorf("błąd przy odczycie pliku %s: %v",filePath, err)
		}
	
		if len(lines) == 0 {
			fmt.Printf("plik %s jest pusty", filePath)
			return "", fmt.Errorf("plik %s jest pusty", filePath)
		}
	
		inputText := strings.Join(lines, "\n")
		return inputText, nil
}

// Function to prepare text for encryption, cleans non-letter characters and converts to lowercase.
func PrepareText(filePath string) (string, error) {
		// Get input text.
		inputText,err := GetText(filePath)
		if err != nil {
			return "", fmt.Errorf("błąd przy odczycie pliku %s: %v", filePath, err)
		}
		preparedText, err := CleanText(inputText)
		if err != nil {
			log.Printf("błąd przy czyszczeniu tekstu: %v", err)
			return "", fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
		}
		return preparedText, nil
}

// Function to validate the key for Vignere cipher
func Validate(text string) (error) {
	// Check if the key is empty.
	if len(text) == 0 {
		return fmt.Errorf("klucz/tekst jest pusty")
	}

	// Check if the key contains only lowercase English letters.
	for _, char := range text {
		if char < 'a' || char > 'z' { // Ensure only 'a' to 'z'
			return fmt.Errorf("klucz/tekst zawiera niedozwolony znak: %c", char)
		}
	}

	return nil
}

// Function to get the key for Vignere cipher.
func GetPreparedKey(keyFile string) (string, error) {

	key, err := PrepareText(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się przygotować klucza")
	}
	// Read alpha key from file and validate key.
	err = Validate(key)
	if err != nil {
		return  "", fmt.Errorf("nie udało się zwalidować klucza")
	}

	return key, nil
}

// Function to get the key for Vignere cipher.
func GetPreparedText(textFile string) (string, error) {
	text, err := PrepareText(textFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się przygotować klucza")
	}
	// Read alpha key from file and validate key.
	err = Validate(text)
	if err != nil {
		return  "", fmt.Errorf("nie udało się zwalidować klucza")
	}

	return text, nil
}
