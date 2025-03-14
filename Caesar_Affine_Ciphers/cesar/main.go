package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"cesar/datafile"
	"cesar/cryptofunc"
	"log"
	)

//Author: Paulina Kimak

func main() {
	
	lines, err := datafile.GetText("plain.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(lines)


	caesarFlag := flag.Bool("c", false, "szyfr cezara")
	affineFlag := flag.Bool("a", false, "szyfr afiniczny")
	encryptFlag := flag.Bool("e", false, "szyfrowanie")
	decryptFlag := flag.Bool("d", false, "deszyfrowanie")
	cryptoAnalizeFlag := flag.Bool("d", false, "deszyfrowanie")

	flag.Parse()

	// Check flags
	if !(*caesarFlag) {
		fmt.Println("Należy wybrać tryb szyfru za pomocą flagi -c")
		return
	}

	if *encryptFlag && *decryptFlag {
		fmt.Println("Nie można jednocześnie wybrać operacji szyfrowania i deszyfrowania")
		return
	}

	// Read plain.txt
	textLines, err := datafile.GetText("plain.txt")
	if err != nil {
		fmt.Println("Błąd przy odczycie pliku:", err)
		return
	}

	plaintext := strings.Join(textLines, "\n")

	// Read key.txt
	keyLines, err := datafile.GetText("key.txt")
	if err != nil {
		fmt.Println("Błąd przy odczycie klucza:", err)
		return
	}
	// Zakładając, że klucz jest jedną liczbą w pliku
	key, err := strconv.Atoi(keyLines[0])
	if err != nil {
		fmt.Println("Błąd przy konwersji klucza:", err)
		return
	}

	// Określamy, czy wykonujemy szyfrowanie czy deszyfrowanie
	encrypt := *encryptFlag
	if *decryptFlag {
		encrypt = false
	}

	// Wykonanie szyfrowania lub deszyfrowania
	result := CaesarCipher(plaintext, key, encrypt)

	// Zapisanie wyniku do pliku
	outputFile := "result.txt"
	err = os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		fmt.Println("Błąd przy zapisywaniu wyniku:", err)
		return
	}

	fmt.Printf("Zakończono operację. Wynik zapisano w %s.\n", outputFile)
}


