package cloudevents_test

import (
	"context"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"

	cloudeventssdk_http "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/common-library/go/event/cloudevents"
)

func TestSend(t *testing.T) {
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

func TestRequest(t *testing.T) {
	address, server := startServer(t)
	defer stopServer(t, server)

	if client, err := cloudevents.NewHttp("http://"+address, nil, nil); err != nil {
		t.Fatal(err)
	} else {
		for i := 0; i < 100; i++ {
			if event, result := client.Request(getEvent(t)); result.IsUndelivered() {
				t.Fatal(result.Error())
			} else if statusCode, err := result.GetHttpStatusCode(); err != nil {
				t.Fatal(err)
			} else if statusCode != http.StatusOK {
				t.Fatal("invalid -", statusCode)
			} else {
				consistencyEvent(t, event)
			}
		}
	}
}

func TestStartReceiver(t *testing.T) {
	port := rand.IntN(1000) + 10000
	httpOption := []cloudeventssdk_http.Option{cloudeventssdk_http.WithPort(port)}
	handler := func(ctx context.Context, event cloudevents.Event) {
		consistencyEvent(t, &event)
	}
	failureFunc := func(err error) { t.Fatal(err) }
	if receiveClient, err := cloudevents.NewHttp("", httpOption, nil); err != nil {
		t.Fatal(err)
	} else {
		receiveClient.StartReceiver(handler, failureFunc)

		for i := 0; i < 100; i++ {
			if sendCient, err := cloudevents.NewHttp("http://:"+strconv.Itoa(port), nil, nil); err != nil {
				t.Fatal(err)
			} else if result := sendCient.Send(getEvent(t)); result.IsUndelivered() {
				t.Fatal(result.Error())
			}
		}

		receiveClient.StopReceiver()
	}
}

func TestStopReceiver(t *testing.T) {
	if client, err := cloudevents.NewHttp("", nil, nil); err != nil {
		t.Fatal(err)
	} else {
		client.StopReceiver()
	}
}
