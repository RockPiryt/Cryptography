# Stegano: HTML-Based Message Hiding Tool (Go)

This Go program implements a **steganography system for HTML files**, enabling secret messages (encoded in hexadecimal) to be embedded in or extracted from a web page using multiple encoding techniques.

---

## Features

| Option | Description                 | Input Files                   | Output Files          |
|--------|-----------------------------|--------------------------------|------------------------|
| `-e`   | Embed message                | `mess.txt`, `cover.html`       | `watermark.html`       |
| `-d`   | Extract hidden message       | `watermark.html`               | `detect.txt`           |
| `-1`   | Use Method 1 for encoding    | (used with `-e` or `-d`)       |                        |
| `-2`   | Use Method 2 for encoding    |                                |                        |
| `-3`   | Use Method 3 for encoding    |                                |                        |
| `-4`   | Use Method 4 for encoding    |                                |                        |

> 💡 When using `-e`, one of the methods `-1` to `-4` **must be specified**.  
> 💡 When using `-d`, the same method used during embedding must be specified.

---

## Message Format

### `mess.txt`

The message is encoded in **hexadecimal**, similar to hash outputs. Each character from `0–9, a–f` represents 4 bits of binary data.

Example: 48656c6c6f20576f726c64 → represents `"Hello World"` in binary.

---

## Embedding Methods

| Method | Description |
|--------|-------------|
| `-1` | Each message bit is encoded as an **extra space at the end of a line**. The message length is limited to the number of lines in the HTML file. |
| `-2` | Each bit is represented as a **single or double space** between words (tab characters are ignored). The message length is limited to the number of unique space regions. |
| `-3` | Each bit is encoded via **typos in CSS attribute names** in `<p style="...">` tags. For example: `"margin-botom"` (for bit 0) or `"lineheight"` (for bit 1). |
| `-4` | Each bit is encoded using **redundant opening and closing tags**, e.g., extra `<FONT>` tags: `bit 1` causes an extra open-close-open pattern, while `bit 0` adds an extra empty close-open at the end. The message length is limited to the number of `<FONT>` tags. |

---

## Extraction Logic

For extraction using `-d`, the same method used to embed the message must be specified using the same flag `-1` to `-4`.  
The program reads the modified HTML file `watermark.html` and reconstructs the original hex message into `detect.txt`.

---

## File Descriptions

| File             | Purpose                                      |
|------------------|----------------------------------------------|
| `mess.txt`       | Contains the message to hide (in hex)        |
| `cover.html`     | Original HTML content (carrier)              |
| `watermark.html` | HTML with the hidden message embedded        |
| `detect.txt`     | Output file containing the extracted message |

---

## Cleanup Requirement

Before embedding any message, the program will automatically **sanitize** the carrier file by:
- Removing trailing spaces
- Removing duplicate or irregular spaces
- Removing existing invalid attribute names
- Removing redundant HTML tag sequences (used in methods 3 and 4)

This ensures **a clean slate**, preventing accidental decoding of leftover or malformed data.

---

## Example Usage

```bash
# Embed a message using method 1
go run stegano.go -e -1 cover.html mess.txt

# Extract the message using the same method
go run stegano.go -d -1 watermark.html detect.txt
