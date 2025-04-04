// Author: Paulina Kimak
package flagfunc

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

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
		return "", fmt.Errorf("error during reading plain text: %v", err)
	}

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("error during reading key: %v", err)
	}

	log.Printf("Plain text: %s\n", plainText)
	log.Printf("Key: %s\n", key)
	
	// Convert plain text and key to byte slices.
	keyBytes := []byte(key)
	lines := strings.Split(plainText, "\n")
	var cryptogramLines []string

	for _, line := range lines {
		lineBytes := []byte(line)
		var encryptedLine []byte
		for i := 0; i < len(lineBytes); i++ {
			encryptedLine = append(encryptedLine, lineBytes[i]^keyBytes[i%len(keyBytes)])
		}

		// Convert encrypted line to hexadecimal representation.
		var encryptedHex string
		for _, b := range encryptedLine {
			encryptedHex += fmt.Sprintf("%02X", b)
		}

		cryptogramLines = append(cryptogramLines, encryptedHex)
	}

	// Join the encrypted lines into a single crypto text.
	cryptogramHex := strings.Join(cryptogramLines, "\n")

	// Save crypto text in hex format to crypto.txt
	err = helpers.SaveOutput(cryptogramHex, cryptoFile)
	if err != nil {
		return "", fmt.Errorf("error during saving cryptogram: %v", err)
	}

	return cryptogramHex, nil
}

// Analyze the ciphertext and try to guess spaces and letters
func AnalyzeXOR(cryptoFile string) (string, error) {
	cryptoText, err := helpers.GetPreparedText(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("error during reading crypto text: %v", err)
	}

	lines := strings.Split(cryptoText, "\n")
	numLines := len(lines)
	if numLines < 2 {
		return "", fmt.Errorf("minimum 2 lines to crypto analyze")
	}

	// Encode each line from hex to bytes
	var cryptoBinary [][]byte
	lineLength := 0
	for _, hexLine := range lines {
		bytesLine, err := hex.DecodeString(hexLine)
		if err != nil {
			return "", fmt.Errorf("encoded error hex: %v", err)
		}
		if lineLength == 0 {
			lineLength = len(bytesLine)
		} else if len(bytesLine) != lineLength {
			return "", fmt.Errorf("length of lines is not the same")
		}
		cryptoBinary = append(cryptoBinary, bytesLine)
	}

	// Create a 2D slice to store guesses for spaces
	spaceGuesses := make([][]bool, numLines)
	for i := range spaceGuesses {
		spaceGuesses[i] = make([]bool, lineLength)
	}

	// Make XOR of each line with every other line and look for positions that could contain spaces
	for i := 0; i < numLines; i++ {
		for j := i + 1; j < numLines; j++ {
			for k := 0; k < lineLength; k++ {
				x := cryptoBinary[i][k] ^ cryptoBinary[j][k]
				// Spacja XOR litera ASCII [a-z] daje wynik z 3-cim bitem ustawionym (czyli x & 0x20 == 0x20)
				// XOR space with ASCII letters gives a result with the 3rd bit set (x & 0x20 == 0x20 == 0010.000)
				if (x >= 0x41 && x <= 0x7A) || (x >= 0x01 && x <= 0x1F) {
					// Możliwe, że jedno z nich to spacja – zaznacz kandydata
					spaceGuesses[i][k] = true
					spaceGuesses[j][k] = true
				}
			}
		}
	}

	// Zgadywanie klucza tam, gdzie jesteśmy pewni spacji
	key := make([]byte, lineLength)
	knownKey := make([]bool, lineLength)

	for i := 0; i < numLines; i++ {
		for k := 0; k < lineLength; k++ {
			if spaceGuesses[i][k] {
				// Jeśli zakładamy, że to spacja (0x20), to klucz = crypto ^ 0x20
				guessedKeyByte := cryptoBinary[i][k] ^ 0x20
				// Zaznacz jako znany tylko jeśli jeszcze nie mamy klucza w tym miejscu
				if !knownKey[k] {
					key[k] = guessedKeyByte
					knownKey[k] = true
				}
			}
		}
	}

	// Decrypt each line using the guessed key
	var output []string
	for _, line := range cryptoBinary {
		var lineText string
		for i, b := range line {
			if knownKey[i] {
				ch := b ^ key[i]
				if ch >= 0x61 && ch <= 0x7A { // a-z
					lineText += string(ch)
				} else if ch == 0x20 {
					lineText += " "
				} else {
					lineText += "_" // doubtful character
				}
			} else {
				lineText += "_" //don't know the key byte
			}
		}
		output = append(output, lineText)
	}

	// Create lines
	result := strings.Join(output, "\n")
	return result, nil
}
