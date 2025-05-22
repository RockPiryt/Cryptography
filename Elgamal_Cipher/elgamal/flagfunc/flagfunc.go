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

func ExecuteCipher(operation string) error {
	switch operation {
	case "k":
		// Generate public and private key for Bolek
		err := GenerateKeys(ElgamalFile)
		if err != nil {
			return fmt.Errorf("error during generation of public and private keys %v", err)
		}
		log.Println("[INFO] keys have been successfully created.")
		return nil

	case "e":
		err := EncryptElgamal(PlainFile, PublicKeyFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil

	case "d":
		// Decrypt crypto.txt using key.txt
		_, err := DecryptElgamal(CryptoFile, PrivateKeyFile, DecryptedFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		fmt.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil
	case "s":
		// Sign the message
		_, err := SignMsg(MessageFile, PrivateKeyFile)
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

func GenerateKeys(ElgamalFile string) error {
	// Read p and g from file
	params, _ := helpers.ReadBigIntsFromFile(ElgamalFile, 2)
	p, g := params[0], params[1]
	b, _ := helpers.RandomBigInt(new(big.Int).Sub(p, big.NewInt(2)))
	b = b.Add(b, big.NewInt(1)) // Ensure 1 <= b < p-1
	
	// Calculate Beta (gᵇ mod p)
	beta := new(big.Int).Exp(g, b, p)

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

	// Read message
	messages, _ := helpers.ReadBigIntsFromFile(PlainFile, 1)
	m := messages[0]

	// Check if m < pa
	if m.Cmp(p) >= 0 {
		fmt.Println("message must be less than p")
		log.Fatal("message must be less than p")
	}

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

func DecryptElgamal(CryptoFile, PrivateKeyFile, DecryptedFile string) (string, error) {
	return "", nil
}

func SignMsg(MessageFile, PrivateKeyFile string) (string, error) {
	return "", nil
}

func VerifySignature(MessageFile, PublicKeyFile, SignatureFile string) (string, error) {
	return "", nil
}
