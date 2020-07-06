package socket

import (
	"testing"
)

const networkServerTest string = "tcp"
const addressServerTest string = "127.0.0.1:11111"

func accept(t *testing.T, server *Server, channel chan Client) {
	Client, err := (*server).accept()
	if err != nil {
		t.Error(err)
	}

	channel <- Client
}

func TestListen(t *testing.T) {
	var server Server
	defer server.Finalize()

	err := server.listen()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}
}

func TestAccept1(t *testing.T) {
	var server Server
	defer server.Finalize()

	_, err := server.accept()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}
}

func TestAccept2(t *testing.T) {
	var server Server
	defer server.Finalize()

	err := server.Initialize(networkServerTest, addressServerTest, 1024, nil)
	if err != nil {
		t.Error(err)
	}

	var channel chan Client = make(chan Client)

	go accept(t, &server, channel)

	var client Client
	defer client.Close()
	err = client.Dial(networkServerTest, addressServerTest)
	if err != nil {
		t.Error(err)
	}

	<-channel
}

func TestInitialize(t *testing.T) {
	var server Server
	defer server.Finalize()

	err := server.Initialize(networkServerTest, addressServerTest, 1024, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var server Server

	err := server.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	jobFunc := func(client Client) {
		const writeContent string = "greeting\n"

		client.Write(writeContent)

		readContent, _ := client.Read(1024)

		client.Write(readContent)
	}

	var server Server

	err := server.Initialize(networkServerTest, addressServerTest, 1024, jobFunc)
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	var client Client
	defer client.Close()
	err = client.Dial(networkServerTest, addressServerTest)
	if err != nil {
		t.Error(err)
	}

	readContent, err := client.Read(1024)
	if err != nil {
		t.Error(err)
	}

	var writeContent string = "greeting\n"

	if readContent != writeContent {
		t.Errorf("read error - writeContent : (%s), readContent : (%s)", writeContent, readContent)
	}

	writeContent = "12345\n"
	err = client.Write(writeContent)
	if err != nil {
		t.Error(err)
	}

	readContent, err = client.Read(1024)
	if err != nil {
		t.Error(err)
	}

	if readContent != writeContent {
		t.Errorf("read error - writeContent : (%s), readContent : (%s)", writeContent, readContent)
	}

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}
