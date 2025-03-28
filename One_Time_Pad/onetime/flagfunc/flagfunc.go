// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"os"

	"onetime/helpers"
)

const (
	orgFile       = "files/org.txt"
	plainFile     = "files/plain.txt"
	keyFile       = "files/key.txt"
	keyOutputFile = "files/key-found.txt"
	cryptoFile    = "files/crypto.txt"
	decryptedFile = "files/decrypt.txt"
	keyFoundFile  = "files/key-found.txt"
)

func ExecuteCipher(operation string) error {
	switch operation {
	case "p":
		// Prepare the text for encryption and save it to plain.txt
		err := CreatePlainFile(orgFile, plainFile)
		if err != nil {
			return fmt.Errorf("error during text preparation %v", err)
		}
		log.Println("[INFO] plain.txt has been successfully created.")
		return nil

	case "e":
		// Ensure plain.txt exists, if not, create it
		if _, err := os.Stat(plainFile); os.IsNotExist(err) {
			if err := CreatePlainFile(orgFile, plainFile); err != nil {
				return fmt.Errorf("error creating plain.txt automatically: %v", err)
			}
			log.Println("[INFO] plain.txt not found. It was automatically created using -p.")
		}

		_, err := EncryptXOR(plainFile, keyFile, cryptoFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}

		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil

	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		// BrakeCipher(cryptoFile, decryptedFile, keyOutputFile)

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
	return nil
}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile(inputFile string, outputFile string) error {
	maxLines := 15
	plainText, err := helpers.PrepareText(inputFile, maxLines)
	if err != nil {
		return fmt.Errorf("error while cleaning the text: %v", err)
	}

	err = helpers.SaveOutput(plainText, outputFile)
	if err != nil {
		return fmt.Errorf("error while saving the text: %v", err)
	}

	return nil
}

// XOR Cipher function to encrypt the plainText using the key
func EncryptXOR(plainFile, keyFile, cryptoFile string) (string, error) {
	plainText, err := helpers.GetPreparedText(plainFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać plain tekstu")
	}

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}

	log.Printf("Plain text: %s\n", plainText)
	log.Printf("Key: %s\n", key)

	// Convert the text and key to hexadecimal representations
	plainHex := helpers.TextToHex(plainText)
	keyHex := helpers.TextToHex(key)

	// Convert to bytes
	plainBytes, err := helpers.HexToBytes(plainHex)
	if err != nil {
		return "", fmt.Errorf("error converting text to bytes: %v", err)
	}

	keyBytes, err := helpers.HexToBytes(keyHex)
	if err != nil {
		return "", fmt.Errorf("error converting key to bytes: %v", err)
	}

	// Generate the cryptogram in hexadecimal format
	var cryptogram []byte
	for i := 0; i < len(plainBytes); i++ {
		cryptogram = append(cryptogram, plainBytes[i]^keyBytes[i%len(keyBytes)])
	}

	var cryptogramHex string
	for _, byteVal := range cryptogram {
		cryptogramHex += fmt.Sprintf("%02X", byteVal)
	}

	// Save the decrypted text to crypto.txt
	err = helpers.SaveOutput(cryptogramHex, cryptoFile)
	if err != nil {
		return "", fmt.Errorf("error saving cryptogram: %v", err)
	}

	return cryptogramHex, nil
}

// Analyze the ciphertext and try to guess spaces and letters
func AnalyzeXOR(cipherText []string) string {
	numLines := len(cipherText)
	numCols := len(cipherText[0])

	// Placeholder for the decrypted message
	decryptedMessage := make([]rune, numCols)

	// Guess key and spaces using XOR results
	key := make([]byte, numCols)

	// Iterate over columns (i.e., positions in the lines)
	for i := 0; i < numCols; i++ {
		// For each column, compare multiple lines
		possibleSpaces := []bool{}
		for line1 := 0; line1 < numLines-1; line1++ {
			for line2 := line1 + 1; line2 < numLines; line2++ {
				// XOR the lines
				line1Bytes := []byte(cipherText[line1])
				line2Bytes := []byte(cipherText[line2])
				xorResult := helpers.XORBytes(line1Bytes, line2Bytes)

				// Check the first bits to determine if it's a space (0x20 in ASCII)
				if xorResult[i] == 0x20 {
					possibleSpaces = append(possibleSpaces, true)
				} else {
					possibleSpaces = append(possibleSpaces, false)
				}
			}
		}

		// If most results indicate a space, set it to space
		if len(possibleSpaces) > 0 && possibleSpaces[0] {
			key[i] = 0x20 // Space in the key
			decryptedMessage[i] = '_'
		} else {
			// Assuming other characters are lowercase letters
			decryptedMessage[i] = 'a' // placeholder for the letter
		}
	}

	return string(decryptedMessage)
}

