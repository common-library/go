// Package zip provides utilities for creating and extracting zip archives.
//
// This package wraps the archive/zip standard library to provide convenient
// functions for working with zip compressed files.
//
// Features:
//   - Compress multiple files and directories into zip
//   - Extract zip archives while preserving structure
//   - Recursive directory processing
//   - File mode preservation
//
// Example usage:
//
//	err := zip.Compress("archive.zip", []string{"./docs", "./images"})
//	err := zip.Decompress("archive.zip", "./extracted")
package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/common-library/go/file"
)

// Compress compresses multiple files and directories into zip format.
//
// Parameters:
//   - name: output zip file path (e.g., "test.zip")
//   - paths: slice of file/directory paths to compress
//
// The function recursively processes directories and preserves file modes.
//
// Example:
//
//	err := zip.Compress("test.zip", []string{"./test", "./test.txt"})
func Compress(name string, paths []string) error {
	filePaths := []string{}
	for _, path := range paths {
		if result, err := file.List(path, true); err != nil {
			return err
		} else {
			filePaths = append(filePaths, result...)
		}
	}

	if err := file.CreateDirectoryAll(filepath.Dir(name), os.ModePerm); err != nil {
		return err
	}

	zipFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	write := func(filePath string) error {
		source, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer source.Close()

		if destination, err := zipWriter.Create(filePath); err != nil {
			return err
		} else if _, err := io.Copy(destination, source); err != nil {
			return err
		}

		return nil
	}

	for _, filePath := range filePaths {
		if err := write(filePath); err != nil {
			return err
		}
	}

	return nil
}

// Decompress extracts a zip archive to the specified directory.
//
// Parameters:
//   - name: input zip file path (e.g., "test.zip")
//   - outputPath: output directory path where files will be extracted
//
// The function preserves directory structure and file modes.
//
// Example:
//
//	err := zip.Decompress("test.zip", "./output")
func Decompress(name, outputPath string) error {
	write := func(zipFile *zip.File) error {
		filePath := filepath.Join(outputPath, zipFile.Name)

		if zipFile.FileInfo().IsDir() {
			if err := file.CreateDirectoryAll(filePath, os.ModePerm); err != nil {
				return err
			} else {
				return nil
			}
		} else if err := file.CreateDirectoryAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		source, err := zipFile.Open()
		if err != nil {
			return err
		}
		defer source.Close()

		flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
		destination, err := os.OpenFile(filePath, flag, zipFile.Mode())
		if err != nil {
			return err
		}
		defer destination.Close()

		if _, err := io.Copy(destination, source); err != nil {
			return err
		}

		return nil
	}

	readCloser, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	for _, file := range readCloser.File {
		if err := write(file); err != nil {
			return err
		}
	}

	return nil
}
