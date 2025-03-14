package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

//Author: Paulina Kimak

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







