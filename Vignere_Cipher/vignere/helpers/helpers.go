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
func GetText(filename string) ([]string, error) {
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


func PrepareText(filePath string) (string, error) {
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			log.Println("plik %s nie istnieje", filePath)
			return "", fmt.Errorf("plik %s nie istnieje", filePath)
		} else if err != nil {
			log.Printf("błąd przy sprawdzaniu istnienia pliku %s %v", filePath, err)
			return "", fmt.Errorf("błąd przy sprawdzaniu istnienia pliku %s: %v", filePath, err)
		}
	
		lines, err := GetText(filePath)
		if err != nil {
			log.Printf("błąd przy odczycie pliku %s: %v", filePath, err)
			return "", fmt.Errorf("błąd przy odczycie pliku %s: %v",filePath, err)
		}
	
		if len(lines) == 0 {
			log.Println("plik %s jest pusty", filePath)
			return "", fmt.Errorf("plik %s jest pusty", filePath)
		}
	
		inputText := strings.Join(lines, "\n")
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

	// Check if the key contains only small letters.
	for _, char := range key {
		if !unicode.IsLetter(char) {
			return "", fmt.Errorf("klucz zawiera niedozwolony znak: %c", char)
		}
	}


	return key, nil
}


// Extended Euclidean algorithm
func ExtendedGCD(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, x1, y1 := ExtendedGCD(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return gcd, x, y
}

// ModInverseExtended calculates the modular inverse of a mod m using the extended Euclidean algorithm.
// If the modular inverse does not exist, it returns an error.
func ModInverseExtended(a, m int) (int, error) {
	gcd, x, _ := ExtendedGCD(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("brak modularnej odwrotności dla a=%d (mod %d)", a, m)
	}
	return (x%m + m) % m, nil // Ensure non-negative result
}