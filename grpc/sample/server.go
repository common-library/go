// Package sample provides sample grpc interface.
package sample

import (
	"context"
	"io"

	"google.golang.org/grpc"
)

// Server is struct that satisfies the implementServer interface and implements the function defined in protobuf IDL.
type Server struct {
	UnimplementedSampleServer
}

// RegisterServer is function to register in grpc server.
//
// the function we need to call is Register{service name of proto file}Server
func (s *Server) RegisterServer(grpcServer *grpc.Server) {
	RegisterSampleServer(grpcServer, s)
}

// Func1 is implementation of the function defined in protobuf IDL.
func (s *Server) Func1(context context.Context, request *Request) (*Reply, error) {
	return &Reply{Data1: request.Data1, Data2: request.Data2}, nil
}

// Func2 is implementation of the function defined in protobuf IDL.
func (s *Server) Func2(stream Sample_Func2Server) error {
	for {
		// Check context cancellation
		if err := stream.Context().Err(); err != nil {
			return err
		}

		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for range 3 {
			err := stream.Send(&Reply{Data1: request.Data1, Data2: request.Data2})
			if err != nil {
				return err
			}
		}
	}
}
