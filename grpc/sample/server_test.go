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
		reply, err := client.Func1(ctx, &sample.Request{Data1: data1, Data2: data2})
		cancel()

		if err != nil {
			t.Fatal(err)
		}
		if reply.Data1 != data1 || reply.Data2 != data2 {
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stream, err := client.Func2(ctx)
	if err != nil {
		t.Fatalf("failed to create stream: %v", err)
	}

	const data1 = 1
	const data2 = "message"

	if err := stream.Send(&sample.Request{Data1: data1, Data2: data2}); err != nil {
		t.Fatalf("failed to send request: %v", err)
	}

	for i := 0; i < 3; i++ {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to receive reply at iteration %d: %v", i, err)
		}
		if reply.Data1 != data1 || reply.Data2 != data2 {
			t.Fatalf("unexpected reply at iteration %d: got %v, want Data1=%d Data2=%s", i, reply, data1, data2)
		}
	}

	go func() {
		for range 3 {
			if err := stream.Send(&sample.Request{Data1: data1, Data2: data2}); err != nil {
				t.Errorf("stream.Send failed: %v", err)
				return
			}
		}

		if err := stream.CloseSend(); err != nil {
			t.Errorf("stream.CloseSend failed: %v", err)
		}
	}()

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if reply, err := stream.Recv(); err == io.EOF {
				return
			} else if err != nil {
				t.Errorf("stream.Recv failed: %v", err)
				return
			} else if reply.Data1 != data1 || reply.Data2 != data2 {
				t.Errorf("unexpected reply: %v", reply)
				return
			}
		}
	}()
	wg.Wait()
}
