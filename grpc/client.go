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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetConnection creates a new gRPC client connection to the specified address.
//
// Parameters:
//   - address: Server address in "host:port" format (e.g., "localhost:50051")
//
// Returns:
//   - *grpc.ClientConn: Client connection that can be used to create service clients
//   - error: Error if connection cannot be established, nil on success
//
// This function creates an insecure connection (no TLS). For production environments,
// consider using secure connections with proper credentials. The caller is responsible
// for closing the connection using conn.Close().
//
// Example:
//
//	conn, err := grpc.GetConnection("localhost:50051")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close()
//
//	// Use connection to create service client
//	client := pb.NewGreeterClient(conn)
//	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})
func GetConnection(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
