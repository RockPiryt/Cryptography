// Author: Paulina Kimak
package helpers

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
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
	log.SetPrefix("[Elgamal] ")
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


func readHexBits(filename string) (string, error) {
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

func bitsToHex(bits string) string {
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