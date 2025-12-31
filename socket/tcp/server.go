// Package tcp provides TCP socket client and server implementations.
//
// This package simplifies network programming with high-level abstractions
// for socket servers and clients, supporting concurrent connection handling
// and automatic resource management.
//
// # Features
//
//   - TCP socket server
//   - Concurrent client connection handling
//   - Client connection pooling
//   - Automatic resource cleanup
//   - Custom accept success/failure handlers
//
// # Basic Server Example
//
//	server := &tcp.Server{}
//	err := server.Start("tcp", ":8080", 100,
//	    func(client tcp.Client) {
//	        data, _ := client.Read(1024)
//	        client.Write("Response: " + data)
//	    },
//	    func(err error) {
//	        log.Printf("Accept error: %v", err)
//	    },
//	)
package tcp

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// Server is a struct that provides server related methods.
type Server struct {
	condition atomic.Bool

	listener net.Listener
	channel  chan Client

	acceptSuccessFunc func(client Client)
	acceptFailureFunc func(err error)

	acceptWaitGroup sync.WaitGroup
	jobWaitGroup    sync.WaitGroup
}

// Start initializes and starts the socket server.
//
// This method creates a listener on the specified network and address,
// accepting incoming connections and handling them concurrently using
// a pool of goroutines.
//
// # Parameters
//
//   - network: Network type ("tcp", "tcp4", "tcp6", "unix")
//   - address: Listen address (e.g., ":8080", "127.0.0.1:8080")
//   - clientPoolSize: Channel buffer size for pending connections. This limits
//     how many accepted connections can be queued before blocking the accept loop.
//     Does not limit total concurrent connections.
//   - acceptSuccessFunc: Callback for each accepted connection
//   - acceptFailureFunc: Callback for accept errors (can be nil)
//
// # Returns
//
//   - error: Error if server start fails, nil on success
//
// # Behavior
//
// The server:
//   - Stops any existing server instance
//   - Creates a listener on the specified address
//   - Accepts connections in a background goroutine
//   - Spawns a goroutine for each connection (up to clientPoolSize buffered)
//   - Calls acceptSuccessFunc for each connection
//   - Calls acceptFailureFunc for accept errors (if not nil)
//
// # Examples
//
// Basic TCP echo server:
//
//	server := &socket.Server{}
//	err := server.Start("tcp", ":8080", 100,
//	    func(client socket.Client) {
//	        data, _ := client.Read(1024)
//	        client.Write("Echo: " + data)
//	    },
//	    func(err error) {
//	        log.Printf("Error: %v", err)
//	    },
//	)
func (s *Server) Start(network, address string, clientPoolSize int, acceptSuccessFunc func(client Client), acceptFailureFunc func(err error)) error {
	s.Stop()

	if len(network) == 0 {
		return errors.New("invalid network")
	}

	if len(address) == 0 {
		return errors.New("invalid address")
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	s.listener = listener
	s.channel = make(chan Client, clientPoolSize)
	s.acceptSuccessFunc = acceptSuccessFunc
	s.acceptFailureFunc = acceptFailureFunc

	s.acceptWaitGroup.Add(1)
	go func() {
		defer s.acceptWaitGroup.Done()
		s.condition.Store(true)
		for s.condition.Load() {
			client, err := s.accept()
			if err != nil {
				if s.condition.Load() && s.acceptFailureFunc != nil {
					s.acceptFailureFunc(err)
				}
				continue
			}

			s.channel <- client

			s.jobWaitGroup.Add(1)
			go s.job()
		}
		close(s.channel)
	}()

	return nil
}

// Stop gracefully shuts down the socket server.
//
// This method stops accepting new connections, waits for all active
// client handlers to complete, and closes the listener.
//
// # Returns
//
//   - error: Error if shutdown fails, nil on successful shutdown
//
// # Behavior
//
// The shutdown process:
//  1. Sets condition to false (stops accepting)
//  2. Waits for all job goroutines to complete
//  3. Waits for channel to drain
//  4. Closes the listener
//
// # Examples
//
//	err := server.Stop()
//	if err != nil {
//	    log.Printf("Stop error: %v", err)
//	}
func (s *Server) Stop() error {
	s.condition.Store(false)

	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	s.acceptWaitGroup.Wait()
	s.jobWaitGroup.Wait()

	return nil
}

// GetCondition returns the server running state.
//
// # Returns
//
//   - bool: true if server is running, false if stopped
//
// # Examples
//
//	if server.GetCondition() {
//	    fmt.Println("Server is running")
//	}
func (s *Server) GetCondition() bool {
	return s.condition.Load()
}

func (s *Server) accept() (Client, error) {
	connection, err := s.listener.Accept()
	return Client{connection: connection}, err
}

func (s *Server) job() {
	defer s.jobWaitGroup.Done()

	client := <-s.channel
	defer client.Close()

	if s.acceptSuccessFunc != nil {
		s.acceptSuccessFunc(client)
	}
}
