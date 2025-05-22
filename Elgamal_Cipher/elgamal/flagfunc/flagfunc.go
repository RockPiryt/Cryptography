// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
)

const (
	plainFile      = "files/plain.txt"
	messageFile    = "files/message.txt"
	elgamalFile    = "files/elgamal.txt"
	privateKeyFile = "files/private.txt"
	publicKeyFile  = "files/publicKey.txt"
	cryptoFile     = "files/crypto.txt"
	decryptedFile  = "files/decrypt.txt"
	signatureFile  = "files/signature.txt"
	verifyFile     = "files/verify.txt"
)

func ExecuteCipher(operation string) error {
	switch operation {
	case "k":
		// Generate public and private key for Bolek
		err := GenerateKeys(elgamalFile)
		if err != nil {
			return fmt.Errorf("error during generation of public and private keys %v", err)
		}
		log.Println("[INFO] keys have been successfully created.")
		return nil

	case "e":
		_, err := EncryptElgamal(plainFile, publicKeyFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil

	case "d":
		// Decrypt crypto.txt using key.txt
		_, err := DecryptElgamal(cryptoFile, privateKeyFile, decryptedFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		fmt.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil
	case "s":
		// Sign the message
		_, err := SignMsg(messageFile, privateKeyFile)
		if err != nil {
			return fmt.Errorf("failed to sing the message: %v", err)
		}
		log.Println("[INFO] Message successfully signed into signature.txt.")
		return nil
	case "v":
		// Verify the signed message
		_, err := VerifySignature(messageFile, publicKeyFile, signatureFile)
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

// func GenerateKeys(elgamalFile string) error {
// 	params, _ := helpers.ReadBigIntsFromFile(elgamalFile, 2)
// 	p, g := params[0], params[1]
// 	b, _ := helpers.RandomBigInt(new(big.Int).Sub(p, big.NewInt(2)))
// 	b = b.Add(b, big.NewInt(1)) // Ensure 1 <= b < p-1

// 	beta := new(big.Int).Exp(g, b, p)

// 	helpers.WriteBigIntsToFile(privateKeyFile, []*big.Int{p, g, b})
// 	helpers.WriteBigIntsToFile(publicKeyFile, []*big.Int{p, g, beta})

// 	return nil
// }

func GenerateKeys(elgamalFile string) error {
	// params, _ := helpers.ReadBigIntsFromFile(elgamalFile, 2)
	// p, g := params[0], params[1]
	// b, _ := helpers.RandomBigInt(new(big.Int).Sub(p, big.NewInt(2)))
	// b = b.Add(b, big.NewInt(1)) // Ensure 1 <= b < p-1

	// beta := new(big.Int).Exp(g, b, p)

	// helpers.WriteBigIntsToFile(privateKeyFile, []*big.Int{p, g, b})
	// helpers.WriteBigIntsToFile(publicKeyFile, []*big.Int{p, g, beta})

	return nil
}

func EncryptElgamal(plainFile, publicKeyFile string) (string, error) {
	return "", nil
}

func DecryptElgamal(cryptoFile, privateKeyFile, decryptedFile string) (string, error) {
	return "", nil
}

func SignMsg(messageFile, privateKeyFile string) (string, error) {
	return "", nil
}

func VerifySignature(messageFile, publicKeyFile, signatureFile string) (string, error) {
	return "", nil
}
