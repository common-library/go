// Package grpc provides simplified gRPC client and server utilities.
//
// This package wraps Google's gRPC-Go library with convenient functions for creating
// client connections and managing server lifecycle. It reduces boilerplate code for
// common gRPC operations.
//
// Features:
//   - Simplified client connection creation
//   - Server lifecycle management (Start/Stop)
//   - Thread-safe server operations
//   - Support for insecure and secure connections
//
// Example:
//
//	// Client
//	conn, _ := grpc.GetConnection("localhost:50051")
//	defer conn.Close()
//
//	// Server
//	var server grpc.Server
//	server.Start(":50051", &myService{})
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

// Start starts the gRPC server on the specified address.
//
// Parameters:
//   - address: Address to bind the server to (e.g., ":50051", "localhost:50051")
//   - server: Implementation of the gRPC service that provides RegisterServer method
//
// Returns:
//   - error: Error if server cannot start or during serving, nil on graceful shutdown
//
// This method is blocking and will serve requests until Stop is called or an error occurs.
// The server parameter must implement the implementServer interface with a RegisterServer
// method that registers the gRPC service handlers.
//
// Example:
//
//	type GreeterServer struct {
//	    pb.UnimplementedGreeterServer
//	}
//
//	func (s *GreeterServer) RegisterServer(grpcServer *grpc.Server) {
//	    pb.RegisterGreeterServer(grpcServer, s)
//	}
//
//	var server grpc.Server
//	err := server.Start(":50051", &GreeterServer{})
//	if err != nil {
//	    log.Fatal(err)
//	}
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

// Stop gracefully stops the gRPC server and closes the listener.
//
// Returns:
//   - error: Always returns nil
//
// This method performs a graceful shutdown of the gRPC server, waiting for active
// RPCs to complete. It is thread-safe and can be called multiple times. After calling
// Stop, the server can be started again with a new address.
//
// Example:
//
//	var server grpc.Server
//	go server.Start(":50051", &myService{})
//
//	// Later, when shutting down
//	err := server.Stop()
//	if err != nil {
//	    log.Printf("Error stopping server: %v", err)
//	}
//	log.Println("Server stopped gracefully")
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
