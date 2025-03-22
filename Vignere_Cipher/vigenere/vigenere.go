//Author: Paulina Kimak
package main

import (
	
	"flag"
	"fmt"
	"os"
	"vigenere/flagfunc"
	"vigenere/helpers"
)



func main() {
	//Set flags
	prepareFlag := flag.Bool("p", false, "przygotowanie tekstu jawnego do szyfrowania")
	encryptFlag := flag.Bool("e", false, "szyfrowanie")
	decryptFlag := flag.Bool("d", false, "deszyfrowanie")
	cryptAnalysisFlag := flag.Bool("k", false, "kryptoanaliza wyłącznie w oparciu o kryptogram")

	flag.Parse()

	// Check flags
	operationFlags := []*bool{prepareFlag,encryptFlag, decryptFlag, cryptAnalysisFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)


	if operationCount != 1 {
		fmt.Println("Błąd: Musisz wybrać dokładnie jedną operację (-p, -e, -d lub -k).")
		os.Exit(1)
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
		fmt.Println("Błąd: Nie wybrano poprawnej operacji (-p, -e, -d, -k).")
		os.Exit(1)
	}

	flagfunc.ExecuteCipher(operation)
}