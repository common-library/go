package file_test

import (
	"os"
	"testing"

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
