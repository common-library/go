package tar_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/common-library/go/archive/tar"
	"github.com/common-library/go/file"
	"github.com/google/uuid"
)

func TestCompress(t *testing.T) {
	name := uuid.New().String() + string(filepath.Separator) + uuid.New().String() + ".tar.gz"
	defer file.RemoveAll(filepath.Dir(name))

	input1 := uuid.New().String() + string(filepath.Separator)
	if err := file.CreateDirectory(input1, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer file.RemoveAll(input1)

	input2 := uuid.New().String() + string(filepath.Separator)
	defer file.RemoveAll(input2)
	input2 += input2
	if err := file.CreateDirectoryAll(input2, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	output := uuid.New().String() + string(filepath.Separator)
	defer file.RemoveAll(output)

	path01 := input1 + uuid.New().String() + ".txt"
	data01 := "aaa"
	if err := file.Write(path01, data01, 0600); err != nil {
		t.Fatal(err)
	}

	path02 := input2 + uuid.New().String() + ".txt"
	data02 := "bbb"
	if err := file.Write(path02, data02, 0600); err != nil {
		t.Fatal(err)
	}

	filePaths := []string{path01, path02}
	if err := tar.Compress(name, filePaths); err != nil {
		t.Fatal(err)
	}

	if err := tar.Decompress(name, output); err != nil {
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
