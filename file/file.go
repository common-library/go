// Package file provides a file interface.
// Package file provides file related implementations.
package file

import (
	"bufio"
	"fmt"
	"os"
)

// Read is get the data of a file.
//
// ex) data, err := file.Read(fileName)
func Read(fileName string) ([]string, error) {
	var lines []string

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Write is write data to file.
//
// ex) Write(fileName, data, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0600)
func Write(fileName string, data []string, flag int, mode uint32) error {
	file, err := os.OpenFile(fileName, flag, os.FileMode(mode))
	defer file.Close()
	if err != nil {
		return err
	}

	for _, value := range data {
		_, err = fmt.Fprintln(file, value)
		if err != nil {
			return err
		}
	}

	return nil
}
