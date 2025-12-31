package udp_test

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/common-library/go/socket/udp"
)

func TestServerStartStop(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}

	// Start server
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		// Simple handler
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check server is listening
	addr := server.GetLocalAddr()
	if addr == nil {
		t.Fatal("server should have local address after start")
	}

	// Stop server
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}

	// Stop again should not error
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestServerAlreadyStarted(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}
	defer server.Stop()

	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Starting again should fail
	err = server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {}, false, nil)
	if err == nil {
		t.Fatal("expected error when starting already running server")
	}
}

func TestServerEchoPacket(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}
	defer server.Stop()

	// Start echo server
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		conn.WriteTo(data, addr)
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	// Create client
	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	// Send test data
	testData := []byte("Hello, Server!")
	if _, err := client.Send(testData); err != nil {
		t.Fatal(err)
	}

	// Receive echo
	received, _, err := client.Receive(1024, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if string(received) != string(testData) {
		t.Fatalf("received %q, expected %q", received, testData)
	}
}

func TestServerMultipleClients(t *testing.T) {
	t.Parallel()

	var receivedCount atomic.Int32
	server := &udp.Server{}
	defer server.Stop()

	// Start counting server
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		receivedCount.Add(1)
		conn.WriteTo(data, addr)
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	// Launch multiple clients
	const numClients = 10
	var wg sync.WaitGroup
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		go func(id int) {
			defer wg.Done()

			client := &udp.Client{}
			defer client.Close()

			if err := client.Connect("udp4", serverAddr); err != nil {
				t.Errorf("client %d connect failed: %v", id, err)
				return
			}

			testData := []byte("test message")
			if _, err := client.Send(testData); err != nil {
				t.Errorf("client %d send failed: %v", id, err)
				return
			}

			_, _, err := client.Receive(1024, 2*time.Second)
			if err != nil {
				t.Errorf("client %d receive failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Verify all packets were received
	if count := receivedCount.Load(); count != numClients {
		t.Fatalf("received %d packets, expected %d", count, numClients)
	}
}

func TestServerLargePacket(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}
	defer server.Stop()

	// Start echo server with 64KB buffer
	bufferSize := 65536
	err := server.Start("udp4", "127.0.0.1:0", bufferSize, func(data []byte, addr net.Addr, conn net.PacketConn) {
		conn.WriteTo(data, addr)
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	// Send large packet (but within typical UDP limits)
	testData := make([]byte, 8192)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	if _, err := client.Send(testData); err != nil {
		t.Fatal(err)
	}

	received, _, err := client.Receive(bufferSize, 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if len(received) != len(testData) {
		t.Fatalf("received %d bytes, expected %d", len(received), len(testData))
	}

	for i := range testData {
		if received[i] != testData[i] {
			t.Fatalf("data mismatch at byte %d", i)
		}
	}
}

func TestServerGetLocalAddr(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}

	// Before start
	if addr := server.GetLocalAddr(); addr != nil {
		t.Fatal("local address should be nil before start")
	}

	// After start
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {}, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	if addr := server.GetLocalAddr(); addr == nil {
		t.Fatal("local address should not be nil after start")
	}
}

func TestServerStopWhileReceiving(t *testing.T) {
	t.Parallel()

	var handlerCalled atomic.Bool
	server := &udp.Server{}

	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		handlerCalled.Store(true)
	}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	// Send one packet to verify handler works
	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	client.Send([]byte("test"))

	// Wait for handler to be called
	time.Sleep(100 * time.Millisecond)

	if !handlerCalled.Load() {
		t.Fatal("handler was not called")
	}

	// Stop server - should not hang
	done := make(chan error, 1)
	go func() {
		done <- server.Stop()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server stop timed out")
	}
}

func TestServerIsRunning(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}

	// Before start
	if server.IsRunning() {
		t.Fatal("server should not be running before start")
	}

	// After start
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !server.IsRunning() {
		t.Fatal("server should be running after start")
	}

	// After stop
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}

	if server.IsRunning() {
		t.Fatal("server should not be running after stop")
	}
}

func TestServerIsRunningConcurrent(t *testing.T) {
	t.Parallel()

	server := &udp.Server{}

	// Start server
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {}, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Test concurrent access to IsRunning
	const goroutines = 100
	const iterations = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = server.IsRunning()
				time.Sleep(time.Microsecond)
			}
		}()
	}

	// Stop server while goroutines are checking IsRunning
	time.Sleep(10 * time.Millisecond)
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func TestServerAsyncHandler(t *testing.T) {
	t.Parallel()

	var processedCount atomic.Int32
	var mu sync.Mutex
	handlerOrder := []int{}

	server := &udp.Server{}
	defer server.Stop()

	// Start server with async handlers
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		handlerOrder = append(handlerOrder, int(data[0]))
		mu.Unlock()

		processedCount.Add(1)
		conn.WriteTo(data, addr)
	}, true, nil) // asyncHandler = true
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	// Send multiple packets quickly
	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	const numPackets = 5
	for i := 0; i < numPackets; i++ {
		client.Send([]byte{byte(i)})
	}

	// Wait for all handlers to complete
	time.Sleep(200 * time.Millisecond)

	if count := processedCount.Load(); count != numPackets {
		t.Fatalf("processed %d packets, expected %d", count, numPackets)
	}
}

