package main

import (
	"flag"
	"fmt"
	"os"
	"vignere/helpers"

)

//Author: Paulina Kimak


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

}