// Package file provides file related implementations.
package file

import (
	"os"
	"path/filepath"
)

// Read is get the data of a file.
//
// ex) data, err := file.Read(fileName)
func Read(fileName string) (string, error) {
	if data, err := os.ReadFile(fileName); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Write is write data to file.
//
// ex 1) err := file.Write(fileName, data, 0600)
// ex 2) err := file.Write(fileName, data, os.ModePerm)
func Write(fileName string, data string, fileMode os.FileMode) error {
	return os.WriteFile(fileName, []byte(data), fileMode)
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
