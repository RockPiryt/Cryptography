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
		fmt.Printf("Znaleziony dzielnik: %s\n", factor.String())
		result = factor.String()
	} else if composite {
		fmt.Printf("Brak dzielnika, ale liczba jest na pewno złożona: composite = %v\n", composite)
		result = "na pewno złożona"
	} else {
		fmt.Printf("Liczba przeszła test: prawdopodobnie pierwsza (composite = %v)\n", composite)
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





// RabinMillerFunc wykonuje probabilistyczny test Rabina-Millera.
// Zwraca:
//   – dzielnik != nil  → znaleziono nietrywialny dzielnik n,
//   – composite == true  → liczba na pewno złożona (bez dzielnika),
//   – composite == false → liczba prawdopodobnie pierwsza.
func RabinMillerFunc(n, r *big.Int) (*big.Int, bool) {
	one := big.NewInt(1)
	two := big.NewInt(2)
	nMinusTwo := new(big.Int).Sub(n, two)

	// 2 jest pierwsza
	if n.Cmp(two) == 0 {
		return nil, false
	}
	// Parzysta? → złożona, dzielnik 2
	if new(big.Int).Mod(n, two).Sign() == 0 {
		return big.NewInt(2), true
	}

	// wykładnik: r (jeśli podany) lub n-1
	exponent := new(big.Int)
	if r != nil {
		exponent.Set(r)
	} else {
		exponent.Sub(n, one) // n-1
	}

	// Rozkład: exponent = m · 2^k  (m – nieparzyste)
	m := new(big.Int).Set(exponent)
	k := 0
	for new(big.Int).And(m, one).Sign() == 0 { // dopóki parzyste
		m.Rsh(m, 1) // m /= 2
		k++
	}

	// Iteracje testu
	for i := 0; i < Iterations; i++ {
		// losowe  a ∈ [2, n−2]
		a, err := helpers.CryptoRandBigIntBetween(two, nMinusTwo)
		if err != nil {
			log.Printf("losowanie a nieudane: %v", err)
			return nil, true
		}

		// szybki dzielnik?
		if d := new(big.Int).GCD(nil, nil, a, n); d.Cmp(one) != 0 {
			return d, true
		}

		// b0 = a^m mod n
		b := new(big.Int).Exp(a, m, n)
		if b.Cmp(one) == 0 || b.Cmp(new(big.Int).Sub(n, one)) == 0 {
			continue // silny pseudopierwszy w tej iteracji
		}

		// powtarzaj: b_j+1 = b_j^2  (mod n)
		prev := new(big.Int).Set(b) // b_j
		strong := false
		for j := 0; j < k-1; j++ {
			b.Exp(prev, two, n) // b = prev^2

			// jeśli b == n−1 → przechodzi tę iterację
			if b.Cmp(new(big.Int).Sub(n, one)) == 0 {
				strong = true
				break
			}

			// jeśli b == 1  i  prev ≠ ±1 → mamy rozkład
			if b.Cmp(one) == 0 {
				// gcd(prev−1, n)
				if d := new(big.Int).GCD(nil, nil,
					new(big.Int).Sub(prev, one), n); d.Cmp(one) != 0 && d.Cmp(n) != 0 {
					return d, true
				}
				// gcd(prev+1, n)
				if d := new(big.Int).GCD(nil, nil,
					new(big.Int).Add(prev, one), n); d.Cmp(one) != 0 && d.Cmp(n) != 0 {
					return d, true
				}
				return nil, true // złożona, bez dzielnika
			}
			prev.Set(b) // przechodzimy do kolejnego b_j
		}

		// jeśli w ogóle nie było -1
		if !strong {
			return nil, true // liczba złożona
		}
	}
	// Po wszystkich iteracjach – prawdopodobnie pierwsza
	return nil, false
}
