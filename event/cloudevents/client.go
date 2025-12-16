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
func (c *client) Send(event Event) Result {
	return Result{result: c.clientOfSdk.Send(c.getContext(), event)}
}

// Send transmits an event and returns a response event.
//
// ex) responseEvent, result := client.Request(event)
func (c *client) Request(event Event) (*Event, Result) {
	responseEvent, result := c.clientOfSdk.Request(c.getContext(), event)

	return responseEvent, Result{result: result}
}

// StartReceiver receives events until StopReceiver is called.
//
// ex)
//
//	httpOption := []cloudeventssdk_http.Option{cloudeventssdk_http.WithPort(port)}
//	receiveclient, err := cloudevents.NewHttp("", httpOption, nil)
//	receiveclient.StartReceiver(handler, failureFunc)
func (c *client) StartReceiver(handler func(context.Context, Event), failureFunc func(error)) {
	c.wgForReceiver.Add(1)
	go func() {
		defer c.wgForReceiver.Done()

		ctx, cancel := context.WithCancel(c.getContext())
		c.cancelFuncForReceiver = cancel

		if err := c.clientOfSdk.StartReceiver(ctx, handler); err != nil {
			failureFunc(err)
		}
	}()
}

// StopReceiver stops receiving events by StartReceiver.
//
// ex)client.StopReceiver()
func (c *client) StopReceiver() {
	if c.cancelFuncForReceiver != nil {
		c.cancelFuncForReceiver()
	}
	c.wgForReceiver.Wait()
}

func (c *client) getContext() context.Context {
	switch c.clientType {
	case clientTypeHttp:
		return cloudeventssdk.ContextWithTarget(context.Background(), c.address)
	default:
		return nil
	}
}
