package main



import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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


	// Read text.
	textLines, err := helpers.GetText("files/plain.txt")
	if err != nil {
		log.Fatalf("Błąd przy odczycie pliku: %v", err)
	}

	plaintext := strings.Join(textLines, "\n")
	fmt.Println(plaintext)
	
	// Read a key.
	keyLines, err := helpers.GetText("files/key.txt")
	if err != nil {
		log.Fatalf("Nie udało się odczytać klucza : %v", err)
	}
	
	// Converse key to integer.
	key, err := strconv.Atoi(keyLines[0])
	if err != nil {
		log.Fatalf("Błąd przy konwersji klucza: %v", err)
	}

	// Determine the cipher function to use
	var cipherFunc func(string, int, bool) string

	if *caesarFlag {
		cipherFunc = cryptofunc.CaesarCipher
	} else if *affineFlag {
		cipherFunc = cryptofunc.AffineCipher
	} else {
		log.Fatal("Błąd: Nie wybrano poprawnego szyfru (-c dla Cezara, -a dla afinicznego).")
	}

	// Determine the operation
	var operationFlag *bool

	switch {
	case *encryptFlag:
		operationFlag = encryptFlag
	case *decryptFlag:
		operationFlag = decryptFlag
	case *explicitCryptAnalysisFlag:
		operationFlag = explicitCryptAnalysisFlag
	case *cryptAnalysisFlag:
		operationFlag = cryptAnalysisFlag
	default:
		fmt.Println("Błąd: Nie wybrano poprawnej operacji (-e, -d, -j, -k).")
		os.Exit(1)
	}

	// Perform encryption or decryption
	result := cipherFunc(plaintext, key, *operationFlag)
	helpers.SaveOutput(result, "files/result.txt")

}
