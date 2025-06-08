// Author: Paulina Kimak
package flagfunc

import (
	"fmt"
	"log"
)

const (
	HtmlFile      		= "files/cover.html"
	MessageFile    		= "files/mess.txt"
	WatermarkFile    	= "files/watermatk.html"
	DetectFile 			= "files/detect.txt"
)


func ExecuteProgram(operation string, method int) error {
	switch operation {
	case "e":
		// IN MessageFile, OUT WatermarkFile 
		err := EmbedMsg(MessageFile, HtmlFile)
		if err != nil {
			return fmt.Errorf("failed to embeded the message: %v", err)
		}
		log.Println("[INFO] Text successfully embeded into watermark.html.")
		return nil

	case "d":
		err := ExtractMsg(WatermarkFile, DetectFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt the text: %v", err)
		}
		log.Println("[INFO] Text successfully decrypted into decrypt.txt.")
		return nil

	default:
		return fmt.Errorf("unsupported operation: %s", operation)
	}
}


func EmbedMsg(MessageFile, HtmlFile string) error {
	return nil
}

func ExtractMsg(WatermarkFile, DetectFile string) error {
	return nil
}