// Author: Paulina Kimak
package main

import (
	"flag"
	"log"
	"xorcipher/flagfunc"
	"xorcipher/helpers"
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

	// // Tests crypto analysis
	// decryptedText, err := flagfunc.AnalyzeXOR("files/crypto.txt")
	// if err != nil {
	// 	log.Fatalf("Error during analysis: %v", err)
	// }
	// fmt.Println("\n\nPlain text after analysis:\n", decryptedText)

}
