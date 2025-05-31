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

	err := flagfunc.ExecuteCipher(operation)
	if err != nil {
		log.Fatalf("Execution error: %v", err)
	}


}
