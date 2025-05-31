// Author: Paulina Kimak
package main

import (
	"flag"
	"log"

	"rabin/flagfunc"
	"rabin/helpers"
)

func main() {
	helpers.SetLogger()

	//Set flags
	fermatFlag := flag.Bool("f", false, "Fermat test")

	flag.Parse()

	// Determine the operation
	var operation string
	if *fermatFlag {
		operation = "f"
	} else {
		operation = "r" // default Rabin-Miller
	}

	err := flagfunc.ExecuteProgram(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}

	err = flagfunc.FermatTest("files/wejscie.txt")
	if err != nil {
		log.Fatalf("Fermat test failed: %v", err)
	}

}
