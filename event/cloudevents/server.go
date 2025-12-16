// Package cloudevents provides cloudevents client and server implementations.
package cloudevents

import (
	"context"
	"time"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/common-library/go/http"
)

// Server is a struct that provides server related methods.
type Server struct {
	server http.Server
}

// Start is start the server.
//
// ex) err := server.Start(address, handler, listenAndServeFailureFunc)
func (s *Server) Start(address string, handler func(Event) (*Event, Result), listenAndServeFailureFunc func(error)) error {
	finalHandler := func(requestEvent Event) (*Event, protocol.Result) {
		responseEvent, result := handler(requestEvent)
		return responseEvent, result.result
	}

	if protocol, err := cloudeventssdk.NewHTTP(); err != nil {
		return err
	} else if eventReceiver, err := cloudeventssdk.NewHTTPReceiveHandler(context.Background(), protocol, finalHandler); err != nil {

		return err
	} else {
		s.server.RegisterPathPrefixHandler("/", eventReceiver)

		return s.server.Start(address, listenAndServeFailureFunc)
	}
}

// Stop is stop the server.
//
// ex) err := server.Stop(10)
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	return s.server.Stop(shutdownTimeout)
}
