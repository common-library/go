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
func (s *Server) Start(network, address string, clientPoolSize int, acceptSuccessFunc func(client Client), acceptFailureFunc func(err error)) error {
	s.Stop()

	if len(network) == 0 {
		return errors.New("invalid network")
	}

	if len(address) == 0 {
		return errors.New("invalid address")
	}

	s.channel = make(chan Client, clientPoolSize)
	s.acceptSuccessFunc = acceptSuccessFunc
	s.acceptFailureFunc = acceptFailureFunc

	listener, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	s.listener = listener

	go func() {
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
	}()

	return nil
}

// Stop is stop the server.
//
// ex) err := server.Stop()
func (s *Server) Stop() error {
	s.condition.Store(false)

	s.jobWaitGroup.Wait()

	if s.channel != nil {
		for len(s.channel) != 0 {
			time.Sleep(time.Millisecond)
		}
		s.channel = nil
	}

	if s.listener != nil {
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCondition is get the condition
//
// ex) condition := server.GetCondition()
func (s *Server) GetCondition() bool {
	return s.condition.Load()
}

func (s *Server) accept() (Client, error) {
	connnetion, err := s.listener.Accept()
	return Client{connnetion}, err
}

func (s *Server) job() {
	defer s.jobWaitGroup.Done()

	client := <-s.channel
	defer client.Close()

	if s.acceptSuccessFunc != nil {
		s.acceptSuccessFunc(client)
	}
}
