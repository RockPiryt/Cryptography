// Author: Paulina Kimak
package flagfunc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"xorcipher/helpers"
)

// Chartype  to store the type of character in the text. Unknown, Space (0x20), Letter (A-Z or a-z).
type CharType int

const (
	orgFile                = "files/orig.txt"
	plainFile              = "files/plain.txt"
	keyFile                = "files/key.txt"
	cryptoFile             = "files/crypto.txt"
	decryptedFile          = "files/decrypt.txt"
	Unknown       CharType = iota
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
		// Encrypt the text from plain.txt using the key from key.txt and save the result to crypto.txt
		_, err := EncryptXORBytes(plainFile, keyFile, cryptoFile)
		if err != nil {
			log.Fatalf("Encryption failed: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil

	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		_, err := AnalyzeXORBytes(cryptoFile, decryptedFile)
		if err != nil {
			log.Fatalf("Analysis failed: %v", err)
		}

		log.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil

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


// XOR Cipher function to encrypt the plainText using the key, save cryptogram in binary format
func EncryptXOR(plainFile, keyFile, cryptoFile string) error {
	plainText, err := helpers.GetPreparedText(plainFile)
	if err != nil {
		return fmt.Errorf("error during reading plain text: %v", err)
	}

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return fmt.Errorf("error during reading key: %v", err)
	}

	log.Printf("Plain text: %s\n", plainText)
	log.Printf("Key: %s\n", key)

	err = helpers.SaveOutput(key, keyFile)
	if err != nil {
		return fmt.Errorf("error during saving key: %v", err)
	}

	keyBytes := []byte(key)
	lines := strings.Split(plainText, "\n")
	var encryptedBytes []byte

	for _, line := range lines {
		lineBytes := []byte(line)
		for i := 0; i < len(lineBytes); i++ {
			encryptedBytes = append(encryptedBytes, lineBytes[i]^keyBytes[i%len(keyBytes)])
		}
		encryptedBytes = append(encryptedBytes, '\n') 
	}

	// Save associated cryptogram in binary format
	err = os.WriteFile(cryptoFile, encryptedBytes, 0644)
	if err != nil {
		return fmt.Errorf("error during saving cryptogram: %v", err)
	}

	return nil
}

// Function to decrypt the text using XOR cipher, cryptotext in binary format
func AnalyzeXOR(cryptoFile string) (string, error) {
	raw, err := os.ReadFile(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("error during reading file %s: %v", cryptoFile, err)
	}

	// Split binary content by newline (\n)
	linesRaw := bytes.Split(bytes.TrimSpace(raw), []byte("\n"))
	numLines := len(linesRaw)
	if numLines < 2 {
		return "", errors.New("minimum 2 lines required")
	}

	lineLen := len(linesRaw[0])
	decoded := make([][]byte, numLines)
	for i, line := range linesRaw {
		if len(line) != lineLen {
			return "", fmt.Errorf("all lines must have the same length; line %d has %d bytes, line 1 has %d",
				i+1, len(line), lineLen)
		}
		decoded[i] = line
	}

	// Character type heuristics
	charType := make([][]CharType, numLines)
	for i := 0; i < numLines; i++ {
		charType[i] = make([]CharType, lineLen)
		for k := 0; k < lineLen; k++ {
			charType[i][k] = Unknown
		}
	}

	top3 := func(x byte) byte {
		return (x & 0xE0) >> 5
	}

	changed := true
	rounds := 0
	for changed {
		changed = false
		rounds++
		for i := 0; i < numLines; i++ {
			for j := i + 1; j < numLines; j++ {
				for k := 0; k < lineLen; k++ {
					x := decoded[i][k] ^ decoded[j][k]
					switch top3(x) {
					case 0, 1:
						if x == 0x00 {
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
						} else {
							if charType[i][k] == Space {
							} else if charType[i][k] == Unknown {
								charType[i][k] = Letter
								changed = true
							}
							if charType[j][k] == Space {
							} else if charType[j][k] == Unknown {
								charType[j][k] = Letter
								changed = true
							}
						}
					case 2, 3:
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
					}
				}
			}
		}
	}
	log.Printf("Propagation finished after %d rounds.\n", rounds)

	// Extract known key
	key := make([]byte, lineLen)
	knownKey := make([]bool, lineLen)
	for i := 0; i < numLines; i++ {
		for k := 0; k < lineLen; k++ {
			if charType[i][k] == Space && !knownKey[k] {
				key[k] = decoded[i][k] ^ 0x20
				knownKey[k] = true
			}
		}
	}

	// Decrypt
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
	decryptText := strings.Join(output, "\n")

	// Save output
	err = helpers.SaveOutput(decryptText, decryptedFile)
	if err != nil {
		return "", fmt.Errorf("error during saving decrypted text: %v", err)
	}

	return decryptText, nil
}

func CheckPlain(plainText string) error {
	err := helpers.FindColumnsWithoutSpaces(plainText)
	if err != nil {
		return fmt.Errorf("error during FindColumnsWithoutSpaces : %v", err)
	}

	err = helpers.PrintSpacePositions(plainText)
	if err != nil {
		return fmt.Errorf("error durig PrintSpacePositions: %v", err)
	}

	return nil
}

// Function to decrypt the text using XOR cipher, cryptotext in hex format
func AnalyzeXORHex(cryptoFile string) (string, error) {
	raw, err := os.ReadFile(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("error during reading file %s: %v", cryptoFile, err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	numLines := len(lines)
	if numLines < 2 {
		return "", errors.New("minimum 2 lines required")
	}

	// Decode hex format to []byte
	decoded := make([][]byte, numLines)
	lineLen := 0
	for i, line := range lines {
		bs, err := hex.DecodeString(line)
		if err != nil {
			return "", fmt.Errorf("error during encoding %d: %v", i+1, err)
		}
		if i == 0 {
			lineLen = len(bs)
		} else if len(bs) != lineLen {
			return "", fmt.Errorf("all lines should has the same lenght, line %d has %d bytes, but 1. line %d",
				i+1, len(bs), lineLen)
		}
		decoded[i] = bs
	}

	// Create a 2D slice to store the character types for each line and each character (Unknown, Space, Letter)
	charType := make([][]CharType, numLines)
	for i := 0; i < numLines; i++ {
		charType[i] = make([]CharType, lineLen)
		for k := 0; k < lineLen; k++ {
			charType[i][k] = Unknown
		}
	}

	// Three top bits of the byte
	top3 := func(x byte) byte {
		return (x & 0xE0) >> 5
	}
	
	// Propagation of the character types
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
	log.Printf("Propagation fished after %d rounds.\n", rounds)

	// Calculate key with known spaces
	key := make([]byte, lineLen)
	knownKey := make([]bool, lineLen)

	for i := 0; i < numLines; i++ {
		for k := 0; k < lineLen; k++ {
			if charType[i][k] == Space && !knownKey[k] {
				key[k] = decoded[i][k] ^ 0x20 // 0x20 = space
				knownKey[k] = true
			}
		}
	}

	// Decrypt the text using the key
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
	decryptText := strings.Join(output, "\n")

	// Save crypto text in hex format to decrypt.txt
	err = helpers.SaveOutput(decryptText, decryptedFile)
	if err != nil {
		return "", fmt.Errorf("error during saving decryptText: %v", err)
	}

	return decryptText, nil
}

// XOR Cipher function to encrypt the plainText using the key, save cryptogram in hex format
func EncryptXORHex(plainFile, keyFile, cryptoFile string) (string, error) {
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

	// Save prepared key to key.txt
	err = helpers.SaveOutput(key, keyFile)
	if err != nil {
		return "", fmt.Errorf("error during saving cryptogram: %v", err)
	}

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

// EncryptXORBytes:
//  1) wczytuje plik plaintext (zakładamy, że każda linia ma 64 bajty –
//     lub dopełniamy ją, jeśli jest krótsza)
//  2) wczytuje plik z kluczem
//  3) XOR‑uje każdą 64‑bajtową linię z kluczem (powtarzanym cyklicznie w razie potrzeby)
//  4) Zapisuje surowe bloki 64‑bajtowe (BEZ znaku '\n') do pliku `cryptoFile`.
func EncryptXORBytes(plainFile, keyFile, cryptoFile string) (string, error) {
    const BlockSize = 64

    // Wczytujemy plaintext w trybie tekstowym
    plainRaw, err := os.ReadFile(plainFile)
    if err != nil {
        return "", fmt.Errorf("error reading plain text: %v", err)
    }
    // Dzielimy na linie
    lines := strings.Split(string(plainRaw), "\n")

    // Wczytujemy klucz (jako ciąg znaków ASCII)
    keyRaw, err := os.ReadFile(keyFile)
    if err != nil {
        return "", fmt.Errorf("error reading key: %v", err)
    }
    // Jeśli w keyFile mamy np. znak nowej linii, warto go przyciąć
    key := strings.TrimSpace(string(keyRaw))
    keyBytes := []byte(key)

    // Bufor do zapisu szyfrogramu
    var cipherBuf bytes.Buffer

    // Przetwarzamy każdą linię
    for lineIndex, line := range lines {
        lineBytes := []byte(line)

        // Jeżeli linia ma mniej niż 64 bajty – dopełnijmy spacjami (lub innym paddingiem)
        if len(lineBytes) < BlockSize {
            padding := make([]byte, BlockSize-len(lineBytes))
            for i := range padding {
                padding[i] = ' ' // dopełniamy spacjami
            }
            lineBytes = append(lineBytes, padding...)
        } else if len(lineBytes) > BlockSize {
            // jeśli jest dłuższa, możemy np. uciąć do 64
            // (albo rzucić błąd, jeśli to niepożądane)
            lineBytes = lineBytes[:BlockSize]
        }

        // XOR‑ujemy 64‑bajtowy blok
        encryptedBlock := make([]byte, BlockSize)
        for i := 0; i < BlockSize; i++ {
            encryptedBlock[i] = lineBytes[i] ^ keyBytes[i%len(keyBytes)]
        }

        // Zapisujemy blok do bufora
        // UWAGA – nie dopisujemy '\n' na końcu
        if _, err := cipherBuf.Write(encryptedBlock); err != nil {
            return "", fmt.Errorf("error writing encrypted block (line %d): %v", lineIndex+1, err)
        }
    }

    // Zapisujemy surowe bajty (wiele bloków 64‑bajtowych) do pliku
    if err := os.WriteFile(cryptoFile, cipherBuf.Bytes(), 0o644); err != nil {
        return "", fmt.Errorf("error saving cryptogram: %v", err)
    }

    return "Encryption complete (raw 64‑byte blocks saved).", nil
}

// AnalyzeXORBytes:
//  1) wczytuje plik `cryptoFile` jako jeden ciąg bajtów (bez dzielenia na '\n')
//  2) dzieli go na bloki 64‑bajtowe
//  3) wykonuje heurystyczną analizę XOR na wielu liniach (każdy blok to "linia")
//  4) odgaduje spacje i litery, odtwarza klucz w znanych kolumnach
//  5) odszyfrowuje i zapisuje wynik do `decryptedFile`
func AnalyzeXORBytes(cryptoFile, decryptedFile string) (string, error) {
    const BlockSize = 64

    raw, err := os.ReadFile(cryptoFile)
    if err != nil {
        return "", fmt.Errorf("error reading ciphertext file %s: %v", cryptoFile, err)
    }

    // Sprawdźmy, czy rozmiar pliku to wielokrotność 64
    if len(raw)%BlockSize != 0 {
        return "", fmt.Errorf("ciphertext size %d is not multiple of %d bytes", len(raw), BlockSize)
    }

    // Policzmy liczbę bloków
    numBlocks := len(raw) / BlockSize
    log.Printf("There are %d blocks, each %d bytes.\n", numBlocks, BlockSize)

    // Przetwarzamy w 2D tablicę: decoded[i][k] = bajt k z linii i
    decoded := make([][]byte, numBlocks)
    for i := 0; i < numBlocks; i++ {
        // weź blok i (64 bajtów)
        start := i * BlockSize
        end := start + BlockSize
        block := raw[start:end]
        // skopiuj do nowej tablicy (niekonieczne, ale przejrzyście)
        decoded[i] = make([]byte, BlockSize)
        copy(decoded[i], block)
    }

    // charType[i][k] => Unknown / Space / Letter
    charType := make([][]CharType, numBlocks)
    for i := 0; i < numBlocks; i++ {
        charType[i] = make([]CharType, BlockSize)
        for k := 0; k < BlockSize; k++ {
            charType[i][k] = Unknown
        }
    }

    // Funkcja pomocnicza – trzy najwyższe bity
    top3 := func(x byte) byte {
        return (x & 0xE0) >> 5
    }

    // Główna pętla propagacji
    changed := true
    rounds := 0
    for changed {
        changed = false
        rounds++
        // wszystkie pary (i,j)
        for i := 0; i < numBlocks; i++ {
            for j := i + 1; j < numBlocks; j++ {
                for k := 0; k < BlockSize; k++ {
                    x := decoded[i][k] ^ decoded[j][k]
                    switch top3(x) {
                    case 0, 1:
                        // x == 0 => spacja⊕spacja lub litera⊕litera (ta sama)
                        // x != 0 => raczej litera⊕litera (inne)
                        if x == 0 {
                            // ten sam typ w (i,k) i (j,k)
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
                        } else {
                            // x != 0 => raczej litera⊕litera
                            if charType[i][k] == Unknown {
                                if charType[i][k] != Space {
                                    charType[i][k] = Letter
                                    changed = true
                                }
                            }
                            if charType[j][k] == Unknown {
                                if charType[j][k] != Space {
                                    charType[j][k] = Letter
                                    changed = true
                                }
                            }
                        }
                    case 2, 3:
                        // litera⊕spacja
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
                    }
                }
            }
        }
    }
    log.Printf("Propagation finished after %d rounds.\n", rounds)

    // Odtwarzamy klucz tam, gdzie w plaintext była spacja
    key := make([]byte, BlockSize)
    knownKey := make([]bool, BlockSize)

    for i := 0; i < numBlocks; i++ {
        for k := 0; k < BlockSize; k++ {
            if charType[i][k] == Space && !knownKey[k] {
                // c = p ^ k => k = c ^ p => p=0x20 (spacja)
                key[k] = decoded[i][k] ^ 0x20
                knownKey[k] = true
            }
        }
    }

    // Odszyfrowanie i zbudowanie czytelnego tekstu z '_' w miejscach niepewnych
    var output []string
    for i := 0; i < numBlocks; i++ {
        var sb strings.Builder
        for k := 0; k < BlockSize; k++ {
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

    decryptedText := strings.Join(output, "\n")

    // Zapisujemy efekt „odgadniętego” tekstu do pliku
    err = os.WriteFile(decryptedFile, []byte(decryptedText), 0o644)
    if err != nil {
        return "", fmt.Errorf("error writing decrypted file: %v", err)
    }

    return decryptedText, nil
}


