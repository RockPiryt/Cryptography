// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
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
	log.SetPrefix("[Rabin-Miller/Fermat] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Function ReadData reads and returns non-empty trimmed lines from a text file.
func ReadData(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

// ParseInput parses lines into a number n and optional exponent r.
func ParseInput(lines []string) (*big.Int, *big.Int, error) {
	if len(lines) < 1 {
		return nil, nil, fmt.Errorf("missing number n in input")
	}

	n := new(big.Int)
	if _, ok := n.SetString(lines[0], 10); !ok {
		return nil, nil, fmt.Errorf("invalid number n: %s", lines[0])
	}

	var r *big.Int
	if len(lines) == 2 {
		r = new(big.Int)
		if _, ok := r.SetString(lines[1], 10); !ok {
			return nil, nil, fmt.Errorf("invalid exponent r: %s", lines[1])
		}
	} else if len(lines) >= 3 {
		r1 := new(big.Int)
		r2 := new(big.Int)
		if _, ok := r1.SetString(lines[1], 10); !ok {
			return nil, nil, fmt.Errorf("invalid value r1: %s", lines[1])
		}
		if _, ok := r2.SetString(lines[2], 10); !ok {
			return nil, nil, fmt.Errorf("invalid value r2: %s", lines[2])
		}
		r = new(big.Int).Mul(r1, r2)
		r.Sub(r, big.NewInt(1))
	}

	return n, r, nil
}
