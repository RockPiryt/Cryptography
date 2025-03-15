package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Author: Paulina Kimak

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

// Function to read text from txt file
func GetText(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = RemovePolishLetters(line)
		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if scanner.Err() != nil{
		return nil, scanner.Err()
	}
	
	return lines, nil
}

// Function to remove Polish letters from text
func RemovePolishLetters(input string) string {
	// Mapa polskich liter z diakrytykami na ich odpowiedniki w alfabecie łacińskim
	replacementMap := map[rune]rune{
		'ą': 'a', 'ć': 'c', 'ę': 'e', 'ł': 'l', 'ń': 'n', 'ó': 'o', 'ś': 's', 'ź': 'z', 'ż': 'z',
		'Ą': 'A', 'Ć': 'C', 'Ę': 'E', 'Ł': 'L', 'Ń': 'N', 'Ó': 'O', 'Ś': 'S', 'Ź': 'Z', 'Ż': 'Z',
	}

	// Iterate through each character in the input string and replace if it's a Polish letter
	var result strings.Builder
	for _, r := range input {
		// If the character is a Polish letter, replace it
		if repl, found := replacementMap[r]; found {
			result.WriteRune(repl)
		} else if unicode.IsLetter(r) || unicode.IsSpace(r) {
			// If it's a letter (ignoring Polish diacriticals) or a space, keep it
			result.WriteRune(r)
		}
	}

	return result.String()
}

func SaveOutput(result string, outputFile string) {
	// Check if the file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Błąd przy tworzeniu pliku: %v", err)
		}
		file.Close()
	}

	// Write the result to the file
	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		log.Fatalf("Błąd przy zapisywaniu wyniku: %v", err)
	}

	//fmt.Println("Zapisano wynik do pliku:", outputFile)
}


// ValidateKey reads and validates the keys for both Caesar and Affine ciphers from a file.
func ValidateKey(filePath string, cipherType string) (int, int) {
	// Read the key from the file.
	keyLines, err := GetText(filePath)
	if err != nil {
		fmt.Printf("Błąd przy odczycie pliku klucza: %v\n", err)
		return -1, -1
	}

	// Ensure the file contains exactly one line.
	if len(keyLines) != 1 {
		fmt.Printf("Błędny klucz: Plik klucza powinien zawierać tylko jedną linię. Znaleziono: %d\n", len(keyLines))
		return -1, -1
	}

	// Split the line into two space-separated numbers.
	parts := strings.Fields(keyLines[0])
	if len(parts) != 2 {
		fmt.Printf("Błędny klucz: Plik klucza musi zawierać dokładnie dwie liczby oddzielone spacją (np. '3 7'). Znaleziono: %s\n", keyLines[0])
		return -1, -1
	}

	// Convert the first number (always used for Caesar and Affine).
	c, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Błędny klucz Cezara: Musi być liczbą całkowitą. Znaleziono: %s\n", parts[0])
		return -1, -1
	}

	// Validate the Caesar cipher key (0-25).
	if c < 0 || c > 25 {
		fmt.Printf("Błędny klucz Cezara: Klucz musi być liczbą z zakresu 0-25. Znaleziono: %d\n", c)
		return -1, -1
	}

	// If Caesar cipher is used, return only the first key, ignoring the second.
	if cipherType == "caesar" {
		return c, -1
	}

	// If the cipherType is not "caesar" or "affine", report an error.
	if cipherType != "affine" {
		fmt.Printf("Błędny typ szyfru: Oczekiwano 'caesar' lub 'affine', znaleziono: %s\n", cipherType)
		return -1, -1
	}

	// Convert the second number (used only for Affine cipher).
	a, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Printf("Błędny klucz afiniczny: Musi być liczbą całkowitą. Znaleziono: %s\n", parts[1])
		return -1, -1
	}

	// Validate that 'a' is coprime with 26 and is in the range [1, 25].
	m := 26
	if a < 1 || a > 25 {
		fmt.Printf("Błędny klucz afiniczny: 'a' musi być liczbą z zakresu 1-25. Znaleziono: %d\n", a)
		return -1, -1
	}

	// Check that 'a' is coprime with 26 (gcd(a, 26) == 1).
	gcd, _, _ := ExtendedGCD(a, m)
	if gcd != 1 {
		fmt.Printf("Błędny klucz afiniczny: Współczynnik 'a' musi być względnie pierwsza z 26. Znaleziono: %d\n", a)
		return -1, -1
	}

	// Return both values: c (for Caesar) and a (for Affine).
	return c, a
}

// Extended Euclidean algorithm
func ExtendedGCD(a, b int) (int, int, int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, x1, y1 := ExtendedGCD(b, a%b)
	x := y1
	y := x1 - (a/b)*y1
	return gcd, x, y
}

// ModInverseExtended calculates the modular inverse of a mod m using the extended Euclidean algorithm.
// If the modular inverse does not exist, it returns an error.
func ModInverseExtended(a, m int) (int, error) {
	gcd, x, _ := ExtendedGCD(a, m)
	if gcd != 1 {
		return 0, fmt.Errorf("brak modularnej odwrotności dla a=%d (mod %d)", a, m)
	}
	return (x%m + m) % m, nil // Ensure non-negative result
}