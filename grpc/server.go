// Package grpc provides grpc interface.
// Used "google.golang.org/grpc".
package grpc

import (
	"errors"
	"google.golang.org/grpc"
	"net"
)

type serverDetail interface {
	RegisterServer(server *grpc.Server)
}

// Server is object that provides server common infomation.
type Server struct {
	address string

	listener net.Listener

	serverDetail serverDetail

	grpcServer *grpc.Server
}

// Initialize is initialize.
//  ex) server.Initialize("127.0.0.1:50051", &Sample.Server{})
func (server *Server) Initialize(address string, serverDetail serverDetail) error {
	server.address = address
	server.serverDetail = serverDetail

	var err error
	server.listener, err = net.Listen("tcp", server.address)
	if err != nil {
		return err
	}

	server.grpcServer = grpc.NewServer()
	if server.grpcServer == nil {
		return errors.New("grpc.NewServer() fail")
	}

	server.serverDetail.RegisterServer(server.grpcServer)

	return nil
}

// Finalize is finalize.
//  ex) server.Finalize()
func (server *Server) Finalize() error {
	if server.grpcServer != nil {
		server.grpcServer.Stop()
		server.grpcServer = nil
	}

	if server.listener != nil {
		server.listener.Close()
		server.listener = nil
	}

	return nil
}

// Run is server run.
// Note that it waits until Finalize() is called.
//  ex 1) server.Run()
//  ex 2) go server.Run()
func (server *Server) Run() error {
	if server.grpcServer == nil {
		return errors.New("please call Initialize first")
	}

	return server.grpcServer.Serve(server.listener)
}
