// Package file provides a file interface
package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// GetContent is get the contents of a file
//  ex) content, err := file.GetContent(fileName)
func GetContent(fileName string) ([]string, error) {
	var lines []string

	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("no such file - (%s)", fileName))
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		return nil, err
	}

	return lines, nil
}

// SetContent is write content to file
//  ex) SetContent(fileName, content, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0600)
func SetContent(fileName string, content []string, flag int, mode uint32) error {
	file, err := os.OpenFile(fileName, flag, os.FileMode(mode))
	defer file.Close()
	if err != nil {
		return err
	}

	for _, value := range content {
		_, err = fmt.Fprintln(file, value)
		if err != nil {
			return err
		}
	}

	return nil
}
