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

// RabinMillerTest reads the input number n (and optional exponent r) from a file,
// performs the Rabin-Miller primality test, and writes the result to an output file.
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

	// Perform Rabin-Miller test
	factor, composite := RabinMillerFunc(n, r)
	var result string
	if factor != nil {
		result = factor.String() // write factor
	} else if composite {
		result = "na pewno złożona"
	} else {
		result = "prawdopodobnie pierwsza"
	}

	// Output result
	fmt.Println("Wynik:", result)
	log.Printf("Rabin-Miller result: %s", result)

	// Write to file
	err = os.WriteFile(OutputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}

// RabinMillerFunc performs the Rabin-Miller probabilistic primality test on a given number n.
// Optionally, a universal exponent r can be provided instead of using n−1.
func RabinMillerFunc(n, r *big.Int) (*big.Int, bool) {
	one := big.NewInt(1)
	two := big.NewInt(2)

	// Handle edge cases
	if n.Cmp(two) == 0 {
		return nil, false // 2 is prime
	}
	if new(big.Int).Mod(n, two).Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(2), true // even => composite
	}

	// Use r if provided, otherwise r = n - 1
	var exponent = new(big.Int)
	if r != nil {
		exponent.Set(r)
	} else {
		exponent.Sub(n, one)
	}

	// Decompose exponent: exponent = m * 2^k
	m := new(big.Int).Set(exponent)
	k := 0
	for new(big.Int).Mod(m, two).Cmp(big.NewInt(0)) == 0 {
		m.Div(m, two)
		k++
	}

	// Run several iterations
	for i := 0; i < Iterations; i++ {
		a, err := helpers.CryptoRandBigIntBetween(two, new(big.Int).Sub(n, two))
		if err != nil {
			log.Printf("Random a generation error: %v", err)
			return nil, true
		}

		// Step 1: check gcd(a, n) for quick factor
		gcd := new(big.Int).GCD(nil, nil, a, n)
		if gcd.Cmp(one) != 0 {
			return gcd, true // non-trivial factor found
		}

		// Step 2: compute a^m mod n
		b := new(big.Int).Exp(a, m, n)
		if b.Cmp(one) == 0 || b.Cmp(new(big.Int).Sub(n, one)) == 0 {
			continue // strong pseudoprime in this iteration
		}

		// Step 3: square b repeatedly: b_j+1 = b_j^2 mod n
		strong := false
		for j := 0; j < k-1; j++ {
			b.Exp(b, two, n)

			if b.Cmp(new(big.Int).Sub(n, one)) == 0 {
				strong = true // passed this round
				break
			}

			if b.Cmp(one) == 0 {
				// Found non-trivial square root of 1 ⇒ factor
				bMinus := new(big.Int).Sub(b, one)
				bPlus := new(big.Int).Add(b, one)
				f1 := new(big.Int).GCD(nil, nil, bMinus, n)
				f2 := new(big.Int).GCD(nil, nil, bPlus, n)

				if f1.Cmp(one) != 0 && f1.Cmp(n) != 0 {
					return f1, true
				}
				if f2.Cmp(one) != 0 && f2.Cmp(n) != 0 {
					return f2, true
				}
				return nil, true // definitely composite
			}
		}

		if !strong {
			return nil, true // definitely composite
		}
	}

	return nil, false // probably prime
}

