# Block Cipher Visualization (ECB vs CBC)

## Task Description

This task demonstrates the visual difference between two block cipher modes:
- **ECB (Electronic Codebook)**
- **CBC (Cipher Block Chaining)**

In this exercise, you will visually understand why ECB is considered insecure in practice, especially for repeated encryption or known plaintext attacks.

---

## Goal

Create a program that performs pseudo-encryption of a **black and white image** using ECB and CBC modes. The image should be divided into small blocks (e.g., 8×8 pixels). Each block will be independently processed using a pseudo-encryption function (e.g., a hash function like SHA256).

The output will be **two encrypted images**:
- `ecb_crypto.bmp`
- `cbc_crypto.bmp`

This will demonstrate how ECB leaks structure while CBC masks it.

## Input and Output

### Input files:
- `plain.bmp` – grayscale or 1-bit black & white BMP image
- `key.txt` *(optional)* – a text file with a custom key

### Output files:
- `ecb_crypto.bmp` – image encrypted in ECB mode
- `cbc_crypto.bmp` – image encrypted in CBC mode

---

## Block Processing Logic

- Divide the image into **blocks** (e.g., 8×8 pixels)
- Process each block:
  - In **ECB**, encrypt each block individually.
  - In **CBC**, encrypt the block XOR'd with the previous encrypted block.
- For encryption, apply a deterministic transformation (e.g., SHA-256 of the block data + key).

---

## How to Run

1. Prepare input image as `plain.bmp`
2. (Optionally) create a `key.txt` with your chosen key
3. Compile and run the Go program:

```bash
go build -o block block.go
./block
```
or 
```
go run  block.go
```

4. Program outputs two BMP images:
   - `ecb_crypto.bmp`
   - `cbc_crypto.bmp`

---

## Notes

- Keep the image **simple**, like a large font letter or logo.
- ECB will likely preserve visual structure (bad!)
- CBC should destroy visual patterns (good!)

---

## Example Use Case

- `plain.bmp` contains a large black letter “A” on white background
- `ecb_crypto.bmp` still shows the shape of “A” (pattern leakage)
- `cbc_crypto.bmp` is visually unrecognizable (structure hidden)


