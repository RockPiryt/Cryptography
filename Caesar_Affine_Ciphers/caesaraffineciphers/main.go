package main



import (
	"flag"
	"fmt"
	"os"
	"log"
	"caesaraffineciphers/helpers"
	"caesaraffineciphers/cryptofunc"
)

//Author: Paulina Kimak


func main() {
	//Set flags
	caesarFlag := flag.Bool("c", false, "szyfr cezara")
	affineFlag := flag.Bool("a", false, "szyfr afiniczny")

	encryptFlag := flag.Bool("e", false, "szyfrowanie")
	decryptFlag := flag.Bool("d", false, "deszyfrowanie")
	explicitCryptAnalysisFlag := flag.Bool("j", false, "kryptoanaliza z tekstem jawnym")
	cryptAnalysisFlag := flag.Bool("k", false, "kryptoanaliza wyłącznie w oparciu o kryptogram")

	flag.Parse()

	// Check flags
	cipherFlags := []*bool{caesarFlag, affineFlag}
	operationFlags := []*bool{encryptFlag, decryptFlag, explicitCryptAnalysisFlag, cryptAnalysisFlag}

	cipherCount := helpers.CountSelectedFlags(cipherFlags)
	operationCount := helpers.CountSelectedFlags(operationFlags)

	if cipherCount != 1 {
		fmt.Println("Błąd: Musisz wybrać dokładnie jeden rodzaj szyfru (-c lub -a).")
		os.Exit(1)
	}

	if operationCount != 1 {
		fmt.Println("Błąd: Musisz wybrać dokładnie jedną operację (-e, -d, -j lub -k).")
		os.Exit(1)
	}


	// Determine the cipher type
	var cipherType string
	if *caesarFlag {
		cipherType = "caesar"
	} else if *affineFlag {
		cipherType = "affine"
	} else {
		log.Fatal("Błąd: Nie wybrano poprawnego szyfru (-c dla Cezara, -a dla afinicznego).")
	}

	// Determine the operation
	var operation string
	switch {
	case *encryptFlag:
		operation = "e"
	case *decryptFlag:
		operation = "d"
	case *explicitCryptAnalysisFlag:
		operation = "j"
	case *cryptAnalysisFlag:
		operation = "k"
	default:
		fmt.Println("Błąd: Nie wybrano poprawnej operacji (-e, -d, -j, -k).")
		os.Exit(1)
	}

	// Execute the cipher operation
	cryptofunc.ExecuteCipher(cipherType, operation)

	// c, a := helpers.ValidateKey("files/key.txt", "affine")
	// fmt.Printf("Przesunięcie c=%d, współczynnik a=%d:", c, a)

	// plaintext := "AZURE RAY"
	// ciphertextfull := cryptofunc.AffineCipher(plaintext, a, c, "e")
	// fmt.Println("Zaszyfrowany tekst full:", ciphertextfull)

	// fmt.Printf("Zapisuję do pliku: %s\n", ciphertextfull)
	// helpers.SaveOutput(ciphertextfull, "output2.txt")


	// ciphertext := "pq" 
	// plaintext := "if"
	// // Znajdowanie klucza (a, c)
	// a, c := cryptofunc.FindAffineKey(ciphertext, plaintext)
	// fmt.Printf("Odgadnięty klucz z dobrego pliku: a = %d, c = %d\n", a, c)
}