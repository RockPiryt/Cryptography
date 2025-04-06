package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/bits"
	"os"
	"os/exec"
	"strings"
)

// Function compareHexBitDiff compares two hex strings and returns the number of differing bits,
func compareHexBitDiff(hex1, hex2 string) (diffBits int, totalBits int, differencePercent float64, err error) {
    bytes1, err := hex.DecodeString(hex1)
    if err != nil {
        return 0, 0, 0, fmt.Errorf("error decoding hex1: %w", err)
    }
    bytes2, err := hex.DecodeString(hex2)
    if err != nil {
        return 0, 0, 0, fmt.Errorf("error decoding hex2: %w", err)
    }
    if len(bytes1) != len(bytes2) {
        return 0, 0, 0, fmt.Errorf("hex strings differ in length: %d vs %d bytes",
            len(bytes1), len(bytes2))
    }
    totalBits = len(bytes1) * 8

    for i := 0; i < len(bytes1); i++ {
        diff := bytes1[i] ^ bytes2[i]
        diffBits += bits.OnesCount8(diff)
    }
    differencePercent = float64(diffBits) / float64(totalBits) * 100
    return diffBits, totalBits, differencePercent, nil
}

// Function runCommand executes a shell command and returns its output as a string.
func runCommand(cmdStr string) (string, error) {
    cmd := exec.Command("sh", "-c", cmdStr)
    out, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("error running command '%s': %w", cmdStr, err)
    }
    return string(out), nil
}

//Function parseHashFromOutput extracts the hash from the output of a hash command.
func parseHashFromOutput(output string) (string, error) {
    fields := strings.Fields(output)
    if len(fields) < 1 {
        return "", fmt.Errorf("unexpected output (no fields): %q", output)
    }
    hashHex := fields[0]

    hashHex = strings.TrimSpace(hashHex)
    return hashHex, nil
}

// Funcction compareAllHashes compares the hashes of two files using different hash commands and writes the results to a file.
func compareAllHashes(hashFile, diffFile *os.File, hashCommands []string, file1, file2 string) error {
    for _, hashCmd := range hashCommands {
        // First command
        cmd1 := fmt.Sprintf("cat %s | %s", file1, hashCmd)
        fmt.Fprintln(hashFile, cmd1) 

        out1, err := runCommand(cmd1)
        if err != nil {
            return fmt.Errorf("error running cmd1: %w", err)
        }
        hash1, err := parseHashFromOutput(out1)
        if err != nil {
            return fmt.Errorf("error parsing hash1: %w", err)
        }

        // Second command
        cmd2 := fmt.Sprintf("cat %s | %s", file2, hashCmd)
        fmt.Fprintln(hashFile, cmd2) 

        out2, err := runCommand(cmd2)
        if err != nil {
            return fmt.Errorf("error running cmd2: %w", err)
        }
        hash2, err := parseHashFromOutput(out2)
        if err != nil {
            return fmt.Errorf("error parsing hash2: %w", err)
        }

        diffBits, totalBits, diffPercent, err := compareHexBitDiff(hash1, hash2)
        if err != nil {
            log.Fatalf("Error comparing hex: %v", err)
        }
		// Save data to hash.txt
        fmt.Fprintln(hashFile, hash1)
        fmt.Fprintln(hashFile, hash2)
		fmt.Fprintln(hashFile)


        // Save data to diff.txt
        fmt.Fprintln(diffFile, cmd1)
        fmt.Fprintln(diffFile, cmd2)
        fmt.Fprintln(diffFile, hash1)
        fmt.Fprintln(diffFile, hash2)
        fmt.Fprintf(diffFile, "Number of differing bits: %d out of %d (%.2f%%)\n\n",
            diffBits, totalBits, diffPercent)
    }

    return nil
}

func main() {
    hashFile, err := os.Create("hash.txt")
    if err != nil {
        fmt.Println("Could not create hash.txt:", err)
        return
    }
    defer hashFile.Close()

    diffFile, err := os.Create("diff.txt")
    if err != nil {
        fmt.Println("Could not create diff.txt:", err)
        return
    }
    defer diffFile.Close()

    hashCommands := []string{
        "md5sum",
        "sha1sum",
        "sha224sum",
        "sha256sum",
        "sha384sum",
        "sha512sum",
        "b2sum",
    }

    file1 := "hash-.pdf personal.txt"
    file2 := "hash-.pdf personal_.txt"

    if err := compareAllHashes(hashFile, diffFile, hashCommands, file1, file2); err != nil {
        fmt.Println("Error comparing all hashes:", err)
        return
    }
}
