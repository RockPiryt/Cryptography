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
		switch {
		case *encryptFlag:
			if *caesarFlag {
				cryptofunc.CaesarEncrypt("e") 
			} else {
				cryptofunc.AffineEncrypt("e") 
			}
	
		case *decryptFlag:
		
			if *caesarFlag {
				cryptofunc.CaesarDecrypt("d") 
			} else {
				cryptofunc.AffineDecrypt("d") 
			}
	
		case *explicitCryptAnalysisFlag:
			
			if *caesarFlag {
				cryptofunc.CaesarExplicitCryptAnalysis("j") 
			} else {
				cryptofunc.AffineExplicitCryptAnalysis("j") 
			}
	
		case *cryptAnalysisFlag:
			
			if *caesarFlag {
				cryptofunc.CaesarCryptAnalysis("k") 
			} else {
				cryptofunc.AffineCryptAnalysis("k") 
			}
	
		default:
			fmt.Println("Błąd: Nie wybrano poprawnej operacji (-e, -d, -j, -k).")
			os.Exit(1)
		}
	}

}

