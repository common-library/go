package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/file"
)

func TestRead(t *testing.T) {
	fileName := uuid.New().String()
	defer os.Remove(fileName)

	if _, err := file.Read(fileName); os.IsExist(err) {
		t.Fatal(err)
	}

	const answer = "aaa\nbbb\nccc\n"

	if err := file.Write(fileName, answer, 0600); err != nil {
		t.Fatal(err)
	}

	if data, err := file.Read(fileName); err != nil {
		t.Fatal(err)
	} else if data != answer {
		t.Fatalf("invalid - (%s)", data)
	}
}

func TestWrite(t *testing.T) {
	TestRead(t)
}

func TestList(t *testing.T) {
	dir01 := uuid.New().String() + string(filepath.Separator)
	defer os.RemoveAll(dir01)

	dir02 := dir01 + uuid.New().String() + string(filepath.Separator)
	if err := os.MkdirAll(dir02, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	dir03 := dir01 + uuid.New().String() + string(filepath.Separator)
	if err := os.MkdirAll(dir03, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	file01 := dir01 + uuid.New().String() + ".txt"
	if err := file.Write(file01, "", 0600); err != nil {
		t.Fatal(err)
	}

	file02 := dir03 + uuid.New().String() + ".txt"
	if err := file.Write(file02, "", 0600); err != nil {
		t.Fatal(err)
	}

	answer := map[string]bool{}

	answer = map[string]bool{dir02: false, dir03: false, file01: false}
	if list, err := file.List(dir01, false); err != nil {
		t.Fatal(err)
	} else {
		if len(list) != len(answer) {
			t.Fatal("invalid list -", list)
		}

		for _, name := range list {
			if _, exist := answer[name]; exist == false {
				t.Fatal("invalid name :", name, list)
			}
		}
	}

	answer = map[string]bool{dir02: true, file01: true, file02: true}
	if list, err := file.List(dir01, true); err != nil {
		t.Fatal(err)
	} else {
		if len(list) != len(answer) {
			t.Fatal("invalid list -", list)
		}

		for _, name := range list {
			if _, exist := answer[name]; exist == false {
				t.Fatal("invalid name :", name, list)
			}
		}
	}
}
