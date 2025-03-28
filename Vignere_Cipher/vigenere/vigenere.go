// Author: Paulina Kimak
package main

import (
	"flag"
	"log"
	"vigenere/flagfunc"
	"vigenere/helpers"
)



func main() {
	//Set flags
	prepareFlag := flag.Bool("p", false, "prepare plaintext for encryption")
	encryptFlag := flag.Bool("e", false, "encrypt the plaintext")
	decryptFlag := flag.Bool("d", false, "decrypt the ciphertext")
	cryptAnalysisFlag := flag.Bool("k", false, "perform cryptanalysis based only on ciphertext")

	flag.Parse()

	// Check flags
	operationFlags := []*bool{prepareFlag,encryptFlag, decryptFlag, cryptAnalysisFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)


	if operationCount != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -p, -e, -d or -k.")
	}

	// Determine the operation
	var operation string
	switch {
	case *prepareFlag:
		operation = "p"
	case *encryptFlag:
		operation = "e"
	case *decryptFlag:
		operation = "d"
	case *cryptAnalysisFlag:
		operation = "k"
	default:
		log.Fatalf("Error: Invalid operation selected.")
	}

	err := flagfunc.ExecuteCipher(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}
}