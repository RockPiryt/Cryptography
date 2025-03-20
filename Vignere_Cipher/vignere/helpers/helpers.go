package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

// Author: Paulina Kimak

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

// Function to read and validate the key for Vignere cipher
func ValidateKey(filePath string) (string, error) {
	key, err := PrepareText(filePath)
	if err != nil {
		return "", fmt.Errorf("nie udało się przygotować klucza")
	}

	// Check if the key is empty.
	if len(key) == 0 {
		return "", fmt.Errorf("klucz jest pusty")
	}

	// Check if the key contains only lowercase English letters.
	for _, char := range key {
		if char < 'a' || char > 'z' { // Ensure only 'a' to 'z'
			return "", fmt.Errorf("klucz zawiera niedozwolony znak: %c", char)
		}
	}

	return key, nil
}

// Function to convert the key to a slice of shift values.
func ConverseKey(key string) ([]int, error) {
	var convertedKey []int

	for _, char := range key {
		// Convert letter to shift value (a = 0, b = 1, ..., z = 25)
		numValue := int(char - 'a')
		convertedKey = append(convertedKey, numValue)
	}

	return convertedKey, nil
}


// Function to get the key for Vignere cipher.
func GetKey(inputKey string) ([]int, error) {
	// Read alpha key from file and validate key.
	key, err := ValidateKey(inputKey)
	if err != nil {
		return nil, fmt.Errorf("nie udało się zwalidować klucza")
	}

	// Convert alpha key to numeric key slice
	numKey,err := ConverseKey(key)
	if err != nil {
		return nil, fmt.Errorf("nie udało się przekonwertować klucza")
	}
	return numKey, nil
}
