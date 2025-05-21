// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"elgamal/helpers"
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
		_, err := EncodeElgamal(plainFile, publicKeyFile)
		if err != nil {
			return fmt.Errorf("failed to encrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully encrypted into crypto.txt.")
		return nil
		
	case "d":
		// Decrypt crypto.txt using key.txt
		_, err = DecryptElgamal(cryptoFile, privateKeyFile, decryptedFile)
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
