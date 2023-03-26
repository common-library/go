// Package socket provides a socket server interface.
package socket

import (
	"errors"
	"net"
	"sync"
	"time"
)

// Server is object that provides server infomation.
type Server struct {
	isRun   bool
	network string
	address string

	listener net.Listener
	channel  chan Client
	jobFunc  func(client Client)

	jobWaitGroup sync.WaitGroup
}

// Initialize is initialize.
//
// ex) server.Initialize("tcp", "127.0.0.1:11111", 1024, func(client Client) {})
func (server *Server) Initialize(network string, address string, clientPoolSize int, jobFunc func(client Client)) error {
	server.Finalize()

	server.isRun = true
	server.network = network
	server.address = address
	server.channel = make(chan Client, clientPoolSize)
	server.jobFunc = jobFunc

	return server.listen()
}

// Finalize is finalize.
//
// server.Finalize()
func (server *Server) Finalize() error {
	server.isRun = false

	var client Client
	client.Connect(server.network, server.address)
	client.Close()

	server.jobWaitGroup.Add(1)
	server.jobWaitGroup.Done()
	server.jobWaitGroup.Wait()

	if server.channel != nil {
		for len(server.channel) != 0 {
			time.Sleep(time.Millisecond)
		}
	}

	if server.listener == nil {
		return nil
	}

	err := server.listener.Close()
	server.listener = nil
	return err
}

// Run is server run.
//
// Note that it waits until Finalize() is called.
//
// ex 1) server.Run()
//
// ex 2) go server.Run()
func (server *Server) Run() error {
	for server.isRun {
		client, err := server.accept()
		if err != nil {
			return err
		}

		server.channel <- client

		server.jobWaitGroup.Add(1)
		go server.job()
	}

	return nil
}

func (server *Server) listen() error {
	if len(server.network) == 0 {
		return errors.New("invalid network")
	}

	if len(server.address) == 0 {
		return errors.New("invalid address")
	}

	var err error
	server.listener, err = net.Listen(server.network, server.address)
	return err
}

func (server *Server) accept() (Client, error) {
	if server.listener == nil {
		return Client{}, errors.New("please call Initialize first")
	}

	connnetion, err := server.listener.Accept()
	return Client{connnetion}, err
}

func (server *Server) job() {
	defer server.jobWaitGroup.Done()

	var client Client = <-server.channel

	if server.jobFunc != nil {
		server.jobFunc(client)
	}

	client.Close()
}
