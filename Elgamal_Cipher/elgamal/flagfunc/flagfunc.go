// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"elgamal/helpers"
)

const (
	PlainFile      = "files/plain.txt"
	MessageFile    = "files/message.txt"
	ElgamalFile    = "files/elgamal.txt"
	PrivateKeyFile = "files/private.txt"
	PublicKeyFile  = "files/publicKey.txt"
	CryptoFile     = "files/crypto.txt"
	DecryptedFile  = "files/decrypt.txt"
	SignatureFile  = "files/signature.txt"
	VerifyFile     = "files/verify.txt"
)


var MessageString string = "Bla"

func ExecuteCipher(operation string) error {
	switch operation {
	case "k":
		// Generate public and private key for Bolek
		// IN ElgamalFile, OUT PrivateKeyFile, PublicKeyFile 
		err := GenerateKeys(ElgamalFile)
		if err != nil {
			return fmt.Errorf("error during generation of public and private keys %v", err)
		}
		log.Println("[INFO] Keys have been successfully created.")
		return nil

	case "e":
		// IN PlainFile, PublicKeyFile, OUT CryptoFile
		err := EncryptElgamal(PlainFile, PublicKeyFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil

	case "d":
		// Decrypt crypto.txt using private.txt
		// IN CryptoFile, PrivateKeyFile, OUT DecryptedFile
		err := DecryptElgamal(CryptoFile, PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil
	case "s":
		// Sign the message
		// IN MessageFile, PrivateKeyFile, OUT SignatureFile
		err := SignMsg(MessageFile, PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to sing the message: %v", err)
		}
		log.Println("[INFO] Message successfully signed into signature.txt.")
		return nil
	case "v":
		// Verify the signed message
		// IN MessageFile, PublicKeyFile, SignatureFile, OUT VerifyFile
		err := VerifySignature(MessageFile, PublicKeyFile, SignatureFile)
		if err != nil {
			return fmt.Errorf("failed to verify the sign: %v", err)
		}
		log.Println("[INFO] Signature successfully verified.")
		return nil

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}

// Encrypting part------------------------------------------------------------------------------------------
// Function generate 2 keys: beta is public key, b is private key and saves these values to file
func GenerateKeys(ElgamalFile string) error {
	// Read p and g from file
	params, _ := helpers.ReadBigIntsFromFile(ElgamalFile, 2)
	p, g := params[0], params[1]

	// Generate b (random), where  1 <= b < p-1
	upperLimit := new(big.Int).Sub(p, big.NewInt(2)) // p - 2
	b, err := helpers.RandomBigInt(upperLimit)       // 0 <= b < p-2
	if err != nil {
		return fmt.Errorf("error generating random b: %v", err)
	}
	b.Add(b, big.NewInt(1)) // 1 <= b < p-1
	log.Printf("[INFO] Generated private key b = %s\n", b.String())
	
	// Calculate Beta (gᵇ mod p)
	beta := new(big.Int).Exp(g, b, p)

	// Checks b < p, beta < p
	if b.Cmp(p) >= 0 {
		return fmt.Errorf("invalid private key: b >= p")
	}
	if beta.Cmp(p) >= 0 {
		return fmt.Errorf("invalid public key: beta >= p")
	}

	log.Println("[INFO] Private b and public beta are < p")

	// Save public and private keys
	log.Printf("[INFO] Public and private keys were saved to file")
	helpers.WriteBigIntsToFile(PrivateKeyFile, []*big.Int{p, g, b})
	helpers.WriteBigIntsToFile(PublicKeyFile, []*big.Int{p, g, beta})

	return nil
}

func EncryptElgamal(PlainFile, PublicKeyFile string) (error) {
	// Read public key with 3 values
	params, _ := helpers.ReadBigIntsFromFile(PublicKeyFile, 3)
	p, g, beta := params[0], params[1], params[2]

	// Read raw content from plain.txt
	content, err := os.ReadFile(PlainFile)
	if err != nil {
		return err
	}
	raw := strings.TrimSpace(string(content))

	m := new(big.Int)

	// Try to parse as number
	_, ok := m.SetString(raw, 10)
	if !ok {
		// If parsing as number failed, treat as string
		m.SetBytes([]byte(raw))
	}
	
	// Check if m < pa
	if m.Cmp(p) >= 0 {
		fmt.Println("message must be less than p")
		log.Fatal("[FATAL] message must be less than p")
	}

	log.Println("[INFO] m < p was successfully checked")

	//Get random number k, where 1 ≤ k < p − 1
	upperLimit := new(big.Int).Sub(p, big.NewInt(2))// Calculate upper limit: p - 2
	k, _ := helpers.RandomBigInt(upperLimit)// Generate random number: 0 <= k < p-2
	k.Add(k, big.NewInt(1))	// Add 1 to make sure 1 <= k < p-1


	// Calucalte c1 = gᵏ mod p
	c1 := new(big.Int).Exp(g, k, p)
	//Calcualate a masking value s = βᵏ mod p
	s := new(big.Int).Exp(beta, k, p)
	//Calculate c2 = m × (βᵏ mod p), 
	c2 := new(big.Int).Mul(m, s)
	c2.Mod(c2, p)

	//Save cryptogram as 2 values to file
	helpers.WriteBigIntsToFile(CryptoFile, []*big.Int{c1, c2})
	return  nil
}

func DecryptElgamal(CryptoFile, PrivateKeyFile string) (error) {
	//Read params from private key file
	params, _ := helpers.ReadBigIntsFromFile(PrivateKeyFile, 3)
	p, _, b := params[0], params[1], params[2]

	//Read cryptogram (c1,c2)
	cipher, _ := helpers.ReadBigIntsFromFile(CryptoFile, 2)
	c1, c2 := cipher[0], cipher[1]

	// Calculate s = c1^b mod p
	s := new(big.Int).Exp(c1, b, p)
	// Calculate s^(-1)
	sInv, err := helpers.ModInverse(s, p)
	if err != nil {
		log.Fatal(err)
	}
	// Get oryginal message (calculate m = c2 · s^(-1) mod p)
	m := new(big.Int).Mul(c2, sInv)
	m.Mod(m, p)

	// Get oryginal string 
	// Attempt to convert decrypted big.Int to ASCII string
	byteMsg := m.Bytes()
	decodedStr := string(byteMsg)

	// Check if all bytes are printable ASCII (32–126)
	isPrintable := true
	for _, b := range byteMsg {
		if b < 32 || b > 126 {
			isPrintable = false
			break
		}
	}

	var output string
	if isPrintable && len(decodedStr) > 0 {
		output = decodedStr // treat as ASCII string
		log.Println("[INFO] Decrypted message interpreted as string.")
	} else {
		output = m.String() // fallback to number
		log.Println("[INFO] Decrypted message interpreted as number.")
	}


	// Save decrypted message
	err = os.WriteFile(DecryptedFile, []byte(output+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("failed to write decrypted message: %v", err)
	}
	
	// -------- Compare to plain.txt --------
	plainRaw, err := os.ReadFile(PlainFile)
	if err != nil {
		return fmt.Errorf("failed to read original plaintext: %v", err)
	}
	plainContent := strings.TrimSpace(string(plainRaw))

	// Try match as number
	plainNum := new(big.Int)
	if _, ok := plainNum.SetString(plainContent, 10); ok {
		// Was number: compare big.Int values
		if m.Cmp(plainNum) != 0 {
			log.Fatal("[FATAL] Decrypted number does not match original number.")
			return fmt.Errorf("decryption mismatch: numeric value differs")
		}
		log.Println("[INFO] Decrypted number matches original.")
	} else {
		// Was string: compare decodedStr with plain text
		if output != plainContent {
			log.Fatal("[FATAL] Decrypted string does not match original string.")
			return fmt.Errorf("decryption mismatch: string differs")
		}
		log.Println("[INFO] Decrypted string matches original.")
	}

	log.Println("[INFO] Plaintext and decrypted text are the same.")

	return nil
}

// Signing part------------------------------------------------------------------------------------------
func SignMsg(MessageFile, PrivateKeyFile string) error {
	// Read p,g,b from private key file
	params, _ := helpers.ReadBigIntsFromFile(PrivateKeyFile, 3)
	p, g, b := params[0], params[1], params[2]
	pm1 := new(big.Int).Sub(p, big.NewInt(1))// p-1

	// Convert sample string to Big Int
	helpers.CreateShortcutSHA(MessageString, MessageFile)
	//Read message to signing
	messages, _ := helpers.ReadBigIntsFromFile(MessageFile, 1)
	m := messages[0]

	//k is random, where 1 ≤ k < p-1, and gcd(k, p-1) = 1 (IsCoprime)
	var k, kInv *big.Int
	for {
		kCandidate, err := helpers.RandomBigInt(pm1) // 0 <= k < p-1
		if err != nil {
			return fmt.Errorf("failed to generate random k: %v", err)
		}
		kCandidate.Add(kCandidate, big.NewInt(1)) // 1 <= k < p
		if helpers.IsCoprime(kCandidate, pm1) {
			k = kCandidate
			break
		}
	}
	//k_inv = inversed k mod (p−1)
	kInv, err := helpers.ModInverse(k, pm1)
	if err != nil || kInv == nil {
		return fmt.Errorf("failed to compute modular inverse of k")
	}

	// Calculate signature (r,x)
	// caculate r = g^k mod p
	r := new(big.Int).Exp(g, k, p)
	//x = (m - b·r)·k_inv mod (p−1)
	x := new(big.Int).Mul(b, r)
	x.Sub(m, x)
	x.Mul(x, kInv)
	x.Mod(x, pm1)

	// Save signature
	err = helpers.WriteBigIntsToFile(SignatureFile, []*big.Int{r, x})
	if err != nil {
		return fmt.Errorf("failed to save signature: %v", err)
	}

	log.Printf("[INFO] Signed message and saved signature (r, x) to %s", SignatureFile)
	return nil
}

func VerifySignature(MessageFile, PublicKeyFile, SignatureFile string) error {
	//Read p,g,beta from public key
	params, _ := helpers.ReadBigIntsFromFile(PublicKeyFile, 3)
	p, g, beta := params[0], params[1], params[2]

	// Read oryginal message
	messages, _ := helpers.ReadBigIntsFromFile(MessageFile, 1)
	m := messages[0]
	// Read signature 
	sig, _ := helpers.ReadBigIntsFromFile(SignatureFile, 2)
	r, x := sig[0], sig[1]
	
	//g^m  ≡ r^x · β^r mod p
	left := new(big.Int).Exp(g, m, p)
	right1 := new(big.Int).Exp(r, x, p)
	right2 := new(big.Int).Exp(beta, r, p)
	right := new(big.Int).Mul(right1, right2)
	right.Mod(right, p)

	result := "N"
	if left.Cmp(right) == 0 {
		result = "T"
	}

	// Save result
	os.WriteFile(VerifyFile, []byte(result), 0644)
	fmt.Println(result)
	return nil
}
