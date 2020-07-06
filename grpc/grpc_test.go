package grpc

import (
	"testing"
)

func TestInitialize(t *testing.T) {
	var server_info ServerInfo

	err := server_info.Initialize("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}

	err = server_info.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestFinalize(t *testing.T) {
	var server_info ServerInfo

	err := server_info.Finalize()
	if err != nil {
		t.Error(err)
	}

	err = server_info.Initialize("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}

	err = server_info.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestServe(t *testing.T) {
	var server_info ServerInfo

	err := server_info.Serve()
	if err.Error() != "please call Initialize first" {
		t.Error(err)
	}

	err = server_info.Initialize("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}

	go server_info.Serve()

	err = server_info.Finalize()
	if err != nil {
		t.Error(err)
	}
}