/////////////////////////////////Notes/////////////////////////////////

// // decryptVigenereSimple decrypts the given cryptoFile using the Vigenère cipher with the provided key.
// func DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile string) (string, error) {
// 	// cryptoText, err := helpers.GetPreparedText(cryptoFile)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("nie udało się odczytać crypto tekstu")
// 	// }
// 	// //fmt.Printf("Crypto Tekst: %s\n", cryptoText)
// 	// //fmt.Printf("Klucz: %s\n", keyFile)
// 	// key, err := helpers.GetPreparedKey(keyFile)
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("nie udało się odczytać klucza")
// 	// }
// 	// //fmt.Printf("Klucz: %s\n", key)
// 	// if len(cryptoText) == 0 || len(key) == 0 {
// 	// 	return "", fmt.Errorf("input plainText, key, or ALPHABET cannot be empty")
// 	// }
// 	// keyLength := len(key)

// 	// //fmt.Printf("Odszyfrowany tekst: %s\n", string(result))

// 	// // Save the decrypted text to decrypt.txt
// 	// err = helpers.SaveOutput(string(result), decryptedFile)
// 	// if err != nil {
// 	// 	log.Printf("błąd przy zapisie tekstu: %v", err)
// 	// 	return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
// 	// }

// 	return "string(result)", nil
// }

// //------------------------------------------------------------CryptoAnalysis------------------------------------------------------------
// // Function finds the key and decrypts the text
// func CryptoAnalysis(message string) []string {
// 	// seqs := findRepeats(message)
// 	// //fmt.Printf("Sequences: %v\n", seqs)
// 	// lengths := findKeyLengths(seqs)
// 	// //fmt.Printf("Lengths: %v\n", lengths)

// 	var allPossibleKeys []string
// 	// for _, length := range lengths {
// 	//     possibleKey := findKey(message, length)
// 	//     possibleKey = removeRepetitions(possibleKey)
// 	//     allPossibleKeys = append(allPossibleKeys, possibleKey)
// 	// }

// 	return allPossibleKeys
// }

// // // Function	find make frequency analysis based on key length
// // func findKey(message string, keyLength int) string {
// //     return ""
// // }

// func BrakeCipher(cryptoFile, decryptedFile, keyOutputFile string) {
// 	// Read the text from crypto.txt
// 	cryptoText, err := helpers.GetText(cryptoFile)
// 	if err != nil {
// 		log.Fatalf("Błąd odczytu crypto.txt: %v", err)
// 	}

// 	// Make analysis of the text from crypto.txt
// 	possibleKeys := CryptoAnalysis(cryptoText)
// 	//fmt.Println("\nMożliwe klucze:")
// 	// for _, k := range possibleKeys {
// 	// 	fmt.Println(" -", k)
// 	// }

// 	if len(possibleKeys) == 0 {
// 		fmt.Println("Nie znaleziono żadnych kluczy!")
// 		return
// 	}

// 	// Get the best key and save it to key-found.txt
// 	bestKey := possibleKeys[0]
// 	err = helpers.SaveOutput(bestKey, keyOutputFile)
// 	if err != nil {
// 		fmt.Printf("błąd przy zapisie tekstu: %v", err)
// 	}

// 	// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
// 	_, err = DecryptVigenereSimple(cryptoFile, keyOutputFile, decryptedFile)
// 	if err != nil {
// 		log.Fatalf("Błąd przy deszyfrowaniu: %v", err)
// 	}
// 	//fmt.Println("[INFO] Odszyfrowany tekst (decrypted.txt):", decrypted)
// }

// // Function to analyze XOR between two strings
// func AnalyzeXOR2(cryptoHex1, cryptoHex2 string) {
// 	// Convert to bytes
// 	bytes1, err := helpers.HexToBytes(cryptoHex1)
// 	if err != nil {
// 		fmt.Println("Error converting the first string:", err)
// 		return
// 	}
// 	bytes2, err := helpers.HexToBytes(cryptoHex2)
// 	if err != nil {
// 		fmt.Println("Error converting the second string:", err)
// 		return
// 	}

// 	for i := 0; i < len(bytes1); i++ {
// 		// XOR operation
// 		xorResult := bytes1[i] ^ bytes2[i]

// 		// Check if XOR starts with 010 (space)
// 		if xorResult == 0x20 { // Space is 0x20 in ASCII
// 			fmt.Printf("Found a space in the pair: '%s' and '%s'\n", string(bytes1[i]), string(bytes2[i]))
// 		}
// 	}
// }
