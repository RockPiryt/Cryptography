// Author: Paulina Kimak
package main

import (
	"flag"
	"log"
	"onetime/flagfunc"
	"onetime/helpers"
)



func main() {
	helpers.SetLogger()
	
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

	// orgFile := "files/org.txt"
	// plainFile := "files/plain.txt"

	// plainText, err := helpers.PrepareText(orgFile)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = helpers.SaveOutput(plainText, plainFile)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// fmt.Println("Text prepared and saved to plain.txt successfully.")
}