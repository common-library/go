package grpc_test

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/common-library/go/grpc"
)

func TestGetConnection(t *testing.T) {
	connection, err := grpc.GetConnection("127.0.0.1:" + strconv.Itoa(10000+rand.IntN(10000)))

	if err != nil {
		t.Error(err)
	}

	defer connection.Close()
}
