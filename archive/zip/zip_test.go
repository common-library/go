package zip_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/archive/zip"
	"github.com/heaven-chp/common-library-go/file"
)

func TestCompress(t *testing.T) {
	name := uuid.New().String() + string(filepath.Separator) + uuid.New().String() + ".zip"
	defer file.RemoveAll(filepath.Dir(name))

	input := uuid.New().String() + string(filepath.Separator)
	if err := file.CreateDirectory(input, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer file.RemoveAll(input)

	output := uuid.New().String() + string(filepath.Separator)
	defer file.RemoveAll(output)

	path01 := input + uuid.New().String() + ".txt"
	data01 := "aaa"
	if err := file.Write(path01, data01, 0600); err != nil {
		t.Fatal(err)
	}

	path02 := input + uuid.New().String() + ".txt"
	data02 := "bbb"
	if err := file.Write(path02, data02, 0600); err != nil {
		t.Fatal(err)
	}

	filePaths := []string{path01, path02}
	if err := zip.Compress(name, filePaths); err != nil {
		t.Fatal(err)
	}

	if err := zip.Decompress(name, output); err != nil {
		t.Fatal(err)
	} else if data, err := file.Read(output + path01); err != nil {
		t.Fatal(err)
	} else if data != data01 {
		t.Fatal("invalid data - ", data, ", ", data01)
	} else if data, err := file.Read(output + path02); err != nil {
		t.Fatal(err)
	} else if data != data02 {
		t.Fatal("invalid data - ", data, ", ", data02)
	}
}

func TestDecompress(t *testing.T) {
	TestCompress(t)
}
