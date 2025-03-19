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

## Usage
The program uses the following files for input and output:

- **plain.txt**: Contains the plaintext
- **crypto.txt**: Contains the encrypted text
- **decrypt.txt**: Contains the decrypted text
- **key.txt**: Contains the encryption key (a short sequence of letters indicating the shifts for each position of the key)
- **orig.txt**: Original text before preparation for encryption
- **key-found.txt**: Contains the key found as a result of cryptanalysis

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

