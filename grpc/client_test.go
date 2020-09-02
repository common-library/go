package grpc

import (
	"github.com/heaven-chp/common-library-go/grpc/sample"
	"testing"
)

func TestGetConnection(t *testing.T) {
	var server Server
	err := server.Initialize("127.0.0.1:50051", &sample.Server{})
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	connection, err := GetConnection("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}
