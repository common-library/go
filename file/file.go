// Package file provides a file interface.
// Package file provides file related implementations.
package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
// ex) err := file.Write(fileName, data, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0600)
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

// List is get file list.
//
// ex) list, err := file.List(dir, true)
func List(path string, recursive bool) ([]string, error) {
	var list func(*[]string, string, bool) error
	list = func(result *[]string, path string, recursive bool) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if fileInfo, err := file.Stat(); err != nil {
			return err
		} else if fileInfo.IsDir() == false {
			*result = append(*result, path)
			return nil
		} else if fileInfos, err := file.Readdir(0); err != nil {
			return err
		} else if len(fileInfos) == 0 {
			*result = append(*result, file.Name()+string(filepath.Separator))
			return nil
		} else {
			for _, fileInfo := range fileInfos {
				dir := filepath.Dir(file.Name() + string(filepath.Separator))
				name := dir + string(filepath.Separator) + fileInfo.Name()

				if fileInfo.IsDir() == false {
					*result = append(*result, name)
				} else if recursive == false {
					*result = append(*result, name+string(filepath.Separator))
				} else {
					list(result, name, recursive)
				}
			}

			return nil
		}
	}

	result := []string{}
	if err := list(&result, path, recursive); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
