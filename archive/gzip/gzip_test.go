package gzip_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/common-library/go/archive/gzip"
	"github.com/common-library/go/file"
	"github.com/google/uuid"
)

func TestCompress(t *testing.T) {
	name := uuid.New().String() + string(filepath.Separator) + uuid.New().String() + ".gz"
	defer file.RemoveAll(filepath.Dir(name))

	input := uuid.New().String() + string(filepath.Separator)
	if err := file.CreateDirectory(input, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer file.RemoveAll(input)

	output := uuid.New().String() + string(filepath.Separator)
	defer file.RemoveAll(output)

	path := input + uuid.New().String() + ".txt"
	data := "aaa"
	if err := file.Write(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	if err := gzip.Compress(name, path); err != nil {
		t.Fatal(err)
	}

	fileName := uuid.New().String() + ".txt"
	if err := gzip.Decompress(name, fileName, output); err != nil {
		t.Fatal(err)
	} else if result, err := file.Read(output + fileName); err != nil {
		t.Fatal(err)
	} else if result != data {
		t.Fatal("invalid data - ", result, ", ", data)
	}
}

func TestDecompress(t *testing.T) {
	TestCompress(t)
}
