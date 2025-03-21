package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// ALPHABET – dostępne znaki.
const ALPHABET = "abcdefghijklmnopqrstuvwxyz"

// frequency – częstości występowania liter w języku angielskim.
// (W razie potrzeby można je zmienić lub dopasować do innego języka.)
var frequency = map[rune]float64{
    'a': 8.2, 'b': 1.5, 'c': 2.8, 'd': 4.3, 'e': 13.0, 'f': 2.2,
    'g': 2.0, 'h': 6.1, 'i': 7.0, 'j': 0.15, 'k': 0.77, 'l': 4.0,
    'm': 2.4, 'n': 6.7, 'o': 7.5, 'p': 1.9, 'q': 0.095, 'r': 6.0,
    's': 6.3, 't': 9.1, 'u': 2.8, 'v': 0.98, 'w': 2.4, 'x': 0.15,
    'y': 2.0, 'z': 0.074,
}

func main() {
    // Przykładowy zaszyfrowany tekst
    message := `klt vugemck kxp zpvitex dj tav zxkegèii rmpavv les tcwd gogjmsirxu
ycfrxropflx rx dre mzqt. xhbj ztvsbfr jwel rw ile dvc p flhto dj txox
pw lhek pw tav taeigkimx. tav tgsbevq lmta klt vugemck kxp zxkegèii
rmpavv xw tarx ile vicexagrpnwt arw hxamzwimctc mcjokdeimog rfdyt myi
ziy xcibinmj(ehwufzrv xhtk xwi befgz sf mvbi ms be e zrope pprgnrkt)
enw klpx igwsgqamzsc aiec ft veycirxew zr ile vztwirmvbi.mf njmck a
miyac rtehdq kxp, awmca zw px lxrwi es efrv es myi trckptiid fvwhegx
rrs ms njis snep scge, myi kmgxeèvt giiyig ms myidvemzgpplr
lrqvetbeqpe. afatzek, zr ilil tehi, im zw ile dvc, cst myi rmpavv,
llivy tgsvbuih grrgxdkrtglxg smiickta, rrs wuvy wnwtxdw pve vfvgicmcc
gifxivth th tsapevkmkilr rw dre-mzqt taw jchxefj, mgvelgirxiov su
ahbtl rmpavvh erx vqeporvh.`

    // Uruchomienie analizy CryptoAnalysis – zwróci listę możliwych kluczy
    possibleKeys := CryptoAnalysis(message)

    fmt.Println("\nMożliwe klucze to:")
    for _, k := range possibleKeys {
        fmt.Println("-", k)
    }

    // Dla każdego klucza pokażmy też próbę odszyfrowania
    for _, k := range possibleKeys {
        fmt.Printf("\nPróba odszyfrowania kluczem: %s\n", k)
        decrypted := vigenere(message, k, "decrypt")
        fmt.Println("Odszyfrowany tekst:", decrypted)
    }
}

// CryptoAnalysis – główna funkcja wykonująca analizę CryptoAnalysis:
// 1. Usuwa zbędne znaki z tekstu.
// 2. Znajduje powtarzające się sekwencje i oblicza możliwe długości klucza.
// 3. Dla każdej długości generuje klucz i usuwa ewentualne powtórzenia w nim.
// 4. Zwraca listę potencjalnych kluczy.
func CryptoAnalysis(message string) []string {
    cleaned := sanitize(message)
    seqs := findRepeats(cleaned)
    lengths := findKeyLengths(seqs)

    var allPossibleKeys []string
    for _, length := range lengths {
        possibleKey := findKey(cleaned, length)
        possibleKey = removeRepetitions(possibleKey)
        allPossibleKeys = append(allPossibleKeys, possibleKey)
    }
    return allPossibleKeys
}

// sanitize – usuwa wszystkie znaki poza a–z i zamienia na małe litery.
func sanitize(msg string) string {
    re := regexp.MustCompile(`[^a-z]+`)
    lower := strings.ToLower(msg)
    return re.ReplaceAllString(lower, "")
}

