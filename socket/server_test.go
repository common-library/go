package socket_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/socket"
)

func TestStart1(t *testing.T) {
	const network = "tcp"
	const address = ":10001"

	server := socket.Server{}

	err := server.Start("", address, 1024, nil, nil)
	if err.Error() != "invalid network" {
		t.Fatal(err)
	}

	err = server.Start(network, "", 1024, nil, nil)
	if err.Error() != "invalid address" {
		t.Fatal(err)
	}

	err = server.Start(network, "invalid_address", 1024, nil, nil)
	if err.Error() != "listen tcp: address invalid_address: missing port in address" {
		t.Fatal(err)
	}

	err = server.Start(network, "invalid_address:10000", 1024, nil, nil)
	if strings.HasPrefix(err.Error(), "listen tcp: lookup invalid_address on") == false {
		t.Fatal(err)
	}
}

func TestStart2(t *testing.T) {
	const network = "tcp"
	const address = ":10002"
	const greeting = "greeting"
	const prefixOfResponse = "[response] "

	acceptSuccessFunc := func(client socket.Client) {
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

	acceptFailureFunc := func(err error) {
		t.Error(err)
	}

	server := socket.Server{}
	err := server.Start(network, address, 100, acceptSuccessFunc, acceptFailureFunc)
	if err != nil {
		t.Error(err)
	}
	for server.GetCondition() == false {
		time.Sleep(100 * time.Millisecond)
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

	err = server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStop(t *testing.T) {
	server := socket.Server{}

	err := server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
