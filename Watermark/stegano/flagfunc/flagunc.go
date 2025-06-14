// Author: Paulina Kimak
package flagfunc

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
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
		log.Println("Cleaned HTML saved as files/clearfile.html")
	}

	input, err := os.ReadFile(ClearedHtml)
	if err != nil {
		return err
	}

	switch method {
	case 1:
		return embedMethod1(input, messageBits)
	case 2:
		return embedMethod2(input, messageBits)
	case 3:
		return embedMethod3(input, messageBits)
	case 4:
		return embedMethod4(input, messageBits)
	default:
		return errors.New("invalid method")
	}
}

func ExtractMsg(WatermarkFile, DetectFile string, method int) error {
	input, err := os.ReadFile(WatermarkFile)
	if err != nil {
		return err
	}

	switch method {
	case 1:
		return extractMethod1(input, DetectFile)
	case 2:
		return extractMethod2(input, DetectFile)
	case 3:
		return extractMethod3(input, DetectFile)
	case 4:
		return extractMethod4(input, DetectFile)
	default:
		return errors.New("invalid method")
	}
}

// Method 1: Embed by adding a space at the end of each line based on bit value
// Method 1: Embed by adding a space at the end of each line based on bit value
func embedMethod1(input []byte, messageBits string) error {
	lines := strings.Split(string(input), "\n")
	if len(messageBits) > len(lines) {
		return errors.New("cover file too small for message (method 1)")
	}
	var outLines []string
	bitIndex := 0
	for _, line := range lines {
		cleanLine := strings.TrimRight(line, " ")
		if bitIndex < len(messageBits) {
			if messageBits[bitIndex] == '1' {
				cleanLine += " "
			}
			bitIndex++
		}
		outLines = append(outLines, cleanLine)
	}
	return os.WriteFile(WatermarkFile, []byte(strings.Join(outLines, "\n")), 0644)
}

// Method 2: Embed bits using single vs. double space regions (e.g., " " for 0, "  " for 1)
func embedMethod2(input []byte, messageBits string) error {
	text := string(input)
	spaceRegions := regexp.MustCompile(`[^\S
]+`).FindAllStringIndex(text, -1)
	if len(messageBits) > len(spaceRegions) {
		return errors.New("cover file too small for message (method 2)")
	}
	var sb strings.Builder
	last := 0
	for i, region := range spaceRegions {
		sb.WriteString(text[last:region[0]])
		if i < len(messageBits) {
			if messageBits[i] == '1' {
				sb.WriteString("  ")
			} else {
				sb.WriteString(" ")
			}
		} else {
			sb.WriteString(text[region[0]:region[1]])
		}
		last = region[1]
	}
	sb.WriteString(text[last:])
	return os.WriteFile(WatermarkFile, []byte(sb.String()), 0644)
}

// Method 3: Embed bits by introducing typos in style attributes of paragraph tags
// 0 → "margin-botom", 1 → "lineheight"
func embedMethod3(input []byte, messageBits string) error {
	text := string(input)
	styleRegex := regexp.MustCompile(`<p[^>]*style="[^"]*"[^>]*>`)
	styleMatches := styleRegex.FindAllStringIndex(text, -1)
	if len(messageBits) > len(styleMatches) {
		return errors.New("cover file too small for message (method 3)")
	}
	var sb strings.Builder
	last := 0
	for i, match := range styleMatches {
		sb.WriteString(text[last:match[0]])
		original := text[match[0]:match[1]]
		modified := original
		if messageBits[i] == '1' {
			modified = strings.Replace(modified, "line-height", "lineheight", 1)
		} else {
			modified = strings.Replace(modified, "margin-bottom", "margin-botom", 1)
		}
		sb.WriteString(modified)
		last = match[1]
	}
	sb.WriteString(text[last:])
	return os.WriteFile(WatermarkFile, []byte(sb.String()), 0644)
}

