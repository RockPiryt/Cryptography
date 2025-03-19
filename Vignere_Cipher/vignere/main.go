package main



import (
	"flag"
	"fmt"
	"vignere/helpers"
)

//Author: Paulina Kimak


func main() {
	//Set flags
	caesarFlag := flag.Bool("p", false, "przygotowanie tekstu jawnego do szyfrowania")
	encryptFlag := flag.Bool("e", false, "szyfrowanie")
	decryptFlag := flag.Bool("d", false, "deszyfrowanie")
	cryptAnalysisFlag := flag.Bool("k", false, "kryptoanaliza wyłącznie w oparciu o kryptogram")

	flag.Parse()

	fmt.Printf("Wybrano %d flagi\n", helpers.CountSelectedFlags([]*bool{caesarFlag, encryptFlag, decryptFlag, cryptAnalysisFlag}))
}