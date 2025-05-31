// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"rabin/helpers"
)

const (
	EntryFile  = "files/wejscie.txt"
	OutputFile = "files/wyjscie.txt"
	Iterations = 40
)

var MessageString string = "Bla"

func ExecuteProgram(operation string) error {
	switch operation {
	case "f":
		// Fermat test
		err := FermatTest(EntryFile)
		if err != nil {
			return fmt.Errorf("error during Fermat test %v", err)
		}
		log.Println("[INFO] Fermat test executed.")
		return nil

	case "r":
		// Rabin-Miller test
		err := RabinMillerTest(EntryFile)
		if err != nil {
			return fmt.Errorf("failed during Rabin-Miller test: %v", err)
		}
		log.Println("[INFO] Rabin-Miller test executed")
		return nil
	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}

func FermatTest(EntryFile string) error {
	log.Println("Fermat test start")

	// Read data
	lines, err := helpers.ReadData(EntryFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	// Parse only 'n' from the input
	n, _, err := helpers.ParseInput(lines)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	fmt.Println("n =", n)

	// Perform Fermat test
	result := "na pewno złożona"
	if FermatFunc(n) {
		result = "prawdopodobnie pierwsza"
	}

	// Output to console
	fmt.Println(result)
	log.Printf("Fermat test result: %s", result)

	// Save result to output file
	err = os.WriteFile(OutputFile, []byte(result), 0644)
	if err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	return nil
}

// FermatFunc performs the Fermat primality test on number n.
func FermatFunc(n *big.Int) bool {
	one := big.NewInt(1)
	nMinusOne := new(big.Int).Sub(n, big.NewInt(1))
	nMinusTwo := new(big.Int).Sub(n, big.NewInt(2))

	for i := 0; i < Iterations; {
		// Generate random a in [2, n−2]
		a, err := helpers.CryptoRandBigIntBetween(big.NewInt(2), nMinusTwo)
		if err != nil {
			log.Printf("Error generating random number: %v", err)
			return false
		}

		// Skip a if not coprime with n (gcd(a, n) != 1)
		if new(big.Int).GCD(nil, nil, a, n).Cmp(one) != 0 {
			continue // skip this iteration without incrementing i
		}

		// Compute a^(n−1) mod n
		// If result != 1, then n is definitely not prime
		if new(big.Int).Exp(a, nMinusOne, n).Cmp(one) != 0 {
			return false // definitely composite
		}

		i++ // increment only when a valid base was tested
	}

	return true // probably prime
}

func RabinMillerTest(EntryFile string) error {
	log.Println("Rabin-Miller test start")

	//Read data
	lines, err := helpers.ReadData(EntryFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	n, r, err := helpers.ParseInput(lines)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	fmt.Println("n =", n)
	if r != nil {
		fmt.Println("r =", r)
	} else {
		fmt.Println("No exponent r provided.")
	}

	return nil
}
