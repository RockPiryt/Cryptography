// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
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

//TODO: remove white signs before read from file
func ReadBigIntsFromFile(path string, count int) ([]*big.Int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var values []*big.Int
	for scanner.Scan() {
		line := scanner.Text()
		n := new(big.Int)
		n, ok := n.SetString(line, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number in %s", path)
		}
		values = append(values, n)
		if len(values) == count {
			break
		}
	}
	if len(values) != count {
		return nil, fmt.Errorf("expected %d numbers in %s", count, path)
	}
	return values, nil
}

func WriteBigIntsToFile(path string, values []*big.Int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, val := range values {
		_, err := file.WriteString(val.String() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func RandomBigInt(max *big.Int) (*big.Int, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Function saves a string message as a big.Int to plain.txt
func SavePlainMessageAsBigInt(message string, filePath string) error {
	// Convert string to []byte, then to big.Int
	m := new(big.Int).SetBytes([]byte(message))

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(m.String() + "\n")
	return err
}