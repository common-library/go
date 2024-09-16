package socket_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/socket"
)

func TestStart1(t *testing.T) {
	t.Parallel()

	const network = "tcp"
	const address = ":10001"

	server := socket.Server{}

	if err := server.Start("", address, 1024, nil, nil); err.Error() != "invalid network" {
		t.Fatal(err)
	}

	if err := server.Start(network, "", 1024, nil, nil); err.Error() != "invalid address" {
		t.Fatal(err)
	}

	if err := server.Start(network, "invalid_address", 1024, nil, nil); err.Error() != "listen tcp: address invalid_address: missing port in address" {
		t.Fatal(err)
	}

	if err := server.Start(network, "invalid_address:10000", 1024, nil, nil); strings.HasPrefix(err.Error(), "listen tcp: lookup invalid_address on") == false {
		t.Fatal(err)
	}
}

func TestStart2(t *testing.T) {
	t.Parallel()

	const network = "tcp"
	const address = ":10002"
	const greeting = "greeting"
	const prefixOfResponse = "[response] "

	acceptSuccessFunc := func(client socket.Client) {
		if writeLen, err := client.Write(greeting); err != nil {
			t.Fatal(err)
		} else if writeLen != len(greeting) {
			t.Fatal(writeLen, ",", len(greeting))
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}

		writeData := prefixOfResponse + readData
		if writeLen, err := client.Write(writeData); err != nil {
			t.Fatal(err)
		} else if writeLen != len(writeData) {
			t.Fatal(writeLen, ",", len(writeData))
		}
	}

	acceptFailureFunc := func(err error) {
		t.Fatal(err)
	}

	server := socket.Server{}
	if err := server.Start(network, address, 100, acceptSuccessFunc, acceptFailureFunc); err != nil {
		t.Fatal(err)
	}
	for server.GetCondition() == false {
		time.Sleep(100 * time.Millisecond)
	}
	defer func() {
		if err := server.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	clientJob := func(wg *sync.WaitGroup) {
		defer wg.Done()

		client := socket.Client{}
		defer client.Close()

		if err := client.Connect(network, address); err != nil {
			t.Fatal(err)
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != greeting {
			t.Fatal(readData, ",", greeting)
		}

		writeData := "test"
		if writeLen, err := client.Write(writeData); err != nil {
			t.Fatal(err)
		} else if writeLen != len(writeData) {
			t.Fatal(writeLen, ",", len(writeData))
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != prefixOfResponse+writeData {
			t.Fatal(writeData, ",", prefixOfResponse+readData)
		}
	}

	wg := sync.WaitGroup{}
	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go clientJob(&wg)
	}
	wg.Wait()
}

func TestStop(t *testing.T) {
	t.Parallel()

	server := socket.Server{}

	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
}
