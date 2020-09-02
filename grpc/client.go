package grpc

import (
	"google.golang.org/grpc"
)

// GetConnection is get connection of grpc
func GetConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
}
