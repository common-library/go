package socket_test

import (
	"sync"
	"testing"

	"github.com/heaven-chp/common-library-go/socket"
)

func TestConnect(t *testing.T) {
	const network = "tcp"
	const address = ":10001"

	var client socket.Client
	defer client.Close()

	err := client.Connect("", address)
	if err.Error() != "dial: unknown network " {
		t.Error(err)
	}

	err = client.Connect("invalid", address)
	if err.Error() != "dial invalid: unknown network invalid" {
		t.Error(err)
	}

	err = client.Connect(network, "")
	if err.Error() != "dial tcp: missing address" {
		t.Error(err)
	}

	err = client.Connect(network, "127.0.0.1")
	if err.Error() != "dial tcp: address 127.0.0.1: missing port in address" {
		t.Error(err)
	}

	err = client.Connect(network, address)
	if err.Error() != "dial tcp "+address+": connect: connection refused" {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestReadWrite(t *testing.T) {
	const network = "tcp"
	const address = ":10002"
	const greeting = "greeting"
	const prefixOfResponse = "[response] "

	serverJob := func(client socket.Client) {
		writeLen, err := client.Write(greeting)
		if err != nil {
			t.Error(err)
		}
		if writeLen != len(greeting) {
			t.Errorf("invalid write - (%d)(%d)", writeLen, len(greeting))
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Error(err)
		}

		writeData := prefixOfResponse + readData
		writeLen, err = client.Write(writeData)
		if err != nil {
			t.Error(err)
		}
		if writeLen != len(writeData) {
			t.Errorf("invalid write - (%d)(%d)", writeLen, len(writeData))
		}
	}

	server := socket.Server{}

	go func() {
		err := server.Start(network, address, 100, serverJob)
		if err != nil {
			t.Error(err)
		}
	}()
	for server.GetCondition() == false {
	}

	clientJob := func(wg *sync.WaitGroup) {
		defer wg.Done()

		client := socket.Client{}
		defer client.Close()

		err := client.Connect(network, address)
		if err != nil {
			t.Fatal(err)
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}
		if readData != greeting {
			t.Fatalf("invalid read - (%s)(%s)", readData, greeting)
		}

		writeData := "test"
		writeLen, err := client.Write(writeData)
		if err != nil {
			t.Fatal(err)
		}
		if writeLen != len(writeData) {
			t.Fatalf("invalid write - (%d)(%d)", writeLen, len(writeData))
		}

		readData, err = client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}
		if readData != prefixOfResponse+writeData {
			t.Fatalf("invalid read - (%s)(%s)", writeData, prefixOfResponse+readData)
		}
	}

	wg := sync.WaitGroup{}
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go clientJob(&wg)
	}
	wg.Wait()

	err := server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClose(t *testing.T) {
	client := socket.Client{}

	err := client.Close()
	if err != nil {
		t.Error(err)
	}
}
