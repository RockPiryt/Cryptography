

# Caesar and Affine  Cipher Encryption / Decryption Program
This Go program implements Caesar cipher encryption and decryption. It allows users to encrypt and decrypt messages using these classical ciphers directly from the command line.

## Features
- Caesar Cipher encryption and decryption
- Affine Cipher encryption and decryption
- Cryptanalysis using either plaintext or ciphertext
- Key reading from a file, and results writing to a file

## Requirements

To run the program, you need:

- Go 1.x or later installed on your machine
- The necessary files (`plain.txt`, `key.txt`) in the required format for each operation

## Usage

The program accepts the following flags:

### Cipher Options:
The program should encrypt and decrypt messages using the Caesar cipher and the affine cipher.
The program, named cezar, should be executable from the command line with the following options:

- `-c`: Caesar cipher
- `-a`: Affine cipher

### Operation Options:
- `-e`: Encryption
- `-d`: Decryption
- `-j`: Cryptanalysis with plaintext (requires both `plain.txt` and `crypto.txt`)
- `-k`: Cryptanalysis with ciphertext (requires `crypto.txt`)

### Example Commands

#### Encrypting with Caesar Cipher:

```bash
go run main.go -c -e
```

### Files
The program should read data from certain files and write to others, with fixed file names as specified:
- plain.txt: file containing the plaintext (one line, letters and spaces),
- crypto.txt: file containing the encrypted text,
- decrypt.txt: file containing the decrypted text,
- key.txt: file containing the key (one line, with the first number representing the shift and the second number for the affine cipher coefficient, the numbers are separated by a space),
- extra.txt: file containing the beginning of the plaintext for cryptanalysis with both plaintext and ciphertext,
- key-found.txt: file containing the found key in case of cryptanalysis with both plaintext and ciphertext.

### Features
- The encryption program reads the plaintext and key, then writes the encrypted text. If the key is invalid, it raises an error.

- The decryption program reads the encrypted text and key, then writes the plaintext. If the key is invalid, it raises an error. For the affine cipher, the task is to find the inverse of the number a given as part of the key – do not assume that the decryption program will receive this inverse.

- The cryptanalysis program with plaintext reads the encrypted text and the helper text, then writes the found key and the decrypted text. If it's impossible to find the key, an error should be raised.

- The cryptanalysis program without plaintext reads only the encrypted text and writes all possible candidates (25 for the Caesar cipher, 311 for the affine cipher) as decrypted text.

- The program should never require the existence of unnecessary files for a given option. Files that need to be written to should be created if they do not exist.
