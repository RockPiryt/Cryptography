# XOR Cipher Cryptanalysis Tool

This Go program demonstrates the weaknesses of the XOR (one-time pad) cipher when the same key is reused. It allows encryption, decryption (partial), and basic cryptanalysis of multiple ciphertexts encrypted with the same key using XOR.

## Principle of the Cipher

The one-time pad (OTP) / XOR cipher relies on the **non-reusability of the key**. Encryption and decryption are performed with the XOR operation:

- **Encryption**: `E(k, m) = k ⊕ m`
- **Decryption**: `D(k, c) = k ⊕ c`

If two messages `m1` and `m2` are encrypted using the same key `k`, then:

c1 = m1 ⊕ k
c2 = m2 ⊕ k
=> m1 ⊕ m2 = c1 ⊕ c2

This means that the XOR of the ciphertexts reveals the XOR of the original messages. Even if the original messages are not fully recovered, any information gained constitutes a **break** of the cipher's security.

## Assumptions

- Only **letters and spaces** are encrypted (for simplicity, possibly only lowercase letters).
- Text is encoded in **ASCII**:
  - Space = 32 (`0x20`)
  - Lowercase letters = 97–122 (`0x61` to `0x7A`)
- In binary, space starts with `001`, letters with `011`.
- XOR of two letters starts with `000...`
- XOR of a letter and a space starts with `010...`
- Therefore, if `m1 ⊕ m2` starts with `010`, we know one character is a space, and `m1 ⊕ m2 ⊕ 0x20` gives the other character — though we may not know which is which.

This property can be used to **analyze ciphertexts**, deduce likely positions of spaces, and infer characters in plaintext.

## Program Overview

The program is named `xor` and supports the following command-line options:

### Command-Line Options

- `-p`: **Prepare** text for demonstration (generate `plain.txt` from `orig.txt`)
- `-e`: **Encrypt** the prepared plaintext using a given key (`key.txt`)
- `-k`: **Cryptanalysis** based only on the ciphertexts (`crypto.txt`), without knowing the key

## File Descriptions

| File        | Description |
|-------------|-------------|
| `orig.txt`  | Any source text for demonstration purposes |
| `plain.txt` | Prepared plaintext with multiple lines of equal length (e.g., 64 characters) |
| `key.txt`   | The encryption key (string of characters with matching line length) |
| `crypto.txt`| The encrypted ciphertext (each line = `plain ⊕ key`) |
| `decrypt.txt` | The decrypted output after cryptanalysis. If a character cannot be recovered, it is replaced with `_`. |

## Cryptanalysis Method

The cryptanalysis (`-k` option) works as follows:

1. **XOR of ciphertexts**:
   - For each pair of ciphertexts `c1` and `c2`, compute `c1 ⊕ c2 = m1 ⊕ m2`.
   - Analyze the result to identify positions likely to contain a space.

2. **Character Inference**:
   - If `m1 ⊕ m2` at some position starts with `010`, then one of them is a space.
   - By XORing with `0x20` (the space character), we get the candidate for the other plaintext character.

3. **Multiple Comparisons**:
   - By comparing many ciphertexts together, patterns and likely characters can be inferred.
   - If two ciphertexts produce the same character at a position (i.e., XOR result is `0x00`), they may be the same letter or both spaces.

4. **Key Recovery**:
   - Once a space is identified in any position of a line, the key at that position can be recovered:
     - `key[i] = ciphertext[i] ⊕ 0x20`
   - With this partial key, other ciphertexts can be decrypted at that position.

5. **Decryption Output**:
   - If a character in plaintext cannot be determined with confidence, the program writes `_` in its place.

> The effectiveness of the cryptanalysis increases with the number of ciphertext lines encrypted with the same key.

## Example Workflow

```bash
# Prepare plain.txt from orig.txt
go run .\xor.go -p

# Encrypt plain.txt using key.txt into crypto.txt
go run .\xor.go -e

# Perform cryptanalysis based only on crypto.txt
go run .\xor.go -k

