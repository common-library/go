package sample_test

import (
	"context"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
)

func setUp(server *grpc.Server) {
	go func() {
		err := server.Start(":10000", &sample.Server{})
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
}

func tearDown(server *grpc.Server) {
	err := server.Stop()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	server := grpc.Server{}

	setUp(&server)

	code := m.Run()

	tearDown(&server)

	os.Exit(code)
}

func TestFunc1(t *testing.T) {
	connection, err := grpc.GetConnection(":10000")
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()

	client := sample.NewSampleClient(connection)

	for i := 0; i < 10; i++ {
		data1 := int64(i)
		data2 := "message " + strconv.Itoa(i)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		reply, err := client.Func1(ctx, &sample.Request{Data1: data1, Data2: data2})
		if err != nil {
			t.Fatal(err)
		}

		if reply.Data1 != data1 || reply.Data2 != data2 {
			t.Errorf("invalid reply - (%d)(%d)(%s)(%s)", data1, reply.Data1, data2, reply.Data2)
		}
	}
}

func TestFunc2(t *testing.T) {
	connection, err := grpc.GetConnection(":10000")
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()

	client := sample.NewSampleClient(connection)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := client.Func2(ctx)
	if err != nil {
		t.Fatal(err)
	}

	const data1 = 1
	const data2 = "message"

	{
		err := stream.Send(&sample.Request{Data1: data1, Data2: data2})
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			reply, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Error(err)
			}

			if reply.Data1 != data1 || reply.Data2 != data2 {
				t.Errorf("invalid reply - (%d)(%d)(%s)(%s)", data1, reply.Data1, data2, reply.Data2)
			}
		}
	}

	go func() {
		for i := 0; i < 3; i++ {
			err := stream.Send(&sample.Request{Data1: data1, Data2: data2})
			if err != nil {
				t.Error(err)
			}
		}

		stream.CloseSend()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			reply, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				t.Error(err)
			}

			if reply.Data1 != data1 || reply.Data2 != data2 {
				t.Errorf("invalid reply - (%d)(%d)(%s)(%s)", data1, reply.Data1, data2, reply.Data2)
			}
		}
	}()
	wg.Wait()
}
