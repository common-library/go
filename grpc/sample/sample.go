// Package sample provides sample grpc interface.
package sample

import (
	"context"

	"google.golang.org/grpc"
)

// Server is struct that satisfies the serverDetail interface and implements the function defined in protobuf IDL.
type Server struct {
	UnimplementedSampleServer
}

// RegisterServer is function to register in grpc server.
//
// auto call.
func (server *Server) RegisterServer(grpcServer *grpc.Server) {
	RegisterSampleServer(grpcServer, server)
}

// Func is implementation of the function defined in protobuf IDL.
//
// auto call.
func (server *Server) Func(context context.Context, request *Request) (*Reply, error) {
	return &Reply{Data1: request.Data1, Data2: request.Data2}, nil
}
