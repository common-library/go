package socket

import (
	"testing"
)

const network string = "tcp"
const address string = "127.0.0.1:22222"

func TestConnect(t *testing.T) {
	var client Client
	defer client.Close()

	err := client.Connect(network, address)
	if err.Error() != "dial tcp "+address+": connect: connection refused" {
		t.Error(err)
	}
}

func TestReadWrite(t *testing.T) {
	var server Server
	defer server.Finalize()

	err := server.Initialize(network, address, 1024, nil)
	if err != nil {
		t.Error(err)
	}

	var channel chan Client = make(chan Client)

	go accept(t, &server, channel)

	var client Client
	defer client.Close()
	err = client.Connect(network, address)
	if err != nil {
		t.Error(err)
	}

	var serverClient Client = <-channel

	const writeData string = "greeting"

	writeLen, err := serverClient.Write(writeData)
	if err != nil {
		t.Error(err)
	}
	if writeLen != len(writeData) {
		t.Errorf("write len is different - writeLen : (%d), len(writeData) : (%d)", writeLen, len(writeData))
	}

	readData, err := client.Read(1024)
	if err != nil {
		t.Error(err)
	}
	if readData != writeData {
		t.Errorf("read error - writeData : (%s), readData : (%s)", writeData, readData)
	}
}
