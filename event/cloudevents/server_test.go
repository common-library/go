package cloudevents_test

import (
	"net/http"
	"testing"

	"github.com/common-library/go/event/cloudevents"
)

func TestStart(t *testing.T) {
	address, server := startServer(t)
	defer stopServer(t, server)

	if client, err := cloudevents.NewHttp("http://"+address, nil, nil); err != nil {
		t.Fatal(err)
	} else {
		for i := 0; i < 100; i++ {
			if result := client.Send(getEvent(t)); result.IsUndelivered() {
				t.Fatal(result.Error())
			} else if statusCode, err := result.GetHttpStatusCode(); err != nil {
				t.Fatal(err)
			} else if statusCode != http.StatusOK {
				t.Fatal("invalid -", statusCode)
			}
		}
	}
}

func TestStop(t *testing.T) {
	server := cloudevents.Server{}

	if err := server.Stop(10); err != nil {
		t.Fatal(err)
	}
}
