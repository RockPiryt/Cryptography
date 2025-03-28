// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"math"
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

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"
var AlphabetLen = len(ALPHABET)

var frequency = map[rune]float64{
    'a': 8.2, 'b': 1.5, 'c': 2.8, 'd': 4.3, 'e': 13.0, 'f': 2.2,
    'g': 2.0, 'h': 6.1, 'i': 7.0, 'j': 0.15, 'k': 0.77, 'l': 4.0,
    'm': 2.4, 'n': 6.7, 'o': 7.5, 'p': 1.9, 'q': 0.095, 'r': 6.0,
    's': 6.3, 't': 9.1, 'u': 2.8, 'v': 0.98, 'w': 2.4, 'x': 0.15,
    'y': 2.0, 'z': 0.074,
}

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

		_, err := EncodeVignere(plainFile, keyFile, cryptoFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil
		
	case "d":
		// Decrypt crypto.txt using key.txt
		_, err := helpers.GetPreparedKey(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read the key: %v", err)
		}

		_, err = DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		fmt.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil 
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		BrakeCipher(cryptoFile, decryptedFile, keyOutputFile)

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}	
	return nil
}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile(inputFile string, outputFile string) error {
	plainText, err := helpers.PrepareText(inputFile)
	if err != nil {
		return fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
	}

	err = helpers.SaveOutput(plainText, outputFile)
	if err != nil {
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
}

// Function to encrypt the plainText using the Vigenère cipher with the provided key.
func EncodeVignere(plainFile, keyFile, cryptoFile string) (string, error) {
	plainText, err := helpers.GetPreparedText(plainFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać plain tekstu")
	}

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}

	//fmt.Printf("Klucz: %s\n", key)
	//fmt.Printf("Plain Tekst: %s\n", plainText)
	
	if len(plainText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, or key cannot be empty")
	}
	
	var result []rune

	for i, char := range plainText {
		index := strings.IndexRune(ALPHABET, char)
		keyIndex := strings.IndexRune(ALPHABET, rune(key[i % len(key)]))

		encryptedIndex := (index + keyIndex) % AlphabetLen
		result = append(result, rune(ALPHABET[encryptedIndex]))
	}

	// Save the decrypted text to crypto.txt
	err = helpers.SaveOutput(string(result), cryptoFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// decryptVigenereSimple decrypts the given cryptoFile using the Vigenère cipher with the provided key.
func DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile string) (string, error) {
	cryptoText, err := helpers.GetPreparedText(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać crypto tekstu")
	}
	//fmt.Printf("Crypto Tekst: %s\n", cryptoText)
	//fmt.Printf("Klucz: %s\n", keyFile)
	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}
	//fmt.Printf("Klucz: %s\n", key)
	if len(cryptoText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, key, or ALPHABET cannot be empty")
	}
	keyLength := len(key)
	var result []rune

	for i, char := range cryptoText {
		index := strings.IndexRune(ALPHABET, char)
		if index == -1 {
			result = append(result, char) // Keep non-ALPHABET characters unchanged
			continue
		}

		keyIndex := strings.IndexRune(ALPHABET, rune(key[i%keyLength]))
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		decryptedIndex := (index - keyIndex + AlphabetLen) % AlphabetLen
		result = append(result, rune(ALPHABET[decryptedIndex]))
	}

	//fmt.Printf("Odszyfrowany tekst: %s\n", string(result))

	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(string(result), decryptedFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}


//------------------------------------------------------------CryptoAnalysis------------------------------------------------------------
// Function finds the key and decrypts the text
func CryptoAnalysis(message string) []string {
	seqs := findRepeats(message)
	//fmt.Printf("Sequences: %v\n", seqs)
	lengths := findKeyLengths(seqs)
	//fmt.Printf("Lengths: %v\n", lengths)

	var allPossibleKeys []string
    for _, length := range lengths {
        possibleKey := findKey(message, length)
        possibleKey = removeRepetitions(possibleKey)
        allPossibleKeys = append(allPossibleKeys, possibleKey)
    }

	return allPossibleKeys
}


// Function findRepeatedSequences finds repeating sequences and returns the distances between them.
func findRepeats(cryptoText string) map[string][]int {	sequences := make(map[string][]int)
    seqLength := 4 

    // Find all sequences of length seqLength
    for seqBegin := 0; seqBegin < len(cryptoText)-seqLength; seqBegin++ {
        seq := cryptoText[seqBegin : seqBegin+seqLength]
        for i := seqBegin + seqLength; i < len(cryptoText)-seqLength; i++ {
            if cryptoText[i:i+seqLength] == seq {
                sequences[seq] = append(sequences[seq], i-seqBegin)
            }
        }
    }
    return sequences
}

// Function findKeyLengths finds the probable key lengths.
func findKeyLengths(sequences map[string][]int) []int {
    potentialKeyAccuracy := make(map[int]float64)
    for length := 2; length < 11; length++ {
        var counter, secondaryCounter float64
        for _, distances := range sequences {
            for _, dist := range distances {
                secondaryCounter++
                if dist%length == 0 {
                    counter++
                }
            }
        }
        // If no repeating sequences were found
        if secondaryCounter == 0 {
            fmt.Println("Brak danych do analizy (zbyt krótki tekst lub brak powtórzeń?).")
        }
        potentialKeyAccuracy[length] = counter / secondaryCounter
    }

    //fmt.Println("Potential Key Accuracy:", potentialKeyAccuracy)

    // Filter out the potential key lengths over 60% accuracy
    var potentialKeys []int
    for length, acc := range potentialKeyAccuracy {
        if acc > 0.60 {
            potentialKeys = append(potentialKeys, length)
        }
    }

    //fmt.Println("Chosen key lengths:", potentialKeys)
    return potentialKeys
}

// Function	find make frequency analysis based on key length
func findKey(message string, keyLength int) string {
    var key strings.Builder

    for i := 0; i < keyLength; i++ {
        // Prepare maps for each letter in key
        positionalDict := make(map[rune]map[rune]int)
        scoredDict := make(map[rune]float64)

        for _, letter := range ALPHABET {
            letterRune := rune(letter)
            positionalDict[letterRune] = make(map[rune]int)
            scoredDict[letterRune] = 0
            // Initialize counter
            for _, letter2 := range ALPHABET {
                positionalDict[letterRune][rune(letter2)] = 0
            }
        }

        // Check each shift.
        for _, letter := range ALPHABET {
            letterRune := rune(letter)
            idx := i
            for idx < len(message) {
                row := strings.IndexRune(ALPHABET, rune(message[idx]))
                col := strings.IndexRune(ALPHABET, letterRune)
                if row == -1 || col == -1 {
                    // If the character is not in the alphabet, skip it
                    idx += keyLength
                    continue
                }
                shifted := rune(ALPHABET[(row-col+26)%26])
                positionalDict[letterRune][shifted]++
                idx += keyLength
            }

            // Calculate the score for each letter
            score := 0.0
            for char, count := range positionalDict[letterRune] {
                freqVal, ok := frequency[char]
                if !ok {
                    freqVal = 0.0
                }
                score += float64(count) * freqVal
            }
            scoredDict[letterRune] = score
        }

        // Find the best letter for the key
        bestLetter := findMaxKey(scoredDict)
        key.WriteRune(bestLetter)
    }

    return key.String()
}

func removeRepetitions(k string) string {
    length := len(k)
    for sublen := 1; sublen <= length; sublen++ {
        if length%sublen == 0 {
            candidate := k[:sublen]
            repeated := strings.Repeat(candidate, length/sublen)
            if repeated == k {
                return candidate
            }
        }
    }
    return k
}

func findMaxKey(scoredDict map[rune]float64) rune {
    var maxRune rune
    maxVal := -math.MaxFloat64
    for k, v := range scoredDict {
        if v > maxVal {
            maxVal = v
            maxRune = k
        }
    }
    return maxRune
}

func BrakeCipher(cryptoFile, decryptedFile, keyOutputFile string){
	// Read the text from crypto.txt
	cryptoText, err := helpers.GetText(cryptoFile)
	if err != nil {
		log.Fatalf("Błąd odczytu crypto.txt: %v", err)
	}

	// Make analysis of the text from crypto.txt
	possibleKeys := CryptoAnalysis(cryptoText)
	//fmt.Println("\nMożliwe klucze:")
	// for _, k := range possibleKeys {
	// 	fmt.Println(" -", k)
	// }

	if len(possibleKeys) == 0 {
		fmt.Println("Nie znaleziono żadnych kluczy!")
		return
	}

	// Get the best key and save it to key-found.txt
	bestKey := possibleKeys[0]
	err = helpers.SaveOutput(bestKey, keyOutputFile)
	if err != nil {
		fmt.Printf("błąd przy zapisie tekstu: %v", err)
	}

	// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
	_, err = DecryptVigenereSimple(cryptoFile, keyOutputFile, decryptedFile)
	if err != nil {
		log.Fatalf("Błąd przy deszyfrowaniu: %v", err)
	}
	//fmt.Println("[INFO] Odszyfrowany tekst (decrypted.txt):", decrypted)
}