func TestServerSyncHandler(t *testing.T) {
	t.Parallel()

	var processedCount atomic.Int32
	var mu sync.Mutex
	handlerOrder := []int{}

	server := &udp.Server{}
	defer server.Stop()

	// Start server with sync handlers
	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		handlerOrder = append(handlerOrder, int(data[0]))
		mu.Unlock()

		processedCount.Add(1)
		conn.WriteTo(data, addr)
	}, false, nil) // asyncHandler = false
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	const numPackets = 5
	for i := 0; i < numPackets; i++ {
		client.Send([]byte{byte(i)})
		time.Sleep(5 * time.Millisecond) // Small delay between sends
	}

	// Wait for all handlers to complete
	time.Sleep(200 * time.Millisecond)

	if count := processedCount.Load(); count != numPackets {
		t.Fatalf("processed %d packets, expected %d", count, numPackets)
	}
}

func TestServerErrorHandler(t *testing.T) {
	t.Parallel()

	var errorCount atomic.Int32
	server := &udp.Server{}

	// Start server with error handler
	err := server.Start("udp4", "127.0.0.1:0", 1024,
		func(data []byte, addr net.Addr, conn net.PacketConn) {
			// Normal handler
		},
		false,
		func(err error) {
			errorCount.Add(1)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Server should be running normally
	if !server.IsRunning() {
		t.Fatal("server should be running")
	}

	// Stop server - this will cause read errors
	server.Stop()

	// Small delay to allow error handler to be called if needed
	time.Sleep(100 * time.Millisecond)

	// Note: Error count may or may not be > 0 depending on timing
	// This test mainly ensures error handler doesn't panic
}

func TestServerStopWaitsForHandlers(t *testing.T) {
	t.Parallel()

	var handlerStarted atomic.Bool
	var handlerCompleted atomic.Bool

	server := &udp.Server{}

	err := server.Start("udp4", "127.0.0.1:0", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
		handlerStarted.Store(true)
		time.Sleep(100 * time.Millisecond) // Slow handler
		handlerCompleted.Store(true)
	}, true, nil)
	if err != nil {
		t.Fatal(err)
	}

	serverAddr := server.GetLocalAddr().String()

	// Send packet
	client := &udp.Client{}
	defer client.Close()

	if err := client.Connect("udp4", serverAddr); err != nil {
		t.Fatal(err)
	}

	client.Send([]byte("test"))

	// Wait for handler to start
	time.Sleep(10 * time.Millisecond)
	if !handlerStarted.Load() {
		t.Fatal("handler should have started")
	}

	// Stop server - should wait for handler to complete
	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}

	// Handler should have completed
	if !handlerCompleted.Load() {
		t.Fatal("handler should have completed before Stop returned")
	}
}
