// Package sample_server provides sample grpc server
package sample_server

import (
	"context"
	"github.com/heaven-chp/common-library-go/grpc"
)

// Server is object that specialized information
type Server struct {
	server_info grpc.ServerInfo
}

// Initialize is initialize.
//  ex) server.Initialize("127.0.0.1:50051")
func (server *Server) Initialize(address string) error {
	return server.server_info.Initialize(address)
}

// Finalize is finalize.
//  ex) server.Finalize()
func (server *Server) Finalize() error {
	return server.server_info.Finalize()
}

// Run is server run.
//  ex) server.Run()
func (server *Server) Run() {
	RegisterSampleServer(server.server_info.GrpcServer, server)

	server.server_info.Serve()
}

// SampleFunc is specialized function.
// auto call.
func (server *Server) SampleFunc(ctx context.Context, in *SampleRequest) (*SampleReply, error) {
	return &SampleReply{Id: in.Id, Msg: in.Msg}, nil
}
