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
	server := grpc.Server{}

	err := server.Start("1.1.1.1:10000", &sample.Server{})

	if err.Error() != "listen tcp 1.1.1.1:10000: bind: cannot assign requested address" {
		t.Error(err)
	}

	go func() {
		err = server.Start(":"+strconv.Itoa(10000+rand.IntN(10000)), &sample.Server{})
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(200 * time.Millisecond)

	err = server.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestStop(t *testing.T) {
	server := grpc.Server{}

	err := server.Stop()
	if err != nil {
		t.Error(err)
	}
}
