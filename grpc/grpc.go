// Package grpc provides grpc interface.
// Used "google.golang.org/grpc".
// You need to include protobuf and create objects with ServerInfo as member variables.
// The sample is sample_server/sample_server.go
package grpc

import (
	"errors"
	"google.golang.org/grpc"
	"net"
)

// ServerInfo is object that provides server common infomation.
type ServerInfo struct {
	address string

	listener net.Listener

	GrpcServer *grpc.Server
}

// Initialize is initialize.
//  ex) server_info.Initialize("127.0.0.1:50051")
func (server_info *ServerInfo) Initialize(address string) error {
	server_info.address = address

	var err error
	server_info.listener, err = net.Listen("tcp", server_info.address)
	if err != nil {
		return err
	}

	server_info.GrpcServer = grpc.NewServer()
	if server_info.GrpcServer == nil {
		return errors.New("grpc.NewServer() fail")
	}

	return nil
}

// Finalize is finalize.
//  ex) server_info.Finalize()
func (server_info *ServerInfo) Finalize() error {
	if server_info.GrpcServer != nil {
		server_info.GrpcServer.Stop()
		server_info.GrpcServer = nil
	}

	if server_info.listener != nil {
		server_info.listener.Close()
		server_info.listener = nil
	}

	return nil
}

// Serve is grpc server serve.
// Note that it waits until Finalize() is called.
//  ex 1) server_info.Run()
//  ex 2) go server_info.Run()
func (server_info *ServerInfo) Serve() error {
	if server_info.GrpcServer == nil {
		return errors.New("please call Initialize first")
	}

	return server_info.GrpcServer.Serve(server_info.listener)
}
