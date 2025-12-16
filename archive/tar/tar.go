// Package tar provides tar implementations.
package tar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/common-library/go/file"
)

// Compress is compression.
//
// ex) err := tar.Compress("test.tar.gz", []string{"./test", "./test.txt"})
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

	tarFile, err := os.Create(name)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzipWriter := gzip.NewWriter(tarFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	write := func(filePath string) error {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		if fileInfo, err := file.Stat(); err != nil {
			return err
		} else if header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name()); err != nil {
			return err
		} else {
			header.Name = filePath

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			} else if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			} else {
				return nil
			}
		}
	}

	for _, filePath := range filePaths {
		if err := write(filePath); err != nil {
			return err
		}
	}

	return nil
}

// Decompress is decompression.
//
// ex) err := tar.Decompress("test.tar.gz", "./output")
func Decompress(name, outputPath string) error {
	write := func(tarReader *tar.Reader, header *tar.Header) error {
		filePath := filepath.Join(outputPath, header.Name)

		switch header.Typeflag {
		case tar.TypeReg:
			if err := file.CreateDirectoryAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}

			flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
			file, err := os.OpenFile(filePath, flag, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		case tar.TypeDir:
			if err := file.CreateDirectoryAll(filePath, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	}

	gzipFile, err := os.Open(name)
	if err != nil {
		return err
	}
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		switch err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}

		if err := write(tarReader, header); err != nil {
			return err
		}
	}
}
