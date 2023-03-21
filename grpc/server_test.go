package grpc

import (
	"testing"

	"github.com/heaven-chp/common-library-go/grpc/sample"
)

func TestInitialize(t *testing.T) {
	var server Server

	err := server.Initialize("1.1.1.1:50051", &sample.Server{})
	if err.Error() != "listen tcp 1.1.1.1:50051: bind: cannot assign requested address" {
		t.Error(err)
	}

	err = server.Initialize("127.0.0.1:50051", &sample.Server{})
	if err != nil {
		t.Error(err)
	}

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var server Server

	err := server.Finalize()
	if err != nil {
		t.Error(err)
	}

	err = server.Initialize("127.0.0.1:50051", &sample.Server{})
	if err != nil {
		t.Error(err)
	}

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestRun(t *testing.T) {
	var server Server

	err := server.Run()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = server.Initialize("127.0.0.1:50051", &sample.Server{})
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}
