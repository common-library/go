// Package cloudevents provides CloudEvents client and server implementations.
//
// This package wraps the official CloudEvents SDK for Go, offering simplified interfaces
// for sending, receiving, and processing CloudEvents over HTTP and other protocols.
//
// Features:
//   - HTTP client for sending and receiving CloudEvents
//   - Request-response pattern support
//   - Asynchronous event receiver with lifecycle management
//   - Result types for event delivery status
//   - CloudEvents v1.0 specification compliance
//
// Example:
//
//	client, _ := cloudevents.NewHttp("http://localhost:8080", nil, nil)
//	event := cloudevents.NewEvent()
//	event.SetType("com.example.event")
//	event.SetSource("example/source")
//	result := client.Send(event)
package cloudevents

import (
	v2 "github.com/cloudevents/sdk-go/v2"
)

// NewEvent creates a new CloudEvent with default attributes.
//
// Returns:
//   - Event: A new CloudEvent instance with generated ID and timestamp
//
// The returned event must have at minimum Type and Source set before being valid
// according to the CloudEvents specification.
//
// Example:
//
//	event := cloudevents.NewEvent()
//	event.SetType("com.example.user.created")
//	event.SetSource("example/users")
//	event.SetData("application/json", userData)
var NewEvent = v2.NewEvent

// Event represents a CloudEvent conforming to the CloudEvents v1.0 specification.
//
// Events have required attributes (Type, Source, ID, SpecVersion) and optional
// attributes (DataContentType, DataSchema, Subject, Time). Events can carry data
// payloads of any type.
//
// Example:
//
//	var event cloudevents.Event
//	event.SetID("abc-123")
//	event.SetType("com.example.object.action")
//	event.SetSource("example/source")
//	event.SetData("application/json", myData)
//
//	// Access attributes
//	eventType := event.Type()
//	eventSource := event.Source()
//	data := event.Data()
type Event = v2.Event
