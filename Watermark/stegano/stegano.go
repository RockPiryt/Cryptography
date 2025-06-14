// Author: Paulina Kimak

package main

import (
	"flag"
	"log"

	"stegano/flagfunc"
	"stegano/helpers"
)

func main() {
	helpers.SetLogger()

	// Flags
	embedFlag := flag.Bool("e", false, "Embed message into cover.html and create watermark.html")
	extractFlag := flag.Bool("d", false, "Extract message from watermark.html")
	oneFlag := flag.Bool("1", false, "Method 1: bit as additional space at line end")
	twoFlag := flag.Bool("2", false, "Method 2: bit as single or double space")
	threeFlag := flag.Bool("3", false, "Method 3: bit as HTML attribute typo")
	fourFlag := flag.Bool("4", false, "Method 4: bit as redundant markup (e.g., FONT tags)")

	flag.Parse()

	// Validate operation selection
	operationFlags := []*bool{embedFlag, extractFlag}
	if helpers.CountSelectedFlags(operationFlags) != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -e or -d")
	}

	// Determine operation
	var operation string
	if *embedFlag {
		operation = "e"
	} else {
		operation = "d"
	}

	// Validate method selection (required in both cases now)
	methodFlags := []*bool{oneFlag, twoFlag, threeFlag, fourFlag}
	selectedMethods := helpers.CountSelectedFlags(methodFlags)

	if selectedMethods != 1 {
		log.Fatalf("Error: You must specify exactly one method (-1, -2, -3, or -4)")
	}

	// Determine selected method
	var method int
	switch {
	case *oneFlag:
		method = 1
	case *twoFlag:
		method = 2
	case *threeFlag:
		method = 3
	case *fourFlag:
		method = 4
	}
	// Create mess.txt with hex input
	var msg string = "bla"
	err := helpers.SaveHexToFile(msg, "files/mess.txt")
	if err != nil {
		log.Println("Error:", err)
	} else {
		log.Println("Hex saved to mess.txt")
	}

	// Execute main program logic
	err = flagfunc.ExecuteProgram(operation, method)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}
}
