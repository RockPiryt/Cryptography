// Author: Paulina Kimak
package flagfunc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"stegano/helpers"
	"strings"
)

const (
	MessageFile    		= "files/mess.txt"
	CoverFile      		= "files/cover.html"
	ClearedHtml			= "files/clearfile.html"
	WatermarkFile    	= "files/watermatk.html"
	DetectFile 			= "files/detect.txt"
)

func ExecuteProgram(operation string, method int) error {
	switch operation {
	case "e":
		// IN MessageFile, OUT WatermarkFile 
		err := EmbedMsg(MessageFile, CoverFile, method)
		if err != nil {
			return fmt.Errorf("failed to embeded the message: %v", err)
		}
		log.Println("[INFO] Text successfully embeded into watermark.html.")
		return nil

	case "d":
		err := ExtractMsg(WatermarkFile, DetectFile, method)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}


func EmbedMsg(MessageFile, CoverFile string, method int) error {
	messageBits, err := helpers.ReadHexBits(MessageFile)
	if err != nil {
		return err
	}

	err = helpers.ClearHtml(CoverFile)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Cleaned HTML saved as files/clearfile.html")
	}

	input, err := os.ReadFile(ClearedHtml)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	// if message is larger than lines in html
	if method == 1 && len(messageBits) > len(lines) {
		return errors.New("cover file too small for message (method 1)")
	}

	var outLines []string
	bitIndex := 0
	for _, line := range lines {
		cleanLine := strings.TrimRight(line, " ") // remove trailing spaces
		if method == 1 && bitIndex < len(messageBits) {
			// if bit in message is one then add empty space at the end
			if messageBits[bitIndex] == '1' {
				cleanLine += " "
			}
			bitIndex++
		}
		outLines = append(outLines, cleanLine)
	}

	return os.WriteFile(WatermarkFile, []byte(strings.Join(outLines, "\n")), 0644)
}


func ExtractMsg(WatermarkFile, DetectFile string, method int) error {
	input, err := os.ReadFile(WatermarkFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	var bits strings.Builder

	for _, line := range lines {
		if method == 1 {
			if strings.HasSuffix(line, " ") {
				bits.WriteByte('1')
			} else {
				bits.WriteByte('0')
			}
		}
	}

	// Convert bits back to hex string
	hexMsg := helpers.BitsToHex(bits.String())
	return os.WriteFile(DetectFile, []byte(hexMsg), 0644)
}