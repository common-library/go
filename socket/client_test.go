package socket_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/socket"
)

const network string = "tcp"
const address string = "127.0.0.1:22222"

func TestConnect(t *testing.T) {
	var client socket.Client
	defer client.Close()

	err := client.Connect(network, address)
	if err.Error() != "dial tcp "+address+": connect: connection refused" {
		t.Error(err)
	}
}

func TestReadWrite(t *testing.T) {
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
