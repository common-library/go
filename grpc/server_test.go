package grpc_test

import (
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
)

func TestStart(t *testing.T) {
	t.Parallel()

	server := grpc.Server{}

	err := server.Start("1.1.1.1:10000", &sample.Server{})
	if err.Error() != "listen tcp 1.1.1.1:10000: bind: cannot assign requested address" {
		t.Fatal(err)
	}

	go func() {
		if err := server.Start(":"+strconv.Itoa(10000+rand.IntN(10000)), &sample.Server{}); err != nil {
			t.Errorf("server.Start failed: %v", err)
			return
		}
	}()
	time.Sleep(200 * time.Millisecond)

	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	server := grpc.Server{}

	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
}
