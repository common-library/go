package cloudevents_test

import (
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
	"time"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	"github.com/common-library/go/event/cloudevents"
)

const eventID = "id-01"
const eventType = "type 01"
const eventSource = "source/01"
const eventSubject = "subject 01"
const eventDataContentType = cloudeventssdk.ApplicationJSON
const eventExtensionsName01 = "name01"
const eventExtensionsValue01 = "1"
const eventExtensionsName02 = "name02"
const eventExtensionsValue02 = "value02"
const eventData = `{"id":1,"message":"Hello, World!"}`

func getEvent(t *testing.T) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetID(eventID)
	event.SetTime(time.Now())
	event.SetType(eventType)
	event.SetSource(eventSource)
	event.SetSubject(eventSubject)
	event.SetExtension(eventExtensionsName01, eventExtensionsValue01)
	event.SetExtension(eventExtensionsName02, eventExtensionsValue02)
	if err := event.SetData(cloudeventssdk.ApplicationJSON, map[string]any{
		"id":      1,
		"message": "Hello, World!",
	}); err != nil {
		t.Fatal(err)
	}

	return event
}

func consistencyEvent(t *testing.T, event *cloudevents.Event) {
	if event.ID() != eventID {
		t.Error("invalid -", event.ID())
	} else if event.Type() != eventType {
		t.Error("invalid -", event.Type())
	} else if event.Source() != eventSource {
		t.Error("invalid -", event.Source())
	} else if event.Subject() != eventSubject {
		t.Error("invalid -", event.Subject())
	} else if event.DataContentType() != eventDataContentType {
		t.Error("invalid -", event.DataContentType())
	} else if event.Extensions()[eventExtensionsName01] == nil ||
		event.Extensions()[eventExtensionsName01].(string) != eventExtensionsValue01 {
		t.Error("invalid -", event.Extensions()[eventExtensionsName01])
	} else if event.Extensions()[eventExtensionsName02] == nil ||
		event.Extensions()[eventExtensionsName02].(string) != eventExtensionsValue02 {
		t.Error("invalid -", event.Extensions()[eventExtensionsName02])
	} else if string(event.Data()) != eventData {
		t.Error("invalid -", string(event.Data()))
	}
}

func startServer(t *testing.T) (string, *cloudevents.Server) {
	handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
		consistencyEvent(t, &event)

		responseEvent := event.Clone()
		return &responseEvent, cloudevents.NewHTTPResult(http.StatusOK, "")
	}

	address := ":" + strconv.Itoa(10000+rand.IntN(1000))
	listenAndServeFailureFunc := func(err error) { t.Fatal(err) }

	server := cloudevents.Server{}
	if err := server.Start(address, handler, listenAndServeFailureFunc); err != nil {
		t.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	return address, &server
}

func stopServer(t *testing.T, server *cloudevents.Server) {
	if err := server.Stop(10); err != nil {
		t.Fatal(err)
	}
}
