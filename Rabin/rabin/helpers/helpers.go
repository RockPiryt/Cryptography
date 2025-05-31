// Author: Paulina Kimak
package helpers

import (
	"log"
	"os"
)

// Function to set logger.
func SetLogger() {
	os.MkdirAll("logs", os.ModePerm)
	logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetPrefix("[Rabin-Miller/Fermat] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
