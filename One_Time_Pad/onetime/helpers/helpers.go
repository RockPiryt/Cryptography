// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"fmt"
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

// Function to read text from txt file.
func ReadText(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// Close file no matter what
	defer file.Close() 

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
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
			return fmt.Errorf("error creating file: %v", err)
		}
		file.Close()
	}

	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("error writing result to file: %v", err)
	}

	return nil
}

// CleanText prepares the text for encryption by removing all non-letter and non-space characters. It keeps letters (both uppercase and lowercase) and space only.
func CleanText(input string) (string, error) {
	var cleanedText []rune

	if len(input) == 0 {
		return "", fmt.Errorf("input text is empty")
	}

	for _, char := range input {
		if unicode.IsLetter(char) || char == ' ' {
			cleanedText = append(cleanedText, char)
		}
	}

	if len(cleanedText) == 0 {
		return "", fmt.Errorf("text contains no valid characters (letters or spaces)")
	}

	return string(cleanedText), nil
}

func GetText(filePath string) (string, error) {
	_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file %s does not exist", filePath)
		} else if err != nil {
			return "", fmt.Errorf("error checking existence of file %s: %v", filePath, err)
		}
	
		lines, err := ReadText(filePath)
		if err != nil {
			return "", fmt.Errorf("error reading file %s: %v", filePath, err)
		}
	
		if len(lines) == 0 {
			return "", fmt.Errorf("file %s is empty", filePath)
		}
	
		inputText := strings.Join(lines, "\n")
		return inputText, nil
}

// Cut text to 15 lines of 64 characters each.
func cutText(text string) (string, error) {
	const lineLength = 64
	const maxLines = 15
	maxChars := lineLength * maxLines

	// Cut the text to 15 lines of 64 characters each
	if len(text) > maxChars {
		text = text[:maxChars]
	}

	// Fill the text with spaces to make it a multiple of lineLength
	if len(text) < maxChars {
		text += strings.Repeat(" ", maxChars-len(text))
	}

	var lines []string
	for i := 0; i < maxChars; i += lineLength {
		lines = append(lines, text[i:i+lineLength])
	}

	result := strings.Join(lines, "\n")

	return result, nil
}

// Function to prepare text for encryption, cleans non-letter characters and converts to lowercase.
func PrepareText(filePath string) (string, error) {
		inputText,err := GetText(filePath)
		if err != nil {
			return "", fmt.Errorf("error reading file %s: %v", filePath, err)
		}
		preparedText, err := CleanText(inputText)
		if err != nil {
			return "", fmt.Errorf("error cleaning text: %v", err)
		}

		cuttedText, err := cutText(preparedText)
		if err != nil {	
			return "", fmt.Errorf("error trimming text: %v", err)
		}
		return cuttedText, nil
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

