package udp_test

import (
	"net"
	"testing"
	"time"

	"github.com/common-library/go/socket/udp"
)

func TestConnect(t *testing.T) {
	t.Parallel()

	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", "localhost:9999"); err != nil {
		t.Fatal(err)
	}
}

func TestSendReceive(t *testing.T) {
	t.Parallel()

	// Create a UDP echo server
	serverConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer serverConn.Close()

	serverAddr := serverConn.LocalAddr().String()

	// Echo server goroutine
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, addr, err := serverConn.ReadFrom(buffer)
			if err != nil {
				return
			}
			serverConn.WriteTo(buffer[:n], addr)
		}
	}()

	// Client test
	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	testData := []byte("Hello, UDP!")
	if n, err := client.Send(testData); err != nil {
		t.Fatal(err)
	} else if n != len(testData) {
		t.Fatalf("sent %d bytes, expected %d", n, len(testData))
	}

	received, addr, err := client.Receive(1024, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if string(received) != string(testData) {
		t.Fatalf("received %q, expected %q", received, testData)
	}

	if addr.String() != serverAddr {
		t.Fatalf("received from %s, expected %s", addr, serverAddr)
	}
}

func TestSendTo(t *testing.T) {
	t.Parallel()

	// Create a UDP server
	serverConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer serverConn.Close()

	serverAddr := serverConn.LocalAddr().String()

	received := make(chan []byte, 1)
	go func() {
		buffer := make([]byte, 1024)
		n, _, err := serverConn.ReadFrom(buffer)
		if err != nil {
			return
		}
		received <- buffer[:n]
	}()

	// Client test
	client := &udp.Client{}
	defer client.Close()

	// Connect to a different address first
	if err := client.Connect("udp4", "localhost:9998"); err != nil {
		t.Fatal(err)
	}

	// Send to the actual server address
	testData := []byte("SendTo test")
	if n, err := client.SendTo(testData, serverAddr); err != nil {
		t.Fatal(err)
	} else if n != len(testData) {
		t.Fatalf("sent %d bytes, expected %d", n, len(testData))
	}

	select {
	case data := <-received:
		if string(data) != string(testData) {
			t.Fatalf("received %q, expected %q", data, testData)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for data")
	}
}

func TestReceiveTimeout(t *testing.T) {
	t.Parallel()

	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", "localhost:9997"); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, _, err := client.Receive(1024, 100*time.Millisecond)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if elapsed < 100*time.Millisecond {
		t.Fatalf("timeout happened too early: %v", elapsed)
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	client := &udp.Client{}

	// Close without connect should not error
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}

	// Connect and close
	if err := client.Connect("udp4", "localhost:9996"); err != nil {
		t.Fatal(err)
	}

	if err := client.Close(); err != nil {
		t.Fatal(err)
	}

	// Operations after close should fail
	if _, err := client.Send([]byte("test")); err == nil {
		t.Fatal("expected error after close")
	}
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	client := &udp.Client{}
	defer client.Close()

	// Before connect
	if addr := client.GetLocalAddr(); addr != nil {
		t.Fatal("local address should be nil before connect")
	}

	if addr := client.GetRemoteAddr(); addr != nil {
		t.Fatal("remote address should be nil before connect")
	}

	// After connect
	if err := client.Connect("udp4", "localhost:9995"); err != nil {
		t.Fatal(err)
	}

	if addr := client.GetLocalAddr(); addr == nil {
		t.Fatal("local address should not be nil after connect")
	}

	if addr := client.GetRemoteAddr(); addr == nil {
		t.Fatal("remote address should not be nil after connect")
	}
}

func TestSetBuffers(t *testing.T) {
	t.Parallel()

	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", "localhost:9994"); err != nil {
		t.Fatal(err)
	}

	if err := client.SetReadBuffer(4096); err != nil {
		t.Fatal(err)
	}

	if err := client.SetWriteBuffer(4096); err != nil {
		t.Fatal(err)
	}
}
