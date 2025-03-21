//Author: Paulina Kimak
package main

import (
	"flag"
	"fmt"
	"os"
	"vignere/flagfunc"
	"vignere/helpers"
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

	// Tests
	//orgFile := "files/org.txt"
	// plainFile := "files/plain.txt"
	// keyFile := "files/key.txt"
	// cryptoFile := "files/crypto.txt"
	// decryptedFile := "files/decrypt.txt"


	// plainFile, err := helpers.PrepareText(inputFile)
	// if err != nil {
	// 	log.Printf("błąd przy czyszczeniu tekstu: %v", err)
	// }
	// fmt.Println(plainFile)

	// err = helpers.SaveOutput(plainFile, plainFile)
	// if err != nil {
	// 	log.Printf("błąd przy zapisie tekstu: %v", err)
	// }	

	
	// Prepare the text for encryption.
	// err := flagfunc.CreatePlainFile(orgFile, plainFile)
	// if err != nil {
	// 	log.Printf("błąd przy przygotowywaniu tekstu: %v", err)
	// }
	
	
	// key, err := helpers.ValidateKey(keyFile)
	// if err != nil {
	// 	log.Printf("nie udało się zwalidować klucza %v", err)
	// }
	// fmt.Println("Zwalidowany Klucz: ", key)


	// numKey,err:=helpers.ConverseKey(key)
	// if err != nil {
	// 	log.Printf("nie udało się przekonwertować klucza %v", err)
	// }
	// fmt.Println("Skonwersowany Klucz  do liczb: ", numKey)


	// numKey, err := flagfunc.GetKey(keyFile)
	// if err != nil {
	// 	log.Printf("nie udało się przekonwertować klucza %v", err)
	// }
	// fmt.Println("Skonwersowany Klucz  do liczb: ", numKey)

	// err := flagfunc.CreatePlainFile(orgFile, plainFile)
	// 	if err != nil {
	// 		fmt.Printf("błąd przy przygotowywaniu tekstu: %v", err)
	// 	}

	// orgFile, err = flagfunc.EncodeText(plainFile, keyFile, outputText)
	// if err != nil {
	// 	log.Printf("nie udało się odczytać poprawnego tekstu %v", err)
	// }
	// fmt.Println("Odczytany tekst: ", orgFile)

	// encodedText,err := flagfunc.EncodeVignere(plainFile, keyFile, cryptoFile)
	// if err != nil {
	// 	fmt.Printf("błąd przy szyfrowaniu tekstu: %v", err)
	// }

	// fmt.Println("Zaszyfrowany tekst: ", encodedText)

	// decodedText,err := flagfunc.DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile)
	// if err != nil {
	// 	fmt.Printf("błąd przy odszyfrowaniu tekstu: %v", err)
	// }

	// fmt.Println("odszyfrowany tekst: ", decodedText)
	//----------------------------------------
	// cryptoText := "wblbxylhrwblwyh" 

	// // Zakładamy długość klucza na podstawie wcześniejszych analiz (np. 5)
	// estimatedKeyLength := 5
	// key := flagfunc.FindKey(cryptoText, estimatedKeyLength)

	// fmt.Printf("Estimated Key: %s\n", key)

}