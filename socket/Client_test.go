package socket

import (
	"testing"
)

const networkClientTest string = "tcp"
const addressClientTest string = "127.0.0.1:22222"

func TestDial(t *testing.T) {
	var client Client
	defer client.Close()

	err := client.Dial(networkClientTest, addressClientTest)
	if err.Error() != "dial tcp "+addressClientTest+": connect: connection refused" {
		t.Error(err)
	}
}

func TestReadWrite(t *testing.T) {
	var server Server
	defer server.Finalize()

	err := server.Initialize(networkClientTest, addressClientTest, 1024, nil)
	if err != nil {
		t.Error(err)
	}

	var channel chan Client = make(chan Client)

	go accept(t, &server, channel)

	var client Client
	defer client.Close()
	err = client.Dial(networkClientTest, addressClientTest)
	if err != nil {
		t.Error(err)
	}

	var serverClient Client = <-channel

	const writeContent string = "greeting"

	err = serverClient.Write(writeContent)
	if err != nil {
		t.Error(err)
	}

	readContent, err := client.Read(1024)
	if err != nil {
		t.Error(err)
	}

	if readContent != writeContent {
		t.Errorf("read error - writeContent : (%s), readContent : (%s)", writeContent, readContent)
	}
}
