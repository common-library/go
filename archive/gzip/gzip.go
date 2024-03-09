// Package gzip provides gzip implementations.
package gzip

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/heaven-chp/common-library-go/file"
)

// Compress is compression.
//
// ex) err := gzip.Compress("test.gz", "test.txt")
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

// Decompress is decompression.
//
// ex) err := gzip.Compress("test.gz", "test.txt", "./output")
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

	if data, err := ioutil.ReadAll(gzipReader); err != nil {
		return err
	} else if err := file.CreateDirectoryAll(outputPath, os.ModePerm); err != nil {
		return err
	} else {
		return file.Write(outputPath+string(filepath.Separator)+fileName, string(data), 0600)
	}
}
