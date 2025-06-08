// Author: Paulina Kimak
package main

import (
	"flag"
	"fmt"
	"log"

	"stegano/flagfunc"
	"stegano/helpers"
)

func main() {
	helpers.SetLogger()

	//Set flags
	embedFlag := flag.Bool("e", false, "embed")
	extractFlag := flag.Bool("d", false, "extract")
	oneFlag := flag.Bool("1", false, "bit as addtional space at the end")
	twoFlag := flag.Bool("2", false, "bit as two spaces")
	threeFlag := flag.Bool("3", false, "bit as error in html attributes")
	fourFlag := flag.Bool("4", false, "bit as addtional markup")

	flag.Parse()

	// Check operation flags 
	operationFlags := []*bool{embedFlag,extractFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)

	if operationCount != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -e or -d")
	}

	// Determine the operation
	var operation string
	switch {
	case *encryptFlag:
		operation = "e"
	case *decryptFlag:
		operation = "d"
	default:
		log.Fatalf("Error: Invalid operation selected.")
	}

	//var option string
	if operation == "d" {
		fmt.Print("Tryb zanurzenia")
	}

	// Check operation flags 
	operationFlags := []*bool{embedFlag,extractFlag}
	operationCount := helpers.CountSelectedFlags(operationFlags)

	if operationCount != 1 {
		log.Fatalf("Error: You must choose exactly one operation: -e or -d")
	}

	// Determine the operation
	var operation string
	switch {
	case *encryptFlag:
		operation = "e"
	case *decryptFlag:
		operation = "d"
	default:
		log.Fatalf("Error: Invalid operation selected.")
	}

	//var option string
	if operation == "d" {
		fmt.Print("Tryb zanurzenia")
	}

	err := flagfunc.ExecuteProgram(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}

}
