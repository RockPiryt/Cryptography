// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"math/big"

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

var PlainTextString string = "Haha"
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
		log.Println("[INFO] keys have been successfully created.")
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
		fmt.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil
	case "s":
		// Sign the message
		err := SignMsg(MessageString, MessageFile, PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to sing the message: %v", err)
		}
		log.Println("[INFO] Message successfully signed into signature.txt.")
		return nil
	case "v":
		// Verify the signed message
		_, err := VerifySignature(MessageFile, PublicKeyFile, SignatureFile)
		if err != nil {
			return fmt.Errorf("failed to verify the sign: %v", err)
		}
		log.Println("[INFO] Signature successfully verified.")
		return nil

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
	return nil
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
	fmt.Printf("Generated private key b = %s\n", b.String())
	
	// Calculate Beta (gᵇ mod p)
	beta := new(big.Int).Exp(g, b, p)

	// Checks b < p, beta < p
	if b.Cmp(p) >= 0 {
		return fmt.Errorf("invalid private key: b >= p")
	}
	if beta.Cmp(p) >= 0 {
		return fmt.Errorf("invalid public key: beta >= p")
	}

	log.Println("private b and public beta are < p")

	// Save public and private keys
	fmt.Printf("Public and private keys were saved to file")
	helpers.WriteBigIntsToFile(PrivateKeyFile, []*big.Int{p, g, b})
	helpers.WriteBigIntsToFile(PublicKeyFile, []*big.Int{p, g, beta})

	return nil
}

func EncryptElgamal(PlainFile, PublicKeyFile string) (error) {
	// Read public key with 3 values
	params, _ := helpers.ReadBigIntsFromFile(PublicKeyFile, 3)
	p, g, beta := params[0], params[1], params[2]

	// Convert sample message to Big Int
	helpers.ConvertStringToBigInt(PlainTextString, PlainFile)
	// Read message
	messages, _ := helpers.ReadBigIntsFromFile(PlainFile, 1)
	m := messages[0]
	// Check if m < pa
	if m.Cmp(p) >= 0 {
		fmt.Println("message must be less than p")
		log.Fatal("message must be less than p")
	}

	log.Println("m < p was successfully checked")

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
	// Save decrypted message
	helpers.WriteBigIntsToFile(DecryptedFile, []*big.Int{m})
	
	// Compare with original plaintext
	original, _ := helpers.ReadBigIntsFromFile(PlainFile, 1)
	originalMsg := original[0]

	if m.Cmp(originalMsg) != 0 {
		log.Println("Plaintext and decrypted text are NOT the same.")
		return fmt.Errorf("decryption mismatch: decrypted != original")
	}

	log.Println("Plaintext and decrypted text are the same.")

	return nil
}

// Signing part------------------------------------------------------------------------------------------
func SignMsg(MessageString, MessageFile, PrivateKeyFile string) error {
	// Read p,g,b from private key file
	params, _ := helpers.ReadBigIntsFromFile(PrivateKeyFile, 3)
	p, g, b := params[0], params[1], params[2]
	pm1 := new(big.Int).Sub(p, big.NewInt(1))// p-1

	// Convert sample string to Big Int
	helpers.ConvertStringToBigInt(MessageString, MessageFile)
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

func VerifySignature(MessageFile, PublicKeyFile, SignatureFile string) (string, error) {
	return "", nil
}
