package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/bits"
)

// Function calculates the number of differing bits between two hexadecimal strings of equal length.
func countDifferentBits(hex1, hex2 string) (int, error) {
	// Decode hex strings into byte slices
	bytes1, err := hex.DecodeString(hex1)
	if err != nil {
		return 0, fmt.Errorf("error decoding hex1: %v", err)
	}
	bytes2, err := hex.DecodeString(hex2)
	if err != nil {
		return 0, fmt.Errorf("error decoding hex2: %v", err)
	}

	// Check if lengths match
	if len(bytes1) != len(bytes2) {
		return 0, fmt.Errorf("hex strings have different lengths")
	}

	// XOR each byte and count differing bits
	count := 0
	for i := 0; i < len(bytes1); i++ {
		xor := bytes1[i] ^ bytes2[i]
		count += bits.OnesCount8(xor)
	}

	return count, nil
}

func main() {
	hex1 := "f82923c0fda16eeb90ed012245cbee4a"
	hex2 := "f6df877b14aa1ee2c024bef98d1fb0c0"

	diff, err := countDifferentBits(hex1, hex2)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Number of differing bits (Hamming distance): %d\n", diff)
}
