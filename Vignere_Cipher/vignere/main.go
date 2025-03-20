package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"vignere/flagfunc"
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

	// Tests
	inputText := "files/org.txt"
	outputFile := "files/plain.txt"
	inputKey := "files/key.txt"
	outputText := "files/crypto.txt"
	preparedText := "files/plain.txt"

	// preparedText, err := helpers.PrepareText(inputFile)
	// if err != nil {
	// 	log.Printf("błąd przy czyszczeniu tekstu: %v", err)
	// }
	// fmt.Println(preparedText)

	// err = helpers.SaveOutput(preparedText, outputFile)
	// if err != nil {
	// 	log.Printf("błąd przy zapisie tekstu: %v", err)
	// }	

	
	// Prepare the text for encryption.
	// err := flagfunc.CreatePlainFile(inputText, outputFile)
	// if err != nil {
	// 	log.Printf("błąd przy przygotowywaniu tekstu: %v", err)
	// }
	
	
	// key, err := helpers.ValidateKey(inputKey)
	// if err != nil {
	// 	log.Printf("nie udało się zwalidować klucza %v", err)
	// }
	// fmt.Println("Zwalidowany Klucz: ", key)


	// numKey,err:=helpers.ConverseKey(key)
	// if err != nil {
	// 	log.Printf("nie udało się przekonwertować klucza %v", err)
	// }
	// fmt.Println("Skonwersowany Klucz  do liczb: ", numKey)


	// numKey, err := flagfunc.GetKey(inputKey)
	// if err != nil {
	// 	log.Printf("nie udało się przekonwertować klucza %v", err)
	// }
	// fmt.Println("Skonwersowany Klucz  do liczb: ", numKey)

	err := flagfunc.CreatePlainFile(inputText, outputFile)
		if err != nil {
			fmt.Printf("błąd przy przygotowywaniu tekstu: %v", err)
		}

	inputText, err = flagfunc.EncodeText(preparedText, inputKey, outputText)
	if err != nil {
		log.Printf("nie udało się odczytać poprawnego tekstu %v", err)
	}
	fmt.Println("Odczytany tekst: ", inputText)

}