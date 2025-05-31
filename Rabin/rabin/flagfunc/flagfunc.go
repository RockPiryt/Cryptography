// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"os"
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
	fmt.Println("fermat")
	return nil
}

func RabinMillerTest(EntryFile string) error {
	log.Println("Rabin-Miller test start")

	data, err := os.ReadFile(EntryFile)
	if err != nil {
		return fmt.Errorf("could not read input file: %v", err)
	}
	fmt.Printf("inputs data: %s", data)

	return nil
}
