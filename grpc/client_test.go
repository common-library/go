package grpc_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/grpc"
	"github.com/heaven-chp/common-library-go/grpc/sample"
)

func TestGetConnection(t *testing.T) {
	var server grpc.Server
	err := server.Initialize("127.0.0.1:50051", &sample.Server{})
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	connection, err := grpc.GetConnection("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}