// Method 4: Embed bits using extra <font> tag patterns
// 1 → duplicate open-close-open pattern, 0 → duplicate closing tags
func embedMethod4(input []byte, messageBits string) error {
	text := string(input)
	fontRegex := regexp.MustCompile(`(?i)<font[^>]*>`)
	matches := fontRegex.FindAllStringIndex(text, -1)
	if len(messageBits) > len(matches) {
		return errors.New("cover file too small for message (method 4)")
	}
	var sb strings.Builder
	last := 0
	for i, match := range matches {
		sb.WriteString(text[last:match[1]])
		if messageBits[i] == '1' {
			sb.WriteString("</font><font>")
		}
		last = match[1]
	}
	sb.WriteString(text[last:])
	return os.WriteFile(WatermarkFile, []byte(sb.String()), 0644)
}


// Extraction for Method 1: Check if line ends with space
func extractMethod1(input []byte, DetectFile string) error {
	text := string(input)
	lines := strings.Split(text, "\n")
	var bits strings.Builder
	bitCount := 0
	for _, line := range lines {
		if strings.HasSuffix(line, " ") {
			bits.WriteByte('1')
			bitCount++
		} else if strings.TrimSpace(line) != "" {
			bits.WriteByte('0')
			bitCount++
		}
		if bitCount%4 == 0 && helpers.IsHex(bits.String()) {
			if helpers.BitsToHex(bits.String()) == helpers.ReadFileContent("files/mess.txt") {
				break
			}
		}
	}
	hexMsg := helpers.BitsToHex(bits.String())
	return os.WriteFile(DetectFile, []byte(hexMsg), 0644)
}

// Extraction for Method 2: Check for single vs. double spaces
func extractMethod2(input []byte, DetectFile string) error {
	text := string(input)
	spaceRegex := regexp.MustCompile(`[ ]{1,2}`)
	spaces := spaceRegex.FindAllString(text, -1)
	var bits strings.Builder
	for _, sp := range spaces {
		if sp == "  " {
			bits.WriteByte('1')
		} else if sp == " " {
			bits.WriteByte('0')
		}
		if bits.Len()%4 == 0 && helpers.IsHex(bits.String()) {
			if helpers.BitsToHex(bits.String()) == helpers.ReadFileContent("files/mess.txt") {
				break
			}
		}
	}
	hexMsg := helpers.BitsToHex(bits.String())
	return os.WriteFile(DetectFile, []byte(hexMsg), 0644)
}

// Extraction for Method 3: Detect style attribute typos
func extractMethod3(input []byte, DetectFile string) error {
	text := string(input)
	var bits strings.Builder
	styleRegex := regexp.MustCompile(`<p[^>]*style="[^"]*"[^>]*>`)
	matches := styleRegex.FindAllString(text, -1)
	for _, m := range matches {
		if strings.Contains(m, "margin-botom") {
			bits.WriteByte('0')
		} else if strings.Contains(m, "lineheight") {
			bits.WriteByte('1')
		}
		if bits.Len()%4 == 0 && helpers.IsHex(bits.String()) {
			if helpers.BitsToHex(bits.String()) == helpers.ReadFileContent("files/mess.txt") {
				break
			}
		}
	}
	hexMsg := helpers.BitsToHex(bits.String())
	return os.WriteFile(DetectFile, []byte(hexMsg), 0644)
}

// Extraction for Method 4: Detect open-close-open <font> tag sequence
func extractMethod4(input []byte, DetectFile string) error {
	text := string(input)
	open := regexp.MustCompile(`(?i)<font[^>]*>`)
	sequence := open.FindAllString(text, -1)
	var bits strings.Builder
	for i := 0; i < len(sequence)-2; i++ {
		if sequence[i] == sequence[i+2] {
			bits.WriteByte('1')
			i += 2
		} else {
			bits.WriteByte('0')
		}
		if bits.Len()%4 == 0 && helpers.IsHex(bits.String()) {
			if helpers.BitsToHex(bits.String()) == helpers.ReadFileContent("files/mess.txt") {
				break
			}
		}
	}
	hexMsg := helpers.BitsToHex(bits.String())
	return os.WriteFile(DetectFile, []byte(hexMsg), 0644)
}
