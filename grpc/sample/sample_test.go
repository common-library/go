package sample

import (
	"context"
	"github.com/heaven-chp/common-library-go/grpc"
	"testing"
	"time"
)

func client_conn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	connection, err := grpc.GetConnection("127.0.0.1:50051")
	if err != nil {
		t.Error(err)
	}
	defer connection.Close()

	client := NewSampleClient(connection)

	const Data1 = 1
	const Data2 = "message"

	request := Request{Data1: Data1, Data2: Data2}
	reply, err := client.Func(ctx, &request)
	if err != nil {
		t.Error(err)
	}

	if reply.Data1 != Data1 || reply.Data2 != Data2 {
		t.Errorf("invalid reply - Data1 : (%d), Data2 : (%s)", reply.Data1, reply.Data2)
	}
}

func TestFunc(t *testing.T) {
	var server grpc.Server

	err := server.Initialize("127.0.0.1:50051", &Server{})
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
