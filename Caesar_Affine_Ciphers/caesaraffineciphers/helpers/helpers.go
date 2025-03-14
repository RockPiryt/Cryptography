package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func SaveOutput(result string, outputFile string) {
	// Check if the file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Błąd przy tworzeniu pliku: %v", err)
		}
		file.Close()
	}

	// Write the result to the file
	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		log.Fatalf("Błąd przy zapisywaniu wyniku: %v", err)
	}

	fmt.Println("Zapisano wynik do pliku:", outputFile)
}


// ValidateCaesarKey reads and validates the key for Caesar cipher from a file.
// It ensures the key is a single integer between 0 and 25.
func ValidateCaesarKey(filePath string) int {
	// Read the key from the file.
	keyLines, err := GetText(filePath)
	if err != nil {
		log.Fatalf("Failed to read the key: %v", err)
	}

	// Ensure there is only one line in the key file.
	if len(keyLines) != 1 {
		log.Fatalf("The key file should only contain one line with a single integer. Found %d lines.", len(keyLines))
	}

	// Convert the key to an integer.
	key, err := strconv.Atoi(strings.TrimSpace(keyLines[0]))
	if err != nil {
		log.Fatalf("Invalid key value. The key must be a valid integer: %v", err)
	}

	// Check if the key is within the valid range for Caesar cipher (0 to 25).
	if key < 0 || key > 25 {
		log.Fatalf("The key for Caesar cipher must be between 0 and 25. Found: %d", key)
	}

	return key
}







