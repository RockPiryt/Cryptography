// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Function to set logger.
func SetLogger() {
	os.MkdirAll("logs", os.ModePerm)
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

	if scanner.Err() != nil {
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
func ValidateKey(text string) error {
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

// Function to validate the text for One Time Pad.
func ValidateText(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("text is empty")
	}

	// Check if the text contains only uppercase, lowercase letters, spaces, and newlines.
	for _, char := range text {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != ' ' && char != '\n' {
			log.Printf("Invalid character found: %c (ASCII: %d)", char, char)
			return fmt.Errorf("text contains invalid character: %c", char)
		}
	}

	// Check if the text contains exactly 15 lines, allowing a newline at the end of the last line.
	lines := strings.Split(text, "\n")

	if len(lines) > 15 {
		return fmt.Errorf("text contains more than 15 lines")
	}

	// If there are exactly 15 lines, the last line should allow an empty string (which would be a newline at the end)
	if len(lines) == 15 && lines[14] != "" && len(lines[14]) != 64 {
		return fmt.Errorf("the 15th line must be exactly 64 characters, found %d characters", len(lines[14]))
	}

	// Ensure each line contains exactly 64 characters, except possibly the last line (which can end with a newline).
	for i, line := range lines[:14] { // Check the first 14 lines, they must be exactly 64 characters.
		if len(line) != 64 {
			return fmt.Errorf("line %d must contain exactly 64 characters, found %d characters", i+1, len(line))
		}
	}

	// Ensure no extra newlines after the 15th line
	if len(lines) == 15 && lines[14] != "" && len(lines[14]) > 64 {
		return fmt.Errorf("the last line must not exceed 64 characters")
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
		return "", fmt.Errorf("failed to validate text: %v", err)
	}

	return text, nil
}

// Function to convert text to hexadecimal representation
func TextToHex(text string) string {
	var hexString string
	for i := 0; i < len(text); i++ {
		hexString += fmt.Sprintf("%02X", text[i])
	}
	return hexString
}

// Function to convert a hexadecimal string to bytes
func HexToBytes(hex string) ([]byte, error) {
	bytes := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		val, err := strconv.ParseUint(hex[i:i+2], 16, 8)
		if err != nil {
			return nil, err
		}
		bytes[i/2] = byte(val)
	}
	return bytes, nil
}

// Function to compute XOR between two byte slices
func XORBytes(b1, b2 []byte) []byte {
	result := make([]byte, len(b1))
	for i := 0; i < len(b1); i++ {
		result[i] = b1[i] ^ b2[i]
	}
	return result
}

// Read the crypto.txt file
func ReadCiphertext(fileName string) ([]string, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	// Return the list of lines
	return lines, nil
}

func ValidateCryptoHexFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	var expectedLength int = -1
	validHex := regexp.MustCompile(`^[0-9A-Fa-f]+$`)

	hasError := false

	for scanner.Scan() {
		line := scanner.Text()

		if !validHex.MatchString(line) {
			fmt.Printf("Line %d contains invalid characters: %s\n", lineNum, line)
			hasError = true
		}

		if len(line)%2 != 0 {
			fmt.Printf("Line %d has odd length: %d\n", lineNum, len(line))
			hasError = true
		}

		if expectedLength == -1 {
			expectedLength = len(line)
		} else if len(line) != expectedLength {
			fmt.Printf("Line %d has a different length (%d) than expected (%d)\n", lineNum, len(line), expectedLength)
			hasError = true
		}

		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading the file: %w", err)
	}

	if hasError {
		return fmt.Errorf("file contains errors")
	}

	return nil
}

// Function to print the positions of spaces in each line of a text file
func PrintSpacePositions(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error during open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		// log.Printf("Linie numer %2d (len: %d): ", lineNumber, len(line))
		var spacePositions []int
		for i, ch := range line {
			if ch == ' ' {
				spacePositions = append(spacePositions, i)
			}
		}
		if len(spacePositions) == 0 {
			fmt.Println("no space")
		} else {
			log.Printf("spaces: %v\n", spacePositions)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("eroor during reading file: %v", err)
	}
	return nil
}

// Function to find columns without spaces in a text file
func FindColumnsWithoutSpaces(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error during open file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error during read file: %v", err)
	}

	if len(lines) == 0 {
		return fmt.Errorf("file is empty")
	}

	lineLength := len(lines[0])
	// Check if all lines have the same length
	for i, line := range lines {
		if len(line) != lineLength {
			return fmt.Errorf("line %d has lenght %d, exepected %d", i+1, len(line), lineLength)
		}
	}

	noSpaceCols := make([]bool, lineLength)
	for i := 0; i < lineLength; i++ {
		noSpaceCols[i] = true
	}

	// check each line for spaces
	for _, line := range lines {
		for i, ch := range line {
			if ch == ' ' {
				noSpaceCols[i] = false
			}
		}
	}

	var result []int
	for i, isClean := range noSpaceCols {
		if isClean {
			result = append(result, i)
		}
	}

	if len(result) == 0 {
		log.Println("Each column contains at least one space.")
	} else {
		log.Printf("Columns without spaces:\n%v\n", result)
	}

	return nil
}
