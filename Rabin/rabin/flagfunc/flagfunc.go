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

	log.Println("n =", n)

	// Perform Fermat test
	result := "na pewno złożona"
	if FermatFunc(n) {
		result = "prawdopodobnie pierwsza"
	}

	// Output to console
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

	log.Println("n =", n)
	if r != nil {
		log.Println("r =", r)
	} else {
		log.Println("No exponent r provided.")
	}

	// Perform Rabin-Miller test
	factor, composite := RabinMillerFunc(n, r)
	var result string
	if factor != nil {
		result = factor.String()
		log.Printf("Rabin-Miller result: %s", result)
	} else if composite {
		log.Printf("Brak dzielnika, ale liczba jest na pewno złożona: composite = %v\n", composite)
		result = "na pewno złożona"
	} else {
		log.Printf("Liczba przeszła test: prawdopodobnie pierwsza (composite = %v)\n", composite)
		result = "prawdopodobnie pierwsza"
	}

	// Write to file
	err = os.WriteFile(OutputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}


// RabinMillerFunc wykonuje probabilistyczny test Rabina-Millera dla liczby n.
// Jeśli podany jest wykładnik uniwersalny r, używa go zamiast n-1.
// Zwraca:
//   • dzielnik (!= nil)  →  znaleziono nietrywialny dzielnik n,
//   • composite == true  →  n jest na pewno złożona (bez ujawnionego dzielnika),
//   • composite == false →  n jest prawdopodobnie pierwsza.
func RabinMillerFunc(n, r *big.Int) (*big.Int, bool) {
	one := big.NewInt(1)
	two := big.NewInt(2)

	//------------------------------------------------------------
	// 1. Szybkie przypadki brzegowe
	//------------------------------------------------------------
	if n.Cmp(two) == 0 {                 // n = 2  → liczba pierwsza
		return nil, false
	}
	if new(big.Int).Mod(n, two).Sign() == 0 { // n parzysta > 2  → złożona, dzielnik 2
		return big.NewInt(2), true
	}

	//------------------------------------------------------------
	// 2. Wyznacz wykładnik: r (lub n−1)  i rozłóż go na m·2^k
	//------------------------------------------------------------
	exponent := new(big.Int)
	if r != nil {
		exponent.Set(r)                  // używamy uniwersalnego r
	} else {
		exponent.Sub(n, one)             // klasycznie: n-1
	}

	// Rozkład: exponent = m · 2^k  (m nieparzyste, k ≥ 0)
	m := new(big.Int).Set(exponent)
	k := 0
	for new(big.Int).Mod(m, two).Sign() == 0 { // dopóki m jest parzyste
		m.Rsh(m, 1)                            // m /= 2
		k++
	}

	//------------------------------------------------------------
	// 3. Iteracyjne testy z losowymi podstawami a
	//------------------------------------------------------------
	for i := 0; i < Iterations; i++ {
		// losowe  a  z przedziału [2, n−2]
		a, err := helpers.CryptoRandBigIntBetween(two, new(big.Int).Sub(n, two))
		if err != nil {
			log.Printf("Błąd losowania a: %v", err)
			return nil, true
		}

		// --- (a) Szybki test Euklidesa: gcd(a, n) > 1 → dzielnik ---
		if d := new(big.Int).GCD(nil, nil, a, n); d.Cmp(one) != 0 {
			if new(big.Int).Mod(n, d).Sign() == 0 { // upewnij się, że faktycznie dzieli n
				return d, true
			}
			continue // jeśli to nie dzielnik – przejdź do kolejnej iteracji
		}

		// --- (b) Oblicz b0 = a^m  (mod n) ---
		b := new(big.Int).Exp(a, m, n)
		if b.Cmp(one) == 0 || b.Cmp(new(big.Int).Sub(n, one)) == 0 {
			continue // silny pseudopierwszy w tej iteracji → przechodzimy dalej
		}

		//--------------------------------------------------------
		// (c) Powtarzaj: bj+1 = bj² (mod n)  i szukaj rozkładu
		//--------------------------------------------------------
		prev := new(big.Int).Set(b) // prev = bj
		strong := false
		for j := 0; j < k-1; j++ {
			b.Exp(prev, two, n)     // b = prev²  →  bj+1

			// jeśli bj+1 = n−1  → iteracja zaliczona, przechodzimy do kolejnej podstawy
			if b.Cmp(new(big.Int).Sub(n, one)) == 0 {
				strong = true
				break
			}

			// jeśli bj+1 = 1  i  bj ≠ ±1  → mamy nietrywialny pierwiastek z 1
			if b.Cmp(one) == 0 {
				// gcd(bj−1, n)
				if d := new(big.Int).GCD(nil, nil,
					new(big.Int).Sub(prev, one), n); d.Cmp(one) != 0 && d.Cmp(n) != 0 &&
					new(big.Int).Mod(n, d).Sign() == 0 {
					return d, true
				}
				// gcd(bj+1, n)
				if d := new(big.Int).GCD(nil, nil,
					new(big.Int).Add(prev, one), n); d.Cmp(one) != 0 && d.Cmp(n) != 0 &&
					new(big.Int).Mod(n, d).Sign() == 0 {
					return d, true
				}
				return nil, true // złożona, ale dzielnika nie udało się wyliczyć
			}
			prev.Set(b) // przechodzimy dalej: bj ← bj+1
		}

		// jeśli w ogóle nie napotkaliśmy wartości n−1  →  liczba złożona
		if !strong {
			return nil, true
		}
	}

	//------------------------------------------------------------
	// 4. Po wszystkich iteracjach: brak świadków złożoności
	//------------------------------------------------------------
	return nil, false // liczba prawdopodobnie pierwsza
}
