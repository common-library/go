package socket_test

import (
	"sync"
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/socket"
)

type TestServer struct {
	Network          string
	Address          string
	Greeting         string
	PrefixOfResponse string

	server socket.Server
}

func (this *TestServer) Start(t *testing.T) {
	this.Network = "tcp"
	this.Address = ":10002"
	this.Greeting = "greeting"
	this.PrefixOfResponse = "[response] "

	acceptSuccessFunc := func(client socket.Client) {
		writeLen, err := client.Write(this.Greeting)
		if err != nil {
			t.Error(err)
		}
		if writeLen != len(this.Greeting) {
			t.Errorf("invalid write - (%d)(%d)", writeLen, len(this.Greeting))
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Error(err)
		}

		writeData := this.PrefixOfResponse + readData
		writeLen, err = client.Write(writeData)
		if err != nil {
			t.Error(err)
		}
		if writeLen != len(writeData) {
			t.Errorf("invalid write - (%d)(%d)", writeLen, len(writeData))
		}
	}

	acceptFailureFunc := func(err error) {
		t.Error(err)
	}

	err := this.server.Start(this.Network, this.Address, 100, acceptSuccessFunc, acceptFailureFunc)
	if err != nil {
		t.Fatal(err)
	}
	for this.server.GetCondition() == false {
		time.Sleep(100 * time.Millisecond)
	}
}

func (this *TestServer) Stop(t *testing.T) {
	err := this.server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

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
	testServer := TestServer{}
	testServer.Start(t)
	defer testServer.Stop(t)

	clientJob := func(wg *sync.WaitGroup) {
		defer wg.Done()

		client := socket.Client{}
		defer client.Close()

		_, err := client.Read(1024)
		if err.Error() != "please call the Connect function first" {
			t.Fatal(err)
		}

		_, err = client.Write("")
		if err.Error() != "please call the Connect function first" {
			t.Fatal(err)
		}

		err = client.Connect(testServer.Network, testServer.Address)
		if err != nil {
			t.Fatal(err)
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}
		if readData != testServer.Greeting {
			t.Fatalf("invalid read - (%s)(%s)", readData, testServer.Greeting)
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
		if readData != testServer.PrefixOfResponse+writeData {
			t.Fatalf("invalid read - (%s)(%s)", writeData, testServer.PrefixOfResponse+readData)
		}
	}

	wg := sync.WaitGroup{}
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go clientJob(&wg)
	}
	wg.Wait()
}

func TestClose(t *testing.T) {
	client := socket.Client{}

	err := client.Close()
	if err != nil {
		t.Error(err)
	}
}
