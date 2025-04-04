// Author: Paulina Kimak
package flagfunc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"onetime/helpers"
)
// CharType – rozpoznajemy: Unknown (jeszcze nie wiemy), Space (0x20), Letter (A-Z lub a-z).
type CharType int

const (
	orgFile       = "files/org.txt"
	plainFile     = "files/plain.txt"
	keyFile       = "files/key.txt"
	keyOutputFile = "files/key-found.txt"
	cryptoFile    = "files/crypto.txt"
	decryptedFile = "files/decrypt.txt"
	keyFoundFile  = "files/key-found.txt"
	Unknown CharType = iota
	Space
	Letter
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

// AnalyzeXOR – kryptoanaliza XOR z założeniem, że plaintext ma tylko spacje i litery (A-Z, a-z).
func AnalyzeXOR(cryptoFile string) (string, error) {
	// 1. Wczytanie pliku
	raw, err := os.ReadFile(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się wczytać pliku %s: %v", cryptoFile, err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	numLines := len(lines)
	if numLines < 2 {
		return "", errors.New("potrzebne co najmniej 2 linie do analizy")
	}

	// 2. Dekodowanie hex -> []byte
	decoded := make([][]byte, numLines)
	lineLen := 0
	for i, line := range lines {
		bs, err := hex.DecodeString(line)
		if err != nil {
			return "", fmt.Errorf("błąd dekodowania w linii %d: %v", i+1, err)
		}
		if i == 0 {
			lineLen = len(bs)
		} else if len(bs) != lineLen {
			return "", fmt.Errorf("wszystkie linie muszą mieć tę samą długość; linia %d ma %d bajtów, a 1. linia %d",
				i+1, len(bs), lineLen)
		}
		decoded[i] = bs
	}

	// 3. charType – czy to (Unknown, Space, Letter) w każdym miejscu
	charType := make([][]CharType, numLines)
	for i := 0; i < numLines; i++ {
		charType[i] = make([]CharType, lineLen)
		for k := 0; k < lineLen; k++ {
			charType[i][k] = Unknown
		}
	}

	// Pomocnicza funkcja: najwyższe 3 bity
	top3 := func(x byte) byte {
		return (x & 0xE0) >> 5
	}

	// 4. Iteracyjna propagacja
	changed := true
	rounds := 0
	for changed {
		changed = false
		rounds++

		// Wszystkie pary (i, j)
		for i := 0; i < numLines; i++ {
			for j := i + 1; j < numLines; j++ {
				for k := 0; k < lineLen; k++ {
					x := decoded[i][k] ^ decoded[j][k]
					switch top3(x) {
					case 0, 1:
						// Może to być litera⊕litera (różne lub te same), albo spacja⊕spacja (tylko x=0).
						// - jeśli x == 0 => obie takie same: (spacja-spacja) lub (litera-litera).
						if x == 0x00 {
							// jeżeli w (i, k) już wiemy: letter => j też letter,
							// jeżeli w (i, k) już wiemy: space => j też space.
							if charType[i][k] == Space && charType[j][k] == Unknown {
								charType[j][k] = Space
								changed = true
							} else if charType[i][k] == Letter && charType[j][k] == Unknown {
								charType[j][k] = Letter
								changed = true
							} else if charType[j][k] == Space && charType[i][k] == Unknown {
								charType[i][k] = Space
								changed = true
							} else if charType[j][k] == Letter && charType[i][k] == Unknown {
								charType[i][k] = Letter
								changed = true
							}
							// Jeśli obie Unknown, nic nie wnioskujemy.

						} else {
							// x != 0, top3(x) <= 1 => wskazuje raczej litera⊕litera (inne litery),
							// bo spacja⊕spacja = 0.
							// Jeżeli ktoś jest Space -> sprzeczność (teoretycznie).
							if charType[i][k] == Space {
								// Kolizja heurystyki; zostawiamy – w pełnej analityce można by to oznaczyć jako błąd
							} else if charType[i][k] == Unknown {
								charType[i][k] = Letter
								changed = true
							}
							if charType[j][k] == Space {
								// Kolizja
							} else if charType[j][k] == Unknown {
								charType[j][k] = Letter
								changed = true
							}
						}
					case 2, 3:
						// litera⊕spacja
						// => (i, k) = letter, (j, k) = space LUB (i, k) = space, (j, k) = letter
						// jeśli już wiemy jedną stronę -> determinujemy drugą
						if charType[i][k] == Letter && charType[j][k] != Space {
							charType[j][k] = Space
							changed = true
						} else if charType[i][k] == Space && charType[j][k] != Letter {
							charType[j][k] = Letter
							changed = true
						} else if charType[j][k] == Letter && charType[i][k] != Space {
							charType[i][k] = Space
							changed = true
						} else if charType[j][k] == Space && charType[i][k] != Letter {
							charType[i][k] = Letter
							changed = true
						}
						// Jeśli obie Unknown, w tej iteracji nic nie rozstrzygamy –
						// być może inna para powie, która jest spacja, która litera
					}
				}
			}
		}
	}
	log.Printf("Propagacja zakończona po %d rundach.\n", rounds)

	// 5. Wyliczamy klucz tam, gdzie pewnie Space
	key := make([]byte, lineLen)
	knownKey := make([]bool, lineLen)

	for i := 0; i < numLines; i++ {
		for k := 0; k < lineLen; k++ {
			if charType[i][k] == Space && !knownKey[k] {
				key[k] = decoded[i][k] ^ 0x20 // 0x20 = spacja
				knownKey[k] = true
			}
		}
	}

	// 6. Deszyfrujemy
	var output []string
	for i := 0; i < numLines; i++ {
		var sb strings.Builder
		for k := 0; k < lineLen; k++ {
			if knownKey[k] {
				ch := decoded[i][k] ^ key[k]
				if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == ' ' {
					sb.WriteByte(ch)
				} else {
					sb.WriteByte('_')
				}
			} else {
				sb.WriteByte('_')
			}
		}
		output = append(output, sb.String())
	}

	return strings.Join(output, "\n"), nil
}
