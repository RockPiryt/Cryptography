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


// Function to set logger.
func SetLogger(){
	os.MkdirAll("logs", os.ModePerm)
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetPrefix("[ONE TIME PAD]")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

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

// Cut text to specified number of lines with 64 characters each.
func formatText(text string, maxLines int) (string, error) {
	const lineLength = 64
	maxChars := lineLength * maxLines

	if len(text) > maxChars {
		text = text[:maxChars]
	}

	if len(text) < maxChars {
		text += strings.Repeat(" ", maxChars-len(text))
	}

	var lines []string
	for i := 0; i < maxChars; i += lineLength {
		lines = append(lines, text[i:i+lineLength])
	}

	return strings.Join(lines, "\n"), nil
}

// PrepareText prepares the input text: reads it, cleans it, and formats to given line count.
func PrepareText(filePath string, maxLines int) (string, error) {
	inputText, err := GetText(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	preparedText, err := CleanText(inputText)
	if err != nil {
		return "", fmt.Errorf("error cleaning text: %v", err)
	}

	formattedText, err := formatText(preparedText, maxLines)
	if err != nil {
		return "", fmt.Errorf("error formatting text: %v", err)
	}

	return formattedText, nil
}

// Function to validate the key.
func ValidateKey(text string) (error) {
	if len(text) == 0 {
		return fmt.Errorf("key/text is empty")
	}

	// Check if the key contains only uppercase, lowercase letters, and spaces.
	for _, char := range text {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != ' ' {
			return fmt.Errorf("key/text contains invalid character: %c", char)
		}
	}

	// Check if the key is a single line.
	if strings.Contains(text, "\n") {
		return fmt.Errorf("key/text contains multiple lines")
	}

	// Check if the key is exactly 64 characters.
	if len(text) != 64 {
		return fmt.Errorf("key/text must be exactly 64 characters")
	}

	return nil
}

// Function to get the prepared key.
func GetPreparedKey(keyFile string) (string, error) {
	key, err := PrepareText(keyFile, 1)
	if err != nil {
		return "", fmt.Errorf("failed to prepare key")
	}

	err = ValidateKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to validate key")
	}

	return key, nil
}

// Function to validate the text (for Vigenère cipher).
func ValidateText(text string) (error) {
	if len(text) == 0 {
		return fmt.Errorf("text is empty")
	}

	// Check if the key contains only uppercase, lowercase letters, and spaces.
	for _, char := range text {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != ' ' {
			return fmt.Errorf("key/text contains invalid character: %c", char)
		}
	}

	// Check if the text contains exactly 15 lines.
	lines := strings.Split(text, "\n")
	if len(lines) != 15 {
		return fmt.Errorf("text must have exactly 15 lines")
	}

	// Ensure each line contains exactly 64 characters.
	for _, line := range lines {
		if len(line) != 64 {
			return fmt.Errorf("each line must contain exactly 64 characters, found a line with %d characters", len(line))
		}
	}

	return nil
}

// Function to get the prepared text for Vigenère cipher.
func GetPreparedText(textFile string) (string, error) {
	text, err := PrepareText(textFile, 15)
	if err != nil {
		return "", fmt.Errorf("failed to prepare text")
	}

	err = ValidateText(text)
	if err != nil {
		return "", fmt.Errorf("failed to validate text")
	}

	return text, nil
}
