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
		if err := server.Start(":10000", &sample.Server{}); err != nil {
			panic(err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
}

func tearDown(server *grpc.Server) {
	if err := server.Stop(); err != nil {
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
	t.Parallel()

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

		if reply, err := client.Func1(ctx, &sample.Request{Data1: data1, Data2: data2}); err != nil {
			t.Fatal(err)
		} else if reply.Data1 != data1 || reply.Data2 != data2 {
			t.Fatal(reply)
		}
	}
}

func TestFunc2(t *testing.T) {
	t.Parallel()

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

	if err := stream.Send(&sample.Request{Data1: data1, Data2: data2}); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		if reply, err := stream.Recv(); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		} else if reply.Data1 != data1 || reply.Data2 != data2 {
			t.Fatal(reply)
		}
	}

	go func() {
		for i := 0; i < 3; i++ {
			if err := stream.Send(&sample.Request{Data1: data1, Data2: data2}); err != nil {
				t.Fatal(err)
			}
		}

		stream.CloseSend()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if reply, err := stream.Recv(); err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			} else if reply.Data1 != data1 || reply.Data2 != data2 {
				t.Fatal(reply)
			}
		}
	}()
	wg.Wait()
}
