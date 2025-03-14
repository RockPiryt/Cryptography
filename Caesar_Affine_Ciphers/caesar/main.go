package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"caesar/helpers"
	"caesar/cryptofunc"
	"log"
	)

//Author: Paulina Kimak


func main() {
	//Set flags
	caesarFlag := flag.Bool("c", false, "szyfr cezara")
	affineFlag := flag.Bool("a", false, "szyfr afiniczny")
	
	encryptFlag := flag.Bool("e", false, "szyfrowanie")
	decryptFlag := flag.Bool("d", false, "deszyfrowanie")
	explicitCryptAnalysisFlag := flag.Bool("j", false, "kryptoanaliza wyłącznie w oparciu o kryptogram")
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

	// Open a file.
	lines, err := helpers.GetText("plain.txt")
	if err != nil {
		log.Fatalf("Nie udało się odczytać pliku: %v", err)
	}
	fmt.Println(lines)

	// Read text.
	textLines, err := helpers.GetText("plain.txt")
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	plaintext := strings.Join(textLines, "\n")

	// Read a key.
	keyLines, err := helpers.GetText("key.txt")
	if err != nil {
		log.Fatalf("Nie udało się odczytać klucza : %v", err)
	}
	// Converse key to integer.
	key, err := strconv.Atoi(keyLines[0])
	if err != nil {
		log.Fatalf("Błąd przy konwersji klucza: %v", err)
	}

	// Caesar encrytion.
	result := cryptofunc.CaesarCipher(plaintext, key, *encryptFlag)

	// Save result to file.
	outputFile := "result.txt"
	
	// Sprawdzenie, czy plik istnieje
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Plik nie istnieje, więc go tworzymy
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Błąd przy tworzeniu pliku:", err)
			return
		}
		file.Close()
	}
	
	// Zapisanie wyniku do pliku
	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Błąd przy zapisywaniu wyniku:", err)
		return
	}

	fmt.Println("Zapisano wynik do pliku:", outputFile)
}