// findRepeats – wyszukuje w tekście powtarzające się sekwencje (domyślnie 4-literowe)
// i zwraca mapę: { "sekwencja": [odległości między powtórzeniami], ... }
func findRepeats(message string) map[string][]int {
    sequences := make(map[string][]int)
    seqLength := 4 // tu można zmieniać długość n-gramów

    // Szukamy powtarzających się 4-literowych fragmentów
    for seqBegin := 0; seqBegin < len(message)-seqLength; seqBegin++ {
        seq := message[seqBegin : seqBegin+seqLength]
        for i := seqBegin + seqLength; i < len(message)-seqLength; i++ {
            if message[i:i+seqLength] == seq {
                sequences[seq] = append(sequences[seq], i-seqBegin)
            }
        }
    }
    return sequences
}

// findKeyLengths – dla każdej długości klucza (2..9) sprawdza, jak często dzieli ona
// odległości między powtórzeniami. Zwraca wszystkie długości, których wynik > 0.30.
func findKeyLengths(sequences map[string][]int) []int {
    potentialKeyAccuracy := make(map[int]float64)
    for length := 2; length < 10; length++ {
        var counter, secondaryCounter float64
        for _, distances := range sequences {
            for _, dist := range distances {
                secondaryCounter++
                if dist%length == 0 {
                    counter++
                }
            }
        }
        // Jeśli w ogóle nie znaleziono powtórek, rzucamy wyjątek
        if secondaryCounter == 0 {
            panic("Brak danych do analizy (zbyt krótki tekst lub brak powtórzeń?).")
        }
        potentialKeyAccuracy[length] = counter / secondaryCounter
    }

    fmt.Println("Potential Key Accuracy:", potentialKeyAccuracy)

    // Filtrujemy długości klucza powyżej progu 0.30
    var potentialKeys []int
    for length, acc := range potentialKeyAccuracy {
        if acc > 0.30 {
            potentialKeys = append(potentialKeys, length)
        }
    }

    fmt.Println("Chosen key lengths:", potentialKeys)
    return potentialKeys
}

// findKey – na podstawie długości klucza wykonuje analizę częstościową
// każdej kolumny i wybiera znak klucza, który najbardziej pasuje (max "score").
func findKey(message string, keyLength int) string {
    var key strings.Builder

    // Dla każdej pozycji w kluczu (0..keyLength-1)
    for i := 0; i < keyLength; i++ {
        // Przygotowujemy mapy do zliczania i oceniania
        positionalDict := make(map[rune]map[rune]int)
        scoredDict := make(map[rune]float64)

        for _, letter := range ALPHABET {
            letterRune := rune(letter)
            positionalDict[letterRune] = make(map[rune]int)
            scoredDict[letterRune] = 0
            // Inicjalizujemy liczniki
            for _, letter2 := range ALPHABET {
                positionalDict[letterRune][rune(letter2)] = 0
            }
        }

        // Sprawdzamy każde możliwe przesunięcie (każdą literę) jako klucz
        for _, letter := range ALPHABET {
            letterRune := rune(letter)
            idx := i
            for idx < len(message) {
                row := strings.IndexRune(ALPHABET, rune(message[idx]))
                col := strings.IndexRune(ALPHABET, letterRune)
                if row == -1 || col == -1 {
                    // Jeśli znak nie jest w alfabecie, pomijamy
                    idx += keyLength
                    continue
                }
                shifted := rune(ALPHABET[(row-col+26)%26])
                positionalDict[letterRune][shifted]++
                idx += keyLength
            }

            // Liczymy "score" – mnożymy wystąpienia przez częstość w języku
            score := 0.0
            for char, count := range positionalDict[letterRune] {
                // Zabezpieczamy się, żeby nie wziąć litery spoza frequency
                freqVal, ok := frequency[char]
                if !ok {
                    freqVal = 0.0
                }
                score += float64(count) * freqVal
            }
            scoredDict[letterRune] = score
        }

        // Znajdujemy znak z najwyższym score
        bestLetter := findMaxKey(scoredDict)
        key.WriteRune(bestLetter)
    }

    return key.String()
}

// findMaxKey – pomocnicza funkcja szukająca klucza o największej wartości (score)
// w mapie scoredDict.
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

// removeRepetitions – jeśli klucz jest kilkakrotnie powtórzony (np. "abcabc"),
// to zwracamy tylko "abc".
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

