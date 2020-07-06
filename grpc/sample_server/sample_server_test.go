package sample_server

import (
	"context"
	"google.golang.org/grpc"
	"testing"
	"time"
)

/*
func TestListen(t *testing.T) {
    var server Server
    defer server.Finalize()

    err := server.listen(addressServerTest)
    if err != nil {
        t.Error(err)
    }
}
*/

func TestInitialize(t *testing.T) {
	var server Server

	err := server.Initialize("127.0.0.1:50051")
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
}

func TestRun(t *testing.T) {
	var server Server

	err := server.Initialize("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}

	go server.Run()

	client_conn(t)

	err = server.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func client_conn(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := NewSampleClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := client.SampleFunc(ctx, &SampleRequest{Id: 1, Msg: "abc"})
	if err != nil {
		t.Errorf("did not client.SayHello : %v", err)
	}

	if reply.Id != 1 || reply.Msg != "abc" {
		t.Errorf("invalid reply - Id : (%d), Msg : (%s)\n", reply.Id, reply.Msg)
	}
}
