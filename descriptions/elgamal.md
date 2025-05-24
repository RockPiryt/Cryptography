# ElGamal Public-Key Cryptosystem (Go)

This Go program implements the **ElGamal public-key cryptosystem**, including functionality for:

- Key generation
- Encryption and decryption
- Digital signing and signature verification

The system operates using large prime numbers and modular arithmetic, providing a foundational demonstration of public-key cryptography principles.

---

## Features

| Option | Description             | Input Files               | Output Files            |
|--------|-------------------------|----------------------------|--------------------------|
| `-k`   | Generate key pair        | `elgamal.txt`              | `private.txt`, `public.txt` |
| `-e`   | Encrypt message          | `plain.txt`, `public.txt` | `crypto.txt`             |
| `-d`   | Decrypt ciphertext       | `crypto.txt`, `private.txt` | `decrypt.txt`            |
| `-s`   | Sign message             | `message.txt`, `private.txt` | `signature.txt`         |
| `-v`   | Verify digital signature | `message.txt`, `public.txt`, `signature.txt` | `verify.txt` |

---

## Input/Output Formats

### `elgamal.txt`
Contains:
1. A large prime number `p`
2. A generator `g` of the multiplicative group Z\*_p

Example:
23
5


### `-k` Key Generation
- Reads `p` and `g` from `elgamal.txt`
- Randomly generates private exponent `b` such that `1 ≤ b < p−1`
- Computes `β = g^b mod p`
- Writes:
  - `private.txt`: `p`, `g`, `b`
  - `public.txt`: `p`, `g`, `β`

---

### `-e` Encryption
- Reads `p`, `g`, `β` from `public.txt`
- Reads message `m` from `plain.txt` (can be a number or ASCII-encoded string)
- Ensures `m < p`
- Randomly generates `k` such that `1 ≤ k < p−1`
- Computes:
  - `c1 = g^k mod p`
  - `c2 = m × β^k mod p`
- Saves `c1`, `c2` to `crypto.txt`

---

### `-d` Decryption
- Reads `p`, `g`, `b` from `private.txt`
- Reads `c1`, `c2` from `crypto.txt`
- Computes:
  - `s = c1^b mod p`
  - `s⁻¹ = modular inverse of s mod p`
  - `m = c2 × s⁻¹ mod p`
- If `m` appears to be ASCII text, converts it to a readable string
- Writes result to `decrypt.txt`

---

###  `-s` Signing
- Reads `p`, `g`, `b` from `private.txt`
- Reads message `m` from `message.txt`
- Picks random `k` such that `gcd(k, p−1) = 1`
- Computes:
  - `r = g^k mod p`
  - `k⁻¹ = modular inverse of k mod (p−1)`
  - `x = (m − b·r)·k⁻¹ mod (p−1)`
- Saves signature `(r, x)` to `signature.txt`

---

###  `-v` Verification
- Reads `p`, `g`, `β` from `public.txt`
- Reads message `m` from `message.txt`
- Reads `r`, `x` from `signature.txt`
- Verifies:  g^m ≡ r^x × β^r mod p

- Writes `T` (true) or `N` (false) to `verify.txt`

---

## Example Workflow

```bash
# Generate key pair
go run main.go -k

# Encrypt a message
go run main.go -e

# Decrypt the ciphertext
go run main.go -d

# Sign a message
go run main.go -s

# Verify the signature
go run main.go -v
```

## Notes
Uses math/big for large-number arithmetic.

Modular inverse is computed using the Extended Euclidean Algorithm.

When decrypting, if the message was originally a string, the tool tries to recover the original text using ASCII.

## File Descriptions
File	            Purpose
elgamal.txt	        Initial parameters p and g
private.txt	        Private key components p, g, b
public.txt	        Public key components p, g, β
plain.txt	        Message to encrypt (number or text)
crypto.txt	        Encrypted message: c1, c2
decrypt.txt	        Decrypted output
message.txt	        Message to sign
signature.txt	    Signature: r, x
verify.txt	        Verification result: T or N

## Requirements
Go 1.18+
No external libraries required