// vigenere – prosta implementacja szyfru/de-szyfru Vigenère’a.
// mode == "encrypt" -> szyfruje
// mode == "decrypt" -> deszyfruje
func vigenere(text, key, mode string) string {
    var result strings.Builder

    // Oczyszczamy i zamieniamy na małe litery, aby pokryć się z ALPHABET
    cleaned := sanitize(text)
    // Upewniamy się, że klucz jest również mały
    key = strings.ToLower(key)

    keyIndex := 0
    keyLen := len(key)

    for _, ch := range cleaned {
        row := strings.IndexRune(ALPHABET, ch)
        col := strings.IndexRune(ALPHABET, rune(key[keyIndex]))
        if row == -1 {
            // znak spoza ALPHABET, pomijamy
            continue
        }
        if col == -1 {
            // znak klucza spoza ALPHABET, pomijamy
            col = 0
        }

        var newPos int
        if mode == "decrypt" {
            // Dekodowanie
            newPos = (row - col + 26) % 26
        } else {
            // Domyślnie kodowanie
            newPos = (row + col) % 26
        }
        result.WriteByte(ALPHABET[newPos])

        keyIndex = (keyIndex + 1) % keyLen
    }

    return result.String()
}

// Function to encrypt the plainText using the Vigenère cipher with the provided key.
func EncodeVignere(plainFile, keyFile, cryptoFile string) (string, error) {
	plainText, err := GetPreparedText(plainFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać plain tekstu")
	}

	key, err := GetPreparedKey(keyFile)
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
		index := strings.IndexRune(Alphabet, char)
		keyIndex := strings.IndexRune(Alphabet, rune(key[i % len(key)]))

		encryptedIndex := (index + keyIndex) % AlphabetLen
		result = append(result, rune(Alphabet[encryptedIndex]))
	}

	// Save the decrypted text to crypto.txt
	err = SaveOutput(string(result), cryptoFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}

// decryptVigenereSimple decrypts the given cryptoFile using the Vigenère cipher with the provided key.
func DecryptVigenereSimple(cryptoFile, keyFile, decryptedFile string) (string, error) {
	cryptoText, err := GetPreparedText(cryptoFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać crypto tekstu")
	}
	fmt.Printf("Crypto Tekst: %s\n", cryptoText)
	fmt.Printf("Klucz: %s\n", keyFile)
	key, err := GetPreparedKey(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się odczytać klucza")
	}
	fmt.Printf("Klucz: %s\n", key)
	if len(cryptoText) == 0 || len(key) == 0 {
		return "", fmt.Errorf("input plainText, key, or Alphabet cannot be empty")
	}
	keyLength := len(key)
	var result []rune

	for i, char := range cryptoText {
		index := strings.IndexRune(Alphabet, char)
		if index == -1 {
			result = append(result, char) // Keep non-Alphabet characters unchanged
			continue
		}

		keyIndex := strings.IndexRune(Alphabet, rune(key[i%keyLength]))
		if keyIndex == -1 {
			return "", fmt.Errorf("invalid character in key")
		}

		decryptedIndex := (index - keyIndex + AlphabetLen) % AlphabetLen
		result = append(result, rune(Alphabet[decryptedIndex]))
	}

	fmt.Printf("Odszyfrowany tekst: %s\n", string(result))

	// Save the decrypted text to decrypt.txt
	err = SaveOutput(string(result), decryptedFile)
	if err != nil {
		log.Printf("błąd przy zapisie tekstu: %v", err)
		return "", fmt.Errorf("błąd przy zapisie tekstu: %v", err)
	}

	return string(result), nil
}



const Alphabet = "abcdefghijklmnopqrstuvwxyz"

var AlphabetLen = len(Alphabet)


// FreqMap represents the frequency of English letters in percentage (approx.)
var FreqMap = map[rune]int{
	'a': 82,  'b': 15,  'c': 28,  'd': 43,  'e': 127,
	'f': 22,  'g': 20,  'h': 61,  'i': 70,  'j': 2,
	'k': 8,   'l': 40,  'm': 24,  'n': 67,  'o': 75,
	'p': 29,  'q': 1,   'r': 60,  's': 63,  't': 91,
	'u': 28,  'v': 10,  'w': 23,  'x': 1,   'y': 20,
	'z': 1,
}

// Function to count selected flags
func CountSelectedFlags(flags []*bool) int {
	count := 0
	for _, f := range flags {
		if *f {
			count++
		}
	}
	return count
}

// Function to read text from txt file
func ReadText(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if scanner.Err() != nil{
		return nil, scanner.Err()
	}
	
	return lines, nil
}


