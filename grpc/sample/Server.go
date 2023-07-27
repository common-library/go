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
func (this *Server) RegisterServer(grpcServer *grpc.Server) {
	RegisterSampleServer(grpcServer, this)
}

// Func1 is implementation of the function defined in protobuf IDL.
func (this *Server) Func1(context context.Context, request *Request) (*Reply, error) {
	return &Reply{Data1: request.Data1, Data2: request.Data2}, nil
}

// Func2 is implementation of the function defined in protobuf IDL.
func (this *Server) Func2(stream Sample_Func2Server) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		for i := 0; i < 3; i++ {
			err := stream.Send(&Reply{Data1: request.Data1, Data2: request.Data2})
			if err != nil {
				return err
			}
		}
	}

	return nil
	// return &Reply{Data1: request.Data1, Data2: request.Data2}, nil
}
