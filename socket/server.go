// Package socket provides socket client and server implementations.
package socket

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Server is a struct that provides server related methods.
type Server struct {
	condition atomic.Bool

	listener          net.Listener
	channel           chan Client
	acceptSuccessFunc func(client Client)
	acceptFailureFunc func(err error)

	jobWaitGroup sync.WaitGroup
}

// Start is start the server.
//
// ex) err := server.Start("tcp", "127.0.0.1:10000", 1024, func(client Client) {...}, func(err error) {...})
func (this *Server) Start(network, address string, clientPoolSize int, acceptSuccessFunc func(client Client), acceptFailureFunc func(err error)) error {
	this.Stop()

	if len(network) == 0 {
		return errors.New("invalid network")
	}

	if len(address) == 0 {
		return errors.New("invalid address")
	}

	this.channel = make(chan Client, clientPoolSize)
	this.acceptSuccessFunc = acceptSuccessFunc
	this.acceptFailureFunc = acceptFailureFunc

	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	this.listener = listener

	go func() {
		this.condition.Store(true)
		for this.condition.Load() {
			client, err := this.accept()
			if err != nil {
				if this.condition.Load() && this.acceptFailureFunc != nil {
					this.acceptFailureFunc(err)
				}

				continue
			}

			this.channel <- client

			this.jobWaitGroup.Add(1)
			go this.job()
		}
	}()

	return nil
}

// Stop is stop the server.
//
// ex) err := server.Stop()
func (this *Server) Stop() error {
	this.condition.Store(false)

	this.jobWaitGroup.Wait()

	if this.channel != nil {
		for len(this.channel) != 0 {
			time.Sleep(time.Millisecond)
		}
		this.channel = nil
	}

	if this.listener != nil {
		err := this.listener.Close()
		this.listener = nil
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCondition is get the condition
//
// ex) condition := server.GetCondition()
func (this *Server) GetCondition() bool {
	return this.condition.Load()
}

func (this *Server) accept() (Client, error) {
	connnetion, err := this.listener.Accept()
	return Client{connnetion}, err
}

func (this *Server) job() {
	defer this.jobWaitGroup.Done()

	client := <-this.channel
	defer client.Close()

	if this.acceptSuccessFunc != nil {
		this.acceptSuccessFunc(client)
	}
}
