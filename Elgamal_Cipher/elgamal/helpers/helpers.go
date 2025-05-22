// Author: Paulina Kimak
package helpers

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

// Function to count selected flags
func CountSelectedFlags(flags []*bool) int {
	count := 0
	for _, f := range flags {
		if *f {
			count++
		}
	}
	return count
}

//TODO: remove white signs before read from file
func ReadBigIntsFromFile(path string, count int) ([]*big.Int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var values []*big.Int
	for scanner.Scan() {
		line := scanner.Text()
		n := new(big.Int)
		n, ok := n.SetString(line, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number in %s", path)
		}
		values = append(values, n)
		if len(values) == count {
			break
		}
	}
	if len(values) != count {
		return nil, fmt.Errorf("expected %d numbers in %s", count, path)
	}
	return values, nil
}

func WriteBigIntsToFile(path string, values []*big.Int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, val := range values {
		_, err := file.WriteString(val.String() + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func RandomBigInt(max *big.Int) (*big.Int, error) {
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	return r, nil
}