//Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"strings"

	"vignere/helpers"
)

// Generic function to handle Vignere cipher
func ExecuteCipher(operation string) {
	switch operation {
	case "p":
		// Prepare the text for encryption and save it to plain.txt
		orgFile := "files/org.txt"
		plainFile := "files/plain.txt"
		err := CreatePlainFile(orgFile, plainFile)
		if err != nil {
			log.Printf("Błąd: nie udało się odczytać klucza: %v", err)
		}
	case "e":
		// Encode the text from plain.txt using the key from key.txt and saves the result to crypto.txt
		plainFile := "files/plain.txt"
		keyFile := "files/key.txt"
		cryptoFile := "files/crypto.txt"
		encodedText, err := EncodeVignere(plainFile, keyFile, cryptoFile)
		if err != nil {
			log.Printf("nie udało się zaszyfrować tekstu %v", err)
		}
		fmt.Println("Zaszyfrowany tekst: ", encodedText)
	case "d":
		// Decode the text from crypto.txt using the key from key.txt and saves the result to decrypt.txt
		cryptoFile := "files/crypto.txt"
		keyFile := "files/key.txt"
		decryptedFile := "files/decrypt.txt"

		key, err := helpers.GetPreparedKey(keyFile)
		if err != nil {
			log.Printf("nie udało się odczytać klucza")
		}
		fmt.Printf("Klucz: %s\n", key)

		decodedText, err := DecryptVigenereSimple(cryptoFile, key, decryptedFile)
		if err != nil {
			log.Printf("nie udało się odszyfrować tekstu %v", err)
		}
		fmt.Println("Odszyfrowany tekst: ", decodedText)
		// DecodeVignere(cryptoFile, keyFile, decryptedFile) 
	case "k":
		// Make cryptanalysis of the text from crypto.txt and saves the result to decrypt.txt
		cryptoFile := "files/crypto.txt"
		plainOutputFile := "files/plain.txt"  
		keyOutputFile := "files/key-found.txt"
		decryptedFile := "files/decrypt.txt"
		CryptAnalysisVignere(cryptoFile, plainOutputFile, keyOutputFile, decryptedFile)

	default:
		fmt.Println("Nieobsługiwana operacja.")
		return
	}	

}

