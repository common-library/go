// Package cloudevents provides cloudevents client and server implementations.
package cloudevents

import (
	v2 "github.com/cloudevents/sdk-go/v2"
)

var NewEvent = v2.NewEvent

type Event = v2.Event
