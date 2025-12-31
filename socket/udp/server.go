// Package udp provides UDP socket server and client implementations.
//
// This package simplifies UDP network programming with packet-oriented
// communication. UDP is connectionless and does not guarantee delivery,
// ordering, or duplicate protection.
//
// # Features
//
//   - UDP socket server with packet handler
//   - UDP socket client
//   - Concurrent packet processing
//   - Graceful shutdown support
//   - Timeout support
//
// # Basic Server Example
//
//	server := &udp.Server{}
//	err := server.Start("udp4", ":8080", 1024, func(data []byte, addr net.Addr, conn net.PacketConn) {
//	    fmt.Printf("Received from %s: %s\n", addr, data)
//	    conn.WriteTo(data, addr)  // Echo back
//	}, true, nil)
//	defer server.Stop()
//
// # Basic Client Example
//
//	client := &udp.Client{}
//	err := client.Connect("udp4", "localhost:8080")
//	client.Send([]byte("Hello"))
//	data, addr, _ := client.Receive(1024, 5*time.Second)
//	client.Close()
package udp

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// PacketHandler is a function type for handling received UDP packets.
// It receives the packet data, source address, and the packet connection.
// Handlers can use the connection to send responses back to the source.
type PacketHandler func(data []byte, addr net.Addr, conn net.PacketConn)

// ErrorHandler is a function type for handling errors during packet reception.
// It receives the error that occurred during ReadFrom operation.
type ErrorHandler func(err error)

// Server is a struct that provides UDP server functionality.
// It handles incoming UDP packets and processes them using a PacketHandler.
type Server struct {
	mutex            sync.Mutex
	stop             atomic.Bool
	handlerWaitGroup sync.WaitGroup

	packetConn net.PacketConn
}

// Start starts the UDP server and begins listening for packets.
//
// The server runs in a goroutine and calls the handler for each received packet.
// The handler is responsible for processing the packet and can send responses
// using the provided PacketConn.
//
// # Parameters
//
//   - network: Network type ("udp", "udp4", "udp6")
//   - address: Local address to bind (e.g., ":8080", "0.0.0.0:8080")
//   - bufferSize: Maximum size for received packets in bytes
//   - handler: Function to handle received packets
//   - asyncHandler: If true, each packet handler runs in a separate goroutine,
//     enabling concurrent packet processing. If false, packets are processed
//     sequentially which may cause receive delays if handlers are slow.
//   - errorHandler: Optional function to handle read errors (can be nil)
//
// # Returns
//
//   - error: Error if server fails to start, nil on success
//
// # Examples
//
//	server := &udp.Server{}
//	err := server.Start("udp4", ":8080", 1024,
//	    func(data []byte, addr net.Addr, conn net.PacketConn) {
//	        log.Printf("Received %d bytes from %s\n", len(data), addr)
//	        conn.WriteTo(data, addr)  // Echo server
//	    },
//	    true,  // Run handlers concurrently
//	    func(err error) {
//	        log.Printf("Read error: %v\n", err)
//	    },
//	)
func (s *Server) Start(network, address string, bufferSize int, handler PacketHandler, asyncHandler bool, errorHandler ErrorHandler) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.packetConn != nil {
		return errors.New("server already started")
	}

	conn, err := net.ListenPacket(network, address)
	if err != nil {
		return err
	}
	s.packetConn = conn

	go func() {
		s.stop.Store(false)
		buffer := make([]byte, bufferSize)

		for !s.stop.Load() {
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				// Check if server is stopping
				if s.stop.Load() {
					break
				}
				// Call error handler if provided
				if errorHandler != nil {
					errorHandler(err)
				}
				continue
			}

			// Create a copy of the data for the handler
			data := make([]byte, n)
			copy(data, buffer[:n])

			if asyncHandler {
				// Run handler in goroutine for concurrent processing
				s.handlerWaitGroup.Add(1)
				go func(d []byte, a net.Addr) {
					defer s.handlerWaitGroup.Done()
					handler(d, a, conn)
				}(data, addr)
			} else {
				// Call handler synchronously
				handler(data, addr, conn)
			}
		}
	}()

	return nil
}

// Stop stops the UDP server and closes the connection.
//
// This method signals the server goroutine to stop and closes the underlying
// packet connection. It waits for all running handlers to complete before returning.
// It is safe to call Stop multiple times.
//
// # Returns
//
//   - error: Error if close fails, nil on success
//
// # Examples
//
//	err := server.Stop()
//	if err != nil {
//	    log.Printf("Error stopping server: %v", err)
//	}
func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stop.Store(true)

	if s.packetConn != nil {
		err := s.packetConn.Close()
		s.packetConn = nil

		// Wait for all handlers to complete
		s.handlerWaitGroup.Wait()

		return err
	}

	return nil
}

// GetLocalAddr returns the local network address the server is listening on.
//
// # Returns
//
//   - net.Addr: Local address, or nil if server is not started
//
// # Examples
//
//	addr := server.GetLocalAddr()
//	if addr != nil {
//	    fmt.Printf("Server listening on %s\n", addr)
//	}
func (s *Server) GetLocalAddr() net.Addr {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.packetConn != nil {
		return s.packetConn.LocalAddr()
	}
	return nil
}

// IsRunning returns whether the server is currently running.
//
// # Returns
//
//   - bool: true if server is running, false otherwise
//
// # Examples
//
//	if server.IsRunning() {
//	    fmt.Println("Server is running")
//	}
func (s *Server) IsRunning() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return !s.stop.Load() && s.packetConn != nil
}