// Function to create a new file (plain.txt) containing prepared text for encryption.
func CreatePlainFile(inputFile string, outputFile string) error {
	plainText, err := helpers.PrepareText(inputFile)
	if err != nil {
		log.Printf("błąd przy czyszczeniu tekstu: %v", err)
		return fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
	}

	err = helpers.SaveOutput(plainText, outputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
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

	fmt.Printf("Klucz: %s\n", key)
	fmt.Printf("Plain Tekst: %s\n", plainText)
	
	if len(plainText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, or key cannot be empty")
	}
	
	var result []rune

	for i, char := range plainText {
		index := strings.IndexRune(helpers.Alphabet, char)
		keyIndex := strings.IndexRune(helpers.Alphabet, rune(key[i % len(key)]))

		encryptedIndex := (index + keyIndex) % helpers.AlphabetLen
		result = append(result, rune(helpers.Alphabet[encryptedIndex]))
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
	fmt.Printf("Crypto Tekst: %s\n", cryptoText)

	key, err := helpers.GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}
	fmt.Printf("Klucz: %s\n", key)
	if len(cryptoText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, key, or helpers.Alphabet cannot be empty")
	}
	keyLength := len(key)
	var result []rune

	for i, char := range cryptoText {
		index := strings.IndexRune(helpers.Alphabet, char)
		if index == -1 {
			result = append(result, char) // Keep non-helpers.Alphabet characters unchanged
			continue
		}

		keyIndex := strings.IndexRune(helpers.Alphabet, rune(key[i%keyLength]))
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		decryptedIndex := (index - keyIndex + helpers.AlphabetLen) % helpers.AlphabetLen
		result = append(result, rune(helpers.Alphabet[decryptedIndex]))
	}

	fmt.Printf("Odszyfrowany tekst: %s\n", string(result))

	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(string(result), decryptedFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// // getReversedKey generates a reversed key for decrypting using the Wikipedia formula: K2(i) = [26 – K(i)] mod 26
// func getReversedKey(key string) (string, error) {
// 	if len(key) == 0  {
// 		return "", fmt.Errorf("key cannot be empty")
// 	}

// 	var reversedKey []rune

// 	for _, char := range key {
// 		keyIndex := strings.IndexRune(helpers.Alphabet, char)
// 		if keyIndex == -1 {
// 			return "", fmt.Errorf("invalid character in key")
// 		}

// 		reversedIndex := (alphabetLength - keyIndex) % alphabetLength
// 		reversedKey = append(reversedKey, rune(helpers.Alphabet[reversedIndex]))
// 	}

// 	return string(reversedKey), nil
// }

// // decryptReversedKey decrypts the plainText using the reversed key.
// func DecodeVignere(cryptoText, reversedKey string) (string, error) {
// 	return EncodeVignere(cryptoText, reversedKey)
// }


//------------------------------------------------------------Kryptoanaliza------------------------------------------------------------
// Function finds the key and decrypts the text
func CryptAnalysisVignere(cryptoFile, plainOutputFile, keyOutputFile, decryptedFile string) error{
	cryptoText, err := helpers.GetPreparedText(cryptoFile)
	if err != nil {
		return fmt.Errorf("nie udało się odczytać crypto tekstu")
	}
	fmt.Printf("Crypto Tekst: %s\n", cryptoText)


	repeatedSequences := findRepeatedSequences(cryptoText)
	var allDistances []int
	for _, distList := range repeatedSequences {
		allDistances = append(allDistances, distList...)
	}

	keyLength := estimateKeyLength(allDistances)
	fmt.Printf("Estimated key length: %d\n", keyLength)

	key := findKey(cryptoText, keyLength)
	fmt.Printf("Estimated key: %s\n", key)

	decryptedText, err := DecryptVigenereSimple(cryptoText, key, decryptedFile)
	if err != nil {
		return fmt.Errorf("nie udało się odszyfrować tekstu %v", err)
	}
	fmt.Printf("Decrypted Text: %s\n", decryptedText)

	
	// Save the decrypted text to decrypt.txt
	err = helpers.SaveOutput(decryptedText, plainOutputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	// Save founded key to key-found.txt
	err = helpers.SaveOutput(key, keyOutputFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return nil
	
}


// Function findRepeatedSequences finds repeating sequences and returns their distances.
func findRepeatedSequences(cryptoText string) map[string][]int {
	sequencePositions := make(map[string][]int)
	sequenceDistances := make(map[string][]int)

	for i := 0; i < len(cryptoText)-2; i++ {
		seq := cryptoText[i : i+3] // 3-letters sequences
		if positions, exists := sequencePositions[seq]; exists {
			sequencePositions[seq] = append(positions, i)
		} else {
			sequencePositions[seq] = []int{i}
		}
	}

	for seq, positions := range sequencePositions {
		if len(positions) > 1 {
			for j := 1; j < len(positions); j++ {
				distance := positions[j] - positions[j-1]
				sequenceDistances[seq] = append(sequenceDistances[seq], distance)
			}
		}
	}

	return sequenceDistances
}


// Function estimateKeyLength finds the probable length of the key
func estimateKeyLength(distances []int) int {
    result := distances[0]
    for _, d := range distances[1:] {
        result = helpers.Gcd(result, d)
    }
    return result
} 


// findKey determines the key using frequency analysis
func findKey(cryptoText string, keyLength int) string {
	key := ""

	for i := 0; i < keyLength; i++ {
		var subText []rune
		// Create a subtext from the cryptoText
		for j := i; j < len(cryptoText); j += keyLength {
			subText = append(subText, rune(cryptoText[j]))
		}

		// Calculate the frequency of each letter in the subtext.
		letterFrequencies := make(map[rune]int)
		for _, letter := range subText {
			letterFrequencies[letter]++
		}

		// Find	the most common letter in the subtext.
		mostCommonLetter := 'a'
		maxCount := 0
		for letter, count := range letterFrequencies {
			if count > maxCount {
				maxCount = count
				mostCommonLetter = letter
			}
		}

		// Find the best match for the most common letter in the subtext.
		bestMatch := 'e' // Deuault value
		minDiff := 1000  

		for letter, freq := range helpers.FreqMap {
			diff := helpers.Absolute(freq - letterFrequencies[mostCommonLetter])
			if diff < minDiff {
				minDiff = diff
				bestMatch = letter
			}
		}

		// Calculate the shift for the key.
		shift := (strings.IndexRune(helpers.Alphabet, mostCommonLetter) - strings.IndexRune(helpers.Alphabet, bestMatch) + helpers.AlphabetLen) % helpers.AlphabetLen
		key += string(helpers.Alphabet[shift])
	}

	return key
}
