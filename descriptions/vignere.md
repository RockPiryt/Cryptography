# Project Name
> Vigenère Cipher Implementation

## Table of Contents
* [General Info](#general-information)
* [Technologies Used](#technologies-used)
* [Screenshots](#screenshots)
* [Setup](#setup)
* [Usage](#usage)
* [Project Status](#project-status)
* [Acknowledgements](#acknowledgements)
* [Contact](#contact)
* [License](#license)
* [Cipher Description](#cipher-description)

## General Information
This project implements encryption and decryption of messages using the **Vigenère cipher**. The plaintext consists of a sequence of lowercase letters without spaces, digits, or punctuation marks. The plaintext is prepared from real text using an appropriate preprocessing tool.

The program, named **vigenere**, supports command-line execution with the following options:

- `-p` (prepare plaintext for encryption)
- `-e` (encryption)
- `-d` (decryption)
- `-k` (cryptanalysis based solely on the ciphertext)

The primary focus of this project is **cryptanalysis**. While short ciphertexts may not yield satisfactory results, longer texts (e.g., encrypted newspaper articles) can often be successfully decrypted without knowing the key.

## Technologies Used
- Go 1.23.5

## Screenshots
![Example screenshot](./img/vigenere.png)

## Setup
1. Clone the repository or download the source code.
    ```bash
    git clone <repository_url>
    ```
2. Navigate to the project directory.
    ```bash
    cd <project_directory>
    ```
3. Run the program using the Go command with appropriate flags. Example:
    ```bash
    go run main.go -e 
    ```
    or 
    ```bash
    ./vigenere.exe -e    
    ```

## Project Status
Project is: _in progress_ 

## Acknowledgements
- Cryptography

## Contact
Created by [@rockpiryt](https://www.paulinakimakcom/) - feel free to contact me!

## License
This project is open source and available under the [MIT License]().

## Cipher Description
For a detailed explanation of the Vigenère cipher, refer to:
- [Vigenère Cipher](./descriptions/vigenere.md)

To check
Get-Content ".\org.txt" | measure -Character


# Vigenère Cipher Encryption / Decryption Program

This Go program implements encryption, decryption, and cryptanalysis of the Vigenère cipher. It allows processing of plaintext without special characters or spaces, and includes functionality for recovering the key using ciphertext alone.

## Features

- Vigenère cipher encryption and decryption  
- Plaintext preparation (removes digits, punctuation, and spaces)  
- Cryptanalysis without known plaintext  
- Automatic key length detection using coincidence analysis  
- Caesar-shift estimation for each key position  
- Command-line interface for all operations  


### Operation Options

- `-p`: Preprocess the original text (`orig.txt`) by removing non-letter characters and writing the result to `plain.txt`  
- `-e`: Encrypt the plaintext using the key from `key.txt`, writing the result to `crypto.txt`  
- `-d`: Decrypt the ciphertext using the key from `key.txt`, writing the result to `decrypt.txt`  
- `-k`: Perform cryptanalysis based solely on `crypto.txt`, attempting to recover the key and plaintext  

## Files

The following fixed filenames are used:

- `orig.txt`: The original text (raw input)  
- `plain.txt`: The processed plaintext (only lowercase letters)  
- `crypto.txt`: The encrypted text  
- `decrypt.txt`: The decrypted result  
- `key.txt`: The encryption key (lowercase letters)  
- `key-found.txt`: The key recovered during cryptanalysis (if successful)  

## Cryptanalysis Method

The core of the project is the cryptanalysis of the Vigenère cipher **without access to the plaintext**. The method involves two main steps:

### Step 1: Key Length Detection

For each shift value \( j = 1, 2, 3, \dots \), the ciphertext is shifted by \( j \) positions and compared to the original. The number of letter matches at the same positions is calculated. Peaks in the number of matches suggest a multiple of the key length. The smallest significant peak gives an estimate of the key length \( n \).  
To ensure meaningful analysis, shifts should start from at least 4.

### Step 2: Frequency Analysis per Key Position

For each position \( i = 0, 1, \dots, n-1 \), extract every \( n \)-th letter from the ciphertext (i.e., letters at positions where `index % n == i`).  

For each such slice:

1. Compute a frequency vector \( V \) of letter occurrences.  
2. Compare \( V \) to the standard frequency distribution \( W \) (based on natural language statistics).  
3. Calculate scalar products \( V \cdot W_j \), where \( W_j \) is the vector \( W \) shifted by \( j \).  
4. The value of \( j \) that gives the highest scalar product indicates the Caesar shift for that key position.  

Combining all estimated shifts reconstructs the full Vigenère key.

> **Note**: Cryptanalysis works best on ciphertexts of several hundred characters or more. Short messages may not yield reliable results.
