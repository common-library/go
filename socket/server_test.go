package socket_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/socket"
)

const networkServerTest string = "tcp"
const addressServerTest string = "127.0.0.1:11111"

func TestInitialize(t *testing.T) {
	var server socket.Server
	defer server.Finalize()

	err := server.Initialize(networkServerTest, addressServerTest, 1024, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var server socket.Server

	err := server.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	jobFunc := func(client socket.Client) {
		const writeData string = "greeting"

		client.Write(writeData)

		readData, _ := client.Read(1024)

		client.Write(readData)
	}

	var server socket.Server

	err := server.Initialize(networkServerTest, addressServerTest, 1024, jobFunc)
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	var client socket.Client
	defer client.Close()
	err = client.Connect(networkServerTest, addressServerTest)
	if err != nil {
		t.Error(err)
	}

	readData, err := client.Read(1024)
	if err != nil {
		t.Error(err)
	}

	var writeData string = "greeting"

	if readData != writeData {
		t.Errorf("read error - writeData : (%s), readData : (%s)", writeData, readData)
	}

	writeData = "12345"
	writeLen, err := client.Write(writeData)
	if err != nil {
		t.Error(err)
	}
	if writeLen != len(writeData) {
		t.Errorf("writeLen !=len(writeData) - writeLen : (%d), len(writeData) : (%d)", writeLen, len(writeData))
	}

	readData, err = client.Read(1024)
	if err != nil {
		t.Error(err)
	}

	if readData != writeData {
		t.Errorf("read error - writeData : (%s), readData : (%s)", writeData, readData)
	}

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}
