// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Function to set logger.
func SetLogger() {
	os.MkdirAll("logs", os.ModePerm)
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetPrefix("[Steganpgraphy] ")
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

func TextToHex(input string) string {
	return hex.EncodeToString([]byte(input))
}


func ReadHexBits(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	hexStr := strings.TrimSpace(string(content))
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	var bits strings.Builder
	for _, b := range bytes {
		bits.WriteString(fmt.Sprintf("%08b", b))
	}
	return bits.String(), nil
}

func BitsToHex(bits string) string {
	var bytes []byte
	for i := 0; i+8 <= len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b <<= 1
			if bits[i+j] == '1' {
				b |= 1
			}
		}
		bytes = append(bytes, b)
	}
	return hex.EncodeToString(bytes)
}

// SaveHexToFile converts the input string to hex and writes it to the given file.
func SaveHexToFile(input, filename string) error {
	hexStr := hex.EncodeToString([]byte(input))
	return os.WriteFile(filename, []byte(hexStr), 0644)
}


// ClearHtml removes HTML comments, empty lines, and trailing spaces from each line, then saves the cleaned content to 'clearfile.html'.
func ClearHtml(htmlFile string) error {
	content, err := os.ReadFile(htmlFile)
	if err != nil {
		return err
	}

	// Remove HTML comments: <!-- ... -->
	commentRegex := regexp.MustCompile(`(?s)<!--.*?-->`)
	cleaned := commentRegex.ReplaceAllString(string(content), "")

	// Process line by line
	var cleanedLines []string
	scanner := bufio.NewScanner(strings.NewReader(cleaned))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, " \t")      // Remove trailing spaces/tabs
		line = strings.TrimLeft(line, "\t")        // Optional: remove leading tabs
		if strings.TrimSpace(line) != "" {         // Skip completely empty lines
			cleanedLines = append(cleanedLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	finalContent := strings.Join(cleanedLines, "\n")
	return os.WriteFile(htmlFile, []byte(finalContent), 0644)
}

// IsHex checks if the bit string has a length that's a multiple of 4
// and contains only '0' or '1'.
func IsHex(bits string) bool {
	bits = strings.TrimSpace(bits)
	if len(bits)%4 != 0 {
		return false
	}
	for _, ch := range bits {
		if ch != '0' && ch != '1' {
			return false
		}
	}
	return true
}

// ReadFileContent reads the content of a file and returns it as a trimmed string.
func ReadFileContent(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
