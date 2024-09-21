package socket_test

import (
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/socket"
)

type TestServer struct {
	Network          string
	Address          string
	Greeting         string
	PrefixOfResponse string

	server socket.Server
}

func (this *TestServer) Start(t *testing.T) {
	t.Parallel()

	this.Network = "tcp"
	this.Address = ":" + strconv.Itoa(10000+rand.IntN(10000))
	this.Greeting = "greeting"
	this.PrefixOfResponse = "[response] "

	acceptSuccessFunc := func(client socket.Client) {
		if writeLen, err := client.Write(this.Greeting); err != nil {
			t.Fatal(err)
		} else if writeLen != len(this.Greeting) {
			t.Fatal(writeLen, ",", len(this.Greeting))
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}

		writeData := this.PrefixOfResponse + readData
		if writeLen, err := client.Write(writeData); err != nil {
			t.Fatal(err)
		} else if writeLen != len(writeData) {
			t.Fatal(writeLen, ",", len(writeData))
		}
	}

	acceptFailureFunc := func(err error) {
		t.Fatal(err)
	}

	if err := this.server.Start(this.Network, this.Address, 100, acceptSuccessFunc, acceptFailureFunc); err != nil {
		t.Fatal(err)
	}

	for this.server.GetCondition() == false {
		time.Sleep(100 * time.Millisecond)
	}
}

func (this *TestServer) Stop(t *testing.T) {
	if err := this.server.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestConnect(t *testing.T) {
	t.Parallel()

	const network = "tcp"
	const address = ":10001"

	var client socket.Client
	defer client.Close()

	if err := client.Connect("", address); err.Error() != "dial: unknown network " {
		t.Fatal(err)
	}

	if err := client.Connect("invalid", address); err.Error() != "dial invalid: unknown network invalid" {
		t.Fatal(err)
	}

	if err := client.Connect(network, ""); err.Error() != "dial tcp: missing address" {
		t.Fatal(err)
	}

	if err := client.Connect(network, "127.0.0.1"); err.Error() != "dial tcp: address 127.0.0.1: missing port in address" {
		t.Fatal(err)
	}

	if err := client.Connect(network, address); err.Error() != "dial tcp "+address+": connect: connection refused" {
		t.Fatal(err)
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

		if _, err := client.Read(1024); err.Error() != "please call the Connect function first" {
			t.Fatal(err)
		} else if _, err := client.Write(""); err.Error() != "please call the Connect function first" {
			t.Fatal(err)
		}

		if err := client.Connect(testServer.Network, testServer.Address); err != nil {
			t.Fatal(err)
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != testServer.Greeting {
			t.Fatal(readData, ",", testServer.Greeting)
		}

		writeData := "test"
		if writeLen, err := client.Write(writeData); err != nil {
			t.Fatal(err)
		} else if writeLen != len(writeData) {
			t.Fatal(writeLen, ",", len(writeData))
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != testServer.PrefixOfResponse+writeData {
			t.Fatal(writeData, ",", testServer.PrefixOfResponse+readData)
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
	t.Parallel()

	client := socket.Client{}

	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestGetLocalAddr(t *testing.T) {
	testServer := TestServer{}
	testServer.Start(t)
	defer testServer.Stop(t)

	client := socket.Client{}

	if addr := client.GetLocalAddr(); addr != nil {
		t.Fatal(addr)
	}

	if err := client.Connect(testServer.Network, testServer.Address); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if _, err := client.Read(1024); err != nil {
		t.Fatal(err)
	} else if _, err := client.Write("test"); err != nil {
		t.Fatal(err)
	} else if _, err := client.Read(1024); err != nil {
		t.Fatal(err)
	}

	if addr := client.GetLocalAddr(); addr == nil {
		t.Fatal(addr)
	} else if addr.Network() != testServer.Network {
		t.Fatal(addr.Network(), ",", testServer.Network)
	}
}

func TestGetRemoteAddr(t *testing.T) {
	testServer := TestServer{}
	testServer.Start(t)
	defer testServer.Stop(t)

	client := socket.Client{}

	if addr := client.GetRemoteAddr(); addr != nil {
		t.Fatal(addr)
	}

	if err := client.Connect(testServer.Network, testServer.Address); err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if _, err := client.Read(1024); err != nil {
		t.Fatal(err)
	} else if _, err := client.Write("test"); err != nil {
		t.Fatal(err)
	} else if _, err := client.Read(1024); err != nil {
		t.Fatal(err)
	}

	if addr := client.GetRemoteAddr(); addr == nil {
		t.Fatal(addr)
	} else if addr.Network() != testServer.Network {
		t.Fatal(addr.Network(), ",", testServer.Network)
	} else if strings.HasSuffix(addr.String(), testServer.Address) == false {
		t.Fatal(addr.String(), ",", testServer.Address)
	}
}
