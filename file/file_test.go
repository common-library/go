package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/file"
)

func TestRead(t *testing.T) {
	_, err := file.Read("./no_such_file")
	if err.Error() != "open ./no_such_file: no such file or directory" {
		t.Error(err)
	}

	readData, err := file.Read("./test.txt")
	if err != nil {
		t.Error(err)
	}

	compareData := []string{"aaa", "bbb", "ccc"}
	for index, value := range compareData {
		if value != readData[index] {
			t.Errorf("different data - (%s)(%s)", value, readData[index])
		}
	}
}

func TestWrite(t *testing.T) {
	const fileName string = "./temp.txt"
	writeData := []string{"abcdefg", "hijklmnop", "qrstuv"}
	flag := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	const mode uint32 = 0600

	err := file.Write(fileName, writeData, -1, mode)
	if err.Error() != "open ./temp.txt: no such file or directory" {
		t.Error(err)
	}

	err = file.Write(fileName, writeData, flag, mode)
	if err != nil {
		t.Error(err)
	}

	readData, err := file.Read(fileName)
	if err != nil {
		t.Error(err)
	}

	for index, value := range readData {
		if value != writeData[index] {
			t.Errorf("different data - (%s)(%s)", value, writeData[index])
		}
	}

	os.Remove(fileName)
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
	flag01 := int(os.O_WRONLY | os.O_APPEND | os.O_CREATE)
	if err := file.Write(file01, []string{}, flag01, 0600); err != nil {
		t.Fatal(err)
	}

	file02 := dir03 + uuid.New().String() + ".txt"
	flag02 := int(os.O_WRONLY | os.O_APPEND | os.O_CREATE)
	if err := file.Write(file02, []string{}, flag02, 0600); err != nil {
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
