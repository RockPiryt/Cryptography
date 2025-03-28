// Author: Paulina Kimak
package main

import (
	"flag"
	"fmt"
	"log"
	"onetime/flagfunc"
	"onetime/helpers"
)

const (
	orgFile       = "files/org.txt"
	plainFile     = "files/plain.txt"
	keyFile       = "files/key.txt"
	keyOutputFile = "files/key-found.txt"
	cryptoFile    = "files/crypto.txt"
	decryptedFile = "files/decrypt.txt"
	keyFoundFile  = "files/key-found.txt"
)

func main() {
	helpers.SetLogger()

	//Set flags
	prepareFlag := flag.Bool("p", false, "prepare plaintext for encryption")
	encryptFlag := flag.Bool("e", false, "encrypt the plaintext")
	cryptAnalysisFlag := flag.Bool("k", false, "perform cryptanalysis based only on ciphertext")

	flag.Parse()

	// Check flags
	operationFlags := []*bool{prepareFlag, encryptFlag, cryptAnalysisFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)

	if operationCount != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -p, -e or -k.")
	}

	// Determine the operation
	var operation string
	switch {
	case *prepareFlag:
		operation = "p"
	case *encryptFlag:
		operation = "e"
	case *cryptAnalysisFlag:
		operation = "k"
	default:
		log.Fatalf("Error: Invalid operation selected.")
	}

	err := flagfunc.ExecuteCipher(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}

	// Tests
	cryptogram, err := flagfunc.XORCipher(plainFile, keyFile, cryptoFile)
	if err != nil {
		log.Fatalf("XORCipher function error: %v", err)
	}

	fmt.Println("Cryptogram (Hex):", cryptogram)

}
