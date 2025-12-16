// Package grpc provides grpc client and server implementations.
package grpc

import (
	"net"
	"sync"

	"google.golang.org/grpc"
)

type implementServer interface {
	RegisterServer(server *grpc.Server)
}

// Server is struct that provides server common infomation.
type Server struct {
	mutex      sync.RWMutex
	listener   net.Listener
	grpcServer *grpc.Server
}

// Start is start the server.
//
// ex) err := server.Start(":10000", &Sample.Server{})
func (grpcSrv *Server) Start(address string, server implementServer) error {
	if err := func() error {
		grpcSrv.mutex.Lock()
		defer grpcSrv.mutex.Unlock()

		grpcServer := grpc.NewServer()
		server.RegisterServer(grpcServer)

		listener, err := net.Listen("tcp", address)
		if err != nil {
			return err
		}
		grpcSrv.grpcServer = grpcServer
		grpcSrv.listener = listener

		return nil
	}(); err != nil {
		return err
	}

	// Serve는 blocking 함수이므로 뮤텍스 해제 후 호출
	return grpcSrv.grpcServer.Serve(grpcSrv.listener)
}

// Stop is stop the server.
//
// ex) err := server.Stop()
func (grpcSrv *Server) Stop() error {
	grpcSrv.mutex.Lock()
	defer grpcSrv.mutex.Unlock()

	if grpcSrv.grpcServer != nil {
		grpcSrv.grpcServer.Stop()
		grpcSrv.grpcServer = nil
	}

	if grpcSrv.listener != nil {
		grpcSrv.listener.Close()
		grpcSrv.listener = nil
	}

	return nil
}
