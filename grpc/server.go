// Package grpc provides grpc server interface.
//
// Used "google.golang.org/grpc".
package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type implementServer interface {
	RegisterServer(server *grpc.Server)
}

// Server is object that provides server common infomation.
type Server struct {
	listener net.Listener

	grpcServer *grpc.Server
}

// Start is start the server.
//
// ex) err := server.Start(":10000", &Sample.Server{})
func (this *Server) Start(address string, server implementServer) error {
	this.Stop()

	this.grpcServer = grpc.NewServer()

	server.RegisterServer(this.grpcServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	this.listener = listener

	return this.grpcServer.Serve(this.listener)
}

// Stop is stop the server.
//
// ex) err := server.Stop()
func (this *Server) Stop() error {
	if this.grpcServer != nil {
		this.grpcServer.Stop()
		this.grpcServer = nil
	}

	if this.listener != nil {
		this.listener.Close()
		this.listener = nil
	}

	return nil
}
