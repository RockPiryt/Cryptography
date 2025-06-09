# Rabin-Miller Primality Test & Factorization Tool (Go)

This Go program implements the **Rabin-Miller probabilistic primality test**, with optional support for discovering a **nontrivial factor** of a composite number under certain conditions.

It is designed to analyze large integers using **witness-based testing** and optional universal exponents. The tool also includes a simpler **Fermat primality test** mode for comparison.

---

## Features

| Option   | Description                                | Input File      | Output File     |
|----------|--------------------------------------------|------------------|------------------|
| *(no flag)* | Full Rabin-Miller test (with factor detection) | `wejscie.txt`    | `wyjscie.txt`    |
| `-f`     | Fermat test only (no factor detection)     | `wejscie.txt`    | `wyjscie.txt`    |

---

## Input Format – `wejscie.txt`

This file must contain **1 to 3 lines**, depending on the use case:

| Line | Content                              | Purpose                                   |
|------|--------------------------------------|-------------------------------------------|
| 1    | An integer `n`                       | Number to test for primality              |
| 2    | *(optional)* A universal exponent `r`| Speeds up factorization (Rabin-Miller)    |
| 3    | *(optional)* Another factor for computing `r = (line2 * line3) - 1` | Used instead of line 2 |

> 💡 If `-f` is used, only the first line is read.

---

## Output Format – `wyjscie.txt`

The program writes one of the following results:

- `"prawdopodobnie pierwsza"` – likely prime (with error probability < 2⁻⁴⁰)
- `"na pewno złożona"` – definitely composite
- `a` – a nontrivial factor of `n` (if one was found)

---

## Algorithm Overview

### 🔍 Rabin-Miller Test

The Rabin-Miller algorithm is based on:

- **Fermat's Little Theorem**:  
  If `n` is prime, then `a^(n−1) ≡ 1 (mod n)` for all `a` such that `gcd(a, n) = 1`

- For a given odd number `n`, the algorithm:
  1. Writes `n−1 = m·2^k`, where `m` is odd.
  2. Picks random bases `a` in `[2, n−2]`.
  3. Computes `b₀ = a^m mod n`, then squares repeatedly: `bⱼ₊₁ = bⱼ^2 mod n`.
  4. If the result chain never hits `−1` (mod `n`) before reaching `1`, then:
     - A nontrivial square root of 1 is found → extract a factor of `n`.

### 🧪 Fermat Test Mode (`-f`)

In this simplified mode:

- Only `a^(n−1) mod n` is computed for random `a`.
- If the result is not 1, then `n` is composite.
- This mode cannot discover factors or Carmichael numbers.

---

## Example Input

718548065973745507
3449
543546506135745129


Expected output:
- Either `740876531` or `969862097`, both are nontrivial factors of `n`.

Why?  
Because `718548065973745507 = 740876531 × 969862097`  
and a square root of 1 mod `n` (i.e., a `bⱼ` such that `bⱼ^2 ≡ 1 mod n`, but `bⱼ ≠ ±1`) reveals a factor.

---

## Carmichael Example

- Input: `561`  
- Output:
  - With `-f`: `"prawdopodobnie pierwsza"` ❌ (false positive)
  - Without `-f`: `"33"` or `"17"` ✅ (reveals a factor)

---

## Limitations & Libraries

- For cryptographic-sized numbers (hundreds of bits), **big integer libraries** are required.
- In Go, the built-in `math/big` package handles this.
- 🚫 Do not use built-in `ProbablyPrime()` or `IsPrime()` functions – implement Rabin-Miller manually.

---

## Example Usage

```bash
# Run full Rabin-Miller test (with possible factorization)
go run rabinmiller.go

# Run Fermat-only primality test
go run rabinmiller.go -f

Requirements
Go 1.18+

Uses math/big for large-number support

No external libraries required