// Author: Paulina Kimak
package main

import (
	"elgamal/flagfunc"
	"elgamal/helpers"
	"flag"
	"fmt"
	"log"
	"math/big"
)

const (
	plainFile      = "files/plain.txt"
	messageFile    = "files/message.txt"
	elgamalFile    = "files/elgamal.txt"
	privateKeyFile = "files/private.txt"
	publicKeyFile  = "files/publicKey.txt"
	cryptoFile     = "files/crypto.txt"
	decryptedFile  = "files/decrypt.txt"
	signatureFile  = "files/signature.txt"
	verifyFile     = "files/verify.txt"
)

func main() {
	//Set flags
	keysFlag := flag.Bool("k", false, "prepare public and private keys")
	encryptFlag := flag.Bool("e", false, "encrypt the message")
	decryptFlag := flag.Bool("d", false, "decrypt the ciphertext")
	signatureFlag := flag.Bool("s", false, "sign the message")
	verifyFlag := flag.Bool("v", false, "verification of signature")

	flag.Parse()

	// Check flags
	operationFlags := []*bool{keysFlag,encryptFlag, decryptFlag, signatureFlag, verifyFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)


	if operationCount != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -k, -e, -d, -s or v.")
	}

	// Determine the operation
	var operation string
	switch {
	case *keysFlag:
		operation = "k"
	case *encryptFlag:
		operation = "e"
	case *decryptFlag:
		operation = "d"
	case *signatureFlag:
		operation = "s"
	case *verifyFlag:
		operation = "v"
	default:
		log.Fatalf("Error: Invalid operation selected.")
	}

	err := flagfunc.ExecuteCipher(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}

	// Read first number p and generator g from file
	params, _ := helpers.ReadBigIntsFromFile(elgamalFile, 2)
	p, g := params[0], params[1]

	fmt.Printf("p=%d g=%d", p, g)

	

	// Write big num to file
	helpers.WriteBigIntsToFile(privateKeyFile, []*big.Int{p, g})
	
}