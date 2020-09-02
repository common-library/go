package file

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	_, err := Read("./no_such_file")
	if err.Error() != "no such file - (./no_such_file)" {
		t.Error(err)
	}

	readData, err := Read("./test.txt")
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

	err := Write(fileName, writeData, flag, mode)
	if err != nil {
		t.Error(err)
	}

	readData, err := Read(fileName)
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
