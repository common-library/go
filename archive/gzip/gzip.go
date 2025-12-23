// Package gzip provides utilities for creating and extracting gzip compressed files.
//
// This package wraps the compress/gzip standard library to provide convenient
// functions for compressing and decompressing single files.
//
// Features:
//   - Compress single files to gzip format
//   - Decompress gzip files to original format
//   - Automatic directory creation
//   - File permission preservation
//
// Example usage:
//
//	err := gzip.Compress("output.gz", "input.txt")
//	err := gzip.Decompress("archive.gz", "extracted.txt")
package gzip

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/common-library/go/file"
)

// Compress compresses a single file into gzip format.
//
// Parameters:
//   - name: output gzip file path (e.g., "test.gz")
//   - path: input file path to compress
//
// Example:
//
//	err := gzip.Compress("test.gz", "test.txt")
func Compress(name string, path string) error {
	if err := file.CreateDirectoryAll(filepath.Dir(name), os.ModePerm); err != nil {
		return err
	}

	gzipFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	gzipWriter := gzip.NewWriter(gzipFile)
	defer gzipWriter.Close()

	if source, err := os.Open(path); err != nil {
		return err
	} else {
		defer source.Close()

		_, err := io.Copy(gzipWriter, source)
		return err
	}
}

// Decompress decompresses a gzip file.
//
// Parameters:
//   - gzipName: input gzip file path (e.g., "test.gz")
//   - fileName: output file name (e.g., "test.txt")
//   - outputPath: output directory path
//
// Example:
//
//	err := gzip.Decompress("test.gz", "test.txt", "./output")
func Decompress(gzipName, fileName, outputPath string) error {
	gzipFile, err := os.Open(gzipName)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	if data, err := io.ReadAll(gzipReader); err != nil {
		return err
	} else if err := file.CreateDirectoryAll(outputPath, os.ModePerm); err != nil {
		return err
	} else {
		return file.Write(outputPath+string(filepath.Separator)+fileName, string(data), 0600)
	}
}
