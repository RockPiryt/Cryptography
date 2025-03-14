package datafile

import (
	"bufio"
	"os"
)

//Author: Paulina Kimak

func GetText(filename string) ([]string, error) {
	var lines []string
	//otwarcie pliku
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	//odczytanie danych
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	if scanner.Err() != nil{
		return nil, scanner.Err()
	}
	
	return lines, nil
}


