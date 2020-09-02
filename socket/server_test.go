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
	err = client.Connect(networkServerTest, addressServerTest)
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
		const writeData string = "greeting"

		client.Write(writeData)

		readData, _ := client.Read(1024)

		client.Write(readData)
	}

	var server Server

	err := server.Initialize(networkServerTest, addressServerTest, 1024, jobFunc)
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	var client Client
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
