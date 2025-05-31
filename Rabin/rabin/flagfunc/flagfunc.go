// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"

	"rabin/helpers"
)

const (
	EntryFile  = "files/wejscie.txt"
	OutputFile = "files/wyjscie.txt"
)

var MessageString string = "Bla"

func ExecuteCipher(operation string) error {
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
	lines, err := helpers.ReadData("files/wejscie.txt")
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
