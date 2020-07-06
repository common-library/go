package socket

import (
	"net"
)

// Client is object that provides client infomation.
type Client struct {
	network    string
	address    string
	connnetion net.Conn
}

// Dial is connect to the address.
//  ex) client.Dial("127.0.0.1:22222")
func (client *Client) Dial(network string, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	client.network = network
	client.address = address
	client.connnetion = conn

	return nil
}

// Read is read data from connection.
//  ex) readContent, err := client.Read(1024)
func (client *Client) Read(recvSize int) (string, error) {
	buffer := make([]byte, recvSize)

	recvLen, err := client.connnetion.Read(buffer)
	if err != nil {
		return "", nil
	}

	return string(buffer[:recvLen]), nil
}

// Write is write data to connection.
//  ex) serverClient.Write("example")
func (client *Client) Write(content string) error {
	_, err := client.connnetion.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

// Close is close the connection.
//  ex) client.Close()
func (client *Client) Close() error {
	if client.connnetion != nil {
		err := client.connnetion.Close()
		client.connnetion = nil
		if err != nil {
			return err
		}
	}

	return nil
}