func SaveOutput(result string, outputFile string) error {
	// Check if the file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// Create the file if it does not exist
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("błąd przy tworzeniu pliku: %v", err)
		}
		file.Close()
	}

	// Write the result to the file
	err := os.WriteFile(outputFile, []byte(result), 0644)
	if err != nil {
		return fmt.Errorf("błąd przy zapisywaniu wyniku: %v", err)
	}

	fmt.Println("Zapisano wynik do pliku:", outputFile)
	return nil
}

// Function to prepare text for encryption, cleans non-letter characters and converts to lowercase.
func CleanText(input string) (string, error) {
	var cleanedText []rune

	// Check if the input text is empty.
	if len(input) == 0 {
		return "", fmt.Errorf("input text is empty")
	}
	// Prepare the text for encryption/
	for _, char := range input {
		if unicode.IsLetter(char) {
			char = unicode.ToLower(char)
			cleanedText = append(cleanedText, char)
		}
	}

	if len(cleanedText) == 0 {
		return "", fmt.Errorf("tekst nie zawiera liter do przetworzenia")

	}
	return string(cleanedText), nil
}

func GetText(filePath string) (string, error) {
	_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			fmt.Printf("plik %s nie istnieje", filePath)
			return "", fmt.Errorf("plik %s nie istnieje", filePath)
		} else if err != nil {
			log.Printf("błąd przy sprawdzaniu istnienia pliku %s %v", filePath, err)
			return "", fmt.Errorf("błąd przy sprawdzaniu istnienia pliku %s: %v", filePath, err)
		}
	
		lines, err := ReadText(filePath)
		if err != nil {
			fmt.Printf("błąd przy odczycie pliku %s: %v", filePath, err)
			return "", fmt.Errorf("błąd przy odczycie pliku %s: %v",filePath, err)
		}
	
		if len(lines) == 0 {
			fmt.Printf("plik %s jest pusty", filePath)
			return "", fmt.Errorf("plik %s jest pusty", filePath)
		}
	
		inputText := strings.Join(lines, "\n")
		return inputText, nil
}

// Function to prepare text for encryption, cleans non-letter characters and converts to lowercase.
func PrepareText(filePath string) (string, error) {
		// Get input text.
		inputText,err := GetText(filePath)
		if err != nil {
			return "", fmt.Errorf("błąd przy odczycie pliku %s: %v", filePath, err)
		}
		preparedText, err := CleanText(inputText)
		if err != nil {
			log.Printf("błąd przy czyszczeniu tekstu: %v", err)
			return "", fmt.Errorf("błąd przy czyszczeniu tekstu: %v", err)
		}
		return preparedText, nil
}

// Function to validate the key for Vignere cipher
func Validate(text string) (error) {
	// Check if the key is empty.
	if len(text) == 0 {
		return fmt.Errorf("klucz/tekst jest pusty")
	}

	// Check if the key contains only lowercase English letters.
	for _, char := range text {
		if char < 'a' || char > 'z' { // Ensure only 'a' to 'z'
			return fmt.Errorf("klucz/tekst zawiera niedozwolony znak: %c", char)
		}
	}

	return nil
}

// Function to get the key for Vignere cipher.
func GetPreparedKey(keyFile string) (string, error) {

	key, err := PrepareText(keyFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się przygotować klucza")
	}
	// Read alpha key from file and validate key.
	err = Validate(key)
	if err != nil {
		return  "", fmt.Errorf("nie udało się zwalidować klucza")
	}

	return key, nil
}

// Function to get the key for Vignere cipher.
func GetPreparedText(textFile string) (string, error) {
	text, err := PrepareText(textFile)
	if err != nil {
		return "", fmt.Errorf("nie udało się przygotować klucza")
	}
	// Read alpha key from file and validate key.
	err = Validate(text)
	if err != nil {
		return  "", fmt.Errorf("nie udało się zwalidować klucza")
	}

	return text, nil
}

// Function Gcd finds the greatest common divisor (NWD) of two numbers
func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// MostFrequentLetter returns the most common letter in the English language
func MostFrequentLetter() rune {
	maxFreq := 0
	var mostCommon rune
	for letter, freq := range FreqMap {
		if freq > maxFreq {
			maxFreq = freq
			mostCommon = letter
		}
	}
	return mostCommon // Now return 'e'
}

// Function returns absolute value of a number.
func Absolute(x int) int {
	if x < 0 {
		return -x
	}
	return x
}