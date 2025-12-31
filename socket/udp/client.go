// Package udp provides UDP socket client and server implementations.
//
// This package simplifies UDP network programming with packet-oriented
// communication. UDP is connectionless and does not guarantee delivery,
// ordering, or duplicate protection.
//
// # Features
//
//   - UDP socket client
//   - UDP socket server with packet handler
//   - Packet send and receive operations
//   - Concurrent packet processing (server)
//   - Timeout support
//   - Broadcast and multicast support
//   - Local and remote address access
//
// # Basic Client Example
//
//	client := &udp.Client{}
//	err := client.Connect("udp4", "localhost:8080")
//	client.Send([]byte("Hello"))
//	data, addr, _ := client.Receive(1024, 5*time.Second)
//	client.Close()
//
// # UDP Server Example
//
//	server := &udp.Server{}
//	err := server.Start("udp4", ":8080", 1024,
//	    func(data []byte, addr net.Addr, conn net.PacketConn) {
//	        fmt.Printf("Received from %s: %s\n", addr, data)
//	        conn.WriteTo(data, addr)  // Echo back
//	    },
//	    true,  // Async handler
//	    nil,   // No error handler
//	)
//	defer server.Stop()
package udp

import (
	"errors"
	"net"
	"time"
)

// Client is a struct that provides UDP client related methods.
type Client struct {
	conn net.PacketConn
	addr net.Addr
}

// Connect creates a UDP connection with specific network type.
//
// # Parameters
//
//   - network: Network type ("udp", "udp4", "udp6")
//   - address: Remote address (e.g., "localhost:8080")
//
// # Returns
//
//   - error: Error if connection fails, nil on success
//
// # Examples
//
//	client := &udp.Client{}
//	err := client.Connect("udp4", "localhost:8080")
func (c *Client) Connect(network, address string) error {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP(network, nil, addr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.addr = addr
	return nil
}

// Send sends data to the remote address.
//
// # Parameters
//
//   - data: Byte array to send
//
// # Returns
//
//   - int: Number of bytes sent
//   - error: Error if send fails, nil on success
//
// # Examples
//
//	n, err := client.Send([]byte("Hello, Server!"))
func (c *Client) Send(data []byte) (int, error) {
	if c.conn == nil {
		return 0, errors.New("please call Connect first")
	}

	// Use Write for connected UDP sockets
	if udpConn, ok := c.conn.(*net.UDPConn); ok {
		return udpConn.Write(data)
	}
	return c.conn.WriteTo(data, c.addr)
}

// SendTo sends data to a specific address.
//
// This method allows sending to a different address than the one specified
// in Connect(). For connected UDP sockets (using net.DialUDP), this creates
// a temporary unconnected socket for the send operation, which incurs some
// overhead. If you need to send to multiple different addresses frequently,
// consider using net.ListenPacket directly instead of Client.
//
// # Parameters
//
//   - data: Byte array to send
//   - address: Target address (can differ from Connect address)
//
// # Returns
//
//   - int: Number of bytes sent
//   - error: Error if send fails, nil on success
//
// # Examples
//
//	// Send to a different address than Connect
//	client.Connect("udp4", "localhost:8080")
//	n, err := client.SendTo([]byte("Hello"), "192.168.1.100:9000")
//
// # Performance Note
//
// WARNING: Each SendTo() call creates and closes a temporary socket, which
// has significant performance overhead. For high-performance scenarios with
// multiple destinations, use net.ListenPacket directly:
//
//	conn, _ := net.ListenPacket("udp", ":0")
//	defer conn.Close()
//	conn.WriteTo(data1, addr1)
//	conn.WriteTo(data2, addr2)  // No temporary socket created
func (c *Client) SendTo(data []byte, address string) (int, error) {
	if c.conn == nil {
		return 0, errors.New("please call Connect first")
	}

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return 0, err
	}

	// For connected UDP sockets, create a temporary unconnected socket
	// This is necessary because connected sockets cannot send to arbitrary addresses
	if _, ok := c.conn.(*net.UDPConn); ok {
		tempConn, err := net.ListenPacket("udp", ":0")
		if err != nil {
			return 0, err
		}
		defer tempConn.Close()
		return tempConn.WriteTo(data, addr)
	}

	return c.conn.WriteTo(data, addr)
}

// Receive receives data from any source.
//
// # Parameters
//
//   - bufferSize: Maximum bytes to read
//   - timeout: Read timeout duration (0 for no timeout)
//
// # Returns
//
//   - []byte: Received data
//   - net.Addr: Source address
//   - error: Error if receive fails, nil on success
//
// # Examples
//
//	data, addr, err := client.Receive(1024, 5*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Received from %s: %s\n", addr, data)
func (c *Client) Receive(bufferSize int, timeout time.Duration) ([]byte, net.Addr, error) {
	if c.conn == nil {
		return nil, nil, errors.New("please call Connect first")
	}

	if timeout > 0 {
		if err := c.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return nil, nil, err
		}
	}

	buffer := make([]byte, bufferSize)
	n, addr, err := c.conn.ReadFrom(buffer)
	if err != nil {
		return nil, nil, err
	}

	return buffer[:n], addr, nil
}

// Close closes the UDP connection.
//
// # Returns
//
//   - error: Error if close fails, nil on success
//
// # Examples
//
//	err := client.Close()
func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.addr = nil
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
	if c.conn == nil {
		return nil
	}
	return c.conn.LocalAddr()
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
	return c.addr
}

// SetReadBuffer sets the size of the operating system's receive buffer.
//
// # Parameters
//
//   - bytes: Buffer size in bytes
//
// # Returns
//
//   - error: Error if setting fails, nil on success
func (c *Client) SetReadBuffer(bytes int) error {
	if c.conn == nil {
		return errors.New("please call Connect first")
	}

	if udpConn, ok := c.conn.(*net.UDPConn); ok {
		return udpConn.SetReadBuffer(bytes)
	}
	return errors.New("connection is not a UDP connection")
}

// SetWriteBuffer sets the size of the operating system's transmit buffer.
//
// # Parameters
//
//   - bytes: Buffer size in bytes
//
// # Returns
//
//   - error: Error if setting fails, nil on success
func (c *Client) SetWriteBuffer(bytes int) error {
	if c.conn == nil {
		return errors.New("please call Connect first")
	}

	if udpConn, ok := c.conn.(*net.UDPConn); ok {
		return udpConn.SetWriteBuffer(bytes)
	}
	return errors.New("connection is not a UDP connection")
}
