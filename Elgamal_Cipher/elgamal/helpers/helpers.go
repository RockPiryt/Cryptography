// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
		line := strings.TrimSpace(scanner.Text()) 
		if line == "" {
			continue 
		}
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
func ConvertStringToBigInt(message string, filePath string) error {
	
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

// Use Extended Euklides Algorithm to find modular inverse 
func ModInverse(a, m *big.Int) (*big.Int, error) {
	g := new(big.Int)
	x := new(big.Int)
	y := new(big.Int)

	//x is modular inverse a mod m, if g=gcd(a, m) == 1
	g.GCD(x, y, a, m)
	//cmp if 2 numbers are equal then return 0
	if g.Cmp(big.NewInt(1)) != 0 {
		return nil, errors.New("no inverse exists")
	}
	return x.Mod(x, m), nil
}

//This function checks whether two numbers a and b are coprime, 
//meaning their greatest common divisor (GCD) is equal to 1.
func IsCoprime(a, b *big.Int) bool {
	gcd := new(big.Int)

	// Calculate GCD (x i y are not used so nil)
	gcd.GCD(nil, nil, a, b)

	isOne := gcd.Cmp(big.NewInt(1)) == 0

	return isOne
}


// Function calculates SHA256 hash of the input string, converts it to a big int, and saves it.
func CreateShortcutSHA(text, filename string) error {
	// Oblicz skrót SHA256
	hash := sha256.Sum256([]byte(text))

	// Zamień skrót (bytes) na BigInt
	m := new(big.Int).SetBytes(hash[:])

	// Zapisz big.Int do pliku
	err := WriteBigIntsToFile(filename, []*big.Int{m})
	if err != nil {
		return fmt.Errorf("failed to write hashed message to file: %v", err)
	}

	log.Printf("[INFO] Hashed text \"%s\" → SHA256 = %s\nSaved to %s", text, hex.EncodeToString(hash[:]), filename)
	return nil
}
