// Package tcp provides TCP socket client and server implementations.
//
// This package simplifies network programming with high-level abstractions
// for socket servers and clients, supporting concurrent connection handling
// and automatic resource management.
//
// # Features
//
//   - TCP socket client
//   - Simple connect, read, write operations
//   - Automatic connection management
//   - Local and remote address access
//   - Resource cleanup
//
// # Basic Client Example
//
//	client := &tcp.Client{}
//	err := client.Connect("tcp", "localhost:8080")
//	client.Write("Hello")
//	data, _ := client.Read(1024)
//	client.Close()
package tcp

import (
	"errors"
	"net"
)

// Client is a struct that provides client related methods.
type Client struct {
	connection net.Conn
}

// Connect establishes a connection to the remote address.
//
// # Parameters
//
//   - network: Network type ("tcp", "tcp4", "tcp6", "unix")
//   - address: Remote address (e.g., "localhost:8080", "192.168.1.1:9000")
//
// # Returns
//
//   - error: Error if connection fails, nil on success
//
// # Examples
//
//	client := &socket.Client{}
//	err := client.Connect("tcp", "localhost:8080")
func (c *Client) Connect(network, address string) error {
	connection, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.connection = connection

	return nil
}

// Read reads data from the connection.
//
// # Parameters
//
//   - recvSize: Maximum bytes to read (buffer size)
//
// # Returns
//
//   - string: Received data
//   - error: Error if read fails, nil on success
//
// # Examples
//
//	data, err := client.Read(1024)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(data)
func (c *Client) Read(recvSize int) (string, error) {
	if c.connection == nil {
		return "", errors.New("please call Connect first")
	}

	buffer := make([]byte, recvSize)
	recvLen, err := c.connection.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:recvLen]), nil
}

// Write writes data to the connection.
//
// # Parameters
//
//   - data: Text data to write
//
// # Returns
//
//   - int: Number of bytes written
//   - error: Error if write fails, nil on success
//
// # Examples
//
//	n, err := client.Write("Hello, Server!")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Wrote %d bytes\n", n)
func (c *Client) Write(data string) (int, error) {
	if c.connection == nil {
		return -1, errors.New("please call Connect first")
	}

	return c.connection.Write([]byte(data))
}

// Close closes the connection.
//
// # Returns
//
//   - error: Error if close fails, nil on success
//
// # Examples
//
//	err := client.Close()
func (c *Client) Close() error {
	if c.connection != nil {
		err := c.connection.Close()
		c.connection = nil
		return err
	}

	return nil
}

// GetLocalAddr returns the local network address.
//
// # Returns
//
//   - net.Addr: Local address, or nil if not connected
//
// # Examples
//
//	addr := client.GetLocalAddr()
//	if addr != nil {
//	    fmt.Println(addr.String())
//	}
func (c *Client) GetLocalAddr() net.Addr {
	if c.connection == nil {
		return nil
	}

	return c.connection.LocalAddr()
}

// GetRemoteAddr returns the remote network address.
//
// # Returns
//
//   - net.Addr: Remote address, or nil if not connected
//
// # Examples
//
//	addr := client.GetRemoteAddr()
//	if addr != nil {
//	    fmt.Println(addr.String())
//	}
func (c *Client) GetRemoteAddr() net.Addr {
	if c.connection == nil {
		return nil
	}

	return c.connection.RemoteAddr()
}
