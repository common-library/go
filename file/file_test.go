package file

import (
	"os"
	"testing"
)

func TestSetGetContent(t *testing.T) {
	const fileName string = "./test.txt"
	content := []string{"abcdefg", "hijklmnop", "qrstuv"}
	flag := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	const mode uint32 = 0600

	err := SetContent(fileName, content, flag, mode)
	if err != nil {
		t.Error(err)
	}

	compareContent, err := GetContent(fileName)
	if err != nil {
		t.Error(err)
	}

	for index, value := range compareContent {
		if value != content[index] {
			t.Errorf("different content - (%s)(%s)", value, content[index])
		}
	}

	os.Remove(fileName)
}
