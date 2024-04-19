// Package cloudevents provides cloudevents client and server implementations.
package cloudevents

import (
	"context"
	"sync"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	cloudeventssdk_client "github.com/cloudevents/sdk-go/v2/client"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
)

type clientType int

const (
	clientTypeHttp clientType = iota + 1
)

// NewHttp creates and returns a http client.
//
// ex) client, err := cloudevents.NewHttp(address, nil, nil)
func NewHttp(address string, httpOption []http.Option, clientOption []cloudeventssdk_client.Option) (*client, error) {
	if protocol, err := cloudeventssdk.NewHTTP(httpOption...); err != nil {
		return nil, err
	} else if clientOfSdk, err := cloudeventssdk.NewClient(protocol, clientOption...); err != nil {
		return nil, err
	} else {
		return &client{clientType: clientTypeHttp, clientOfSdk: clientOfSdk, address: address}, nil
	}
}

// client is a struct that provides client related methods.
type client struct {
	clientType clientType

	clientOfSdk cloudeventssdk_client.Client
	address     string

	wgForReceiver         sync.WaitGroup
	cancelFuncForReceiver context.CancelFunc
}

// Send transmits an event.
//
// ex) result := client.Send(event)
func (this *client) Send(event Event) Result {
	return Result{result: this.clientOfSdk.Send(this.getContext(), event)}
}

// Send transmits an event and returns a response event.
//
// ex) responseEvent, result := client.Request(event)
func (this *client) Request(event Event) (*Event, Result) {
	responseEvent, result := this.clientOfSdk.Request(this.getContext(), event)

	return responseEvent, Result{result: result}
}

// StartReceiver receives events until StopReceiver is called.
//
// ex)
//
//	httpOption := []cloudeventssdk_http.Option{cloudeventssdk_http.WithPort(port)}
//	receiveclient, err := cloudevents.NewHttp("", httpOption, nil)
//	receiveclient.StartReceiver(handler, failureFunc)
func (this *client) StartReceiver(handler func(context.Context, Event), failureFunc func(error)) {
	this.wgForReceiver.Add(1)
	go func() {
		defer this.wgForReceiver.Done()

		ctx, cancel := context.WithCancel(this.getContext())
		this.cancelFuncForReceiver = cancel

		if err := this.clientOfSdk.StartReceiver(ctx, handler); err != nil {
			failureFunc(err)
		}
	}()
}

// StopReceiver stops receiving events by StartReceiver.
//
// ex)client.StopReceiver()
func (this *client) StopReceiver() {
	if this.cancelFuncForReceiver != nil {
		this.cancelFuncForReceiver()
	}
	this.wgForReceiver.Wait()
}

func (this *client) getContext() context.Context {
	switch this.clientType {
	case clientTypeHttp:
		return cloudeventssdk.ContextWithTarget(context.Background(), this.address)
	default:
		return nil
	}
}
