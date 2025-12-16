// Package cloudevents provides cloudevents client and server implementations.
package cloudevents

import (
	"errors"

	cloudeventssdk "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cloudeventssdk_http "github.com/cloudevents/sdk-go/v2/protocol/http"
)

// NewResult creates and returns a result.
//
// ex) result := cloudevents.NewResult("ok")
func NewResult(format string, arguments ...any) Result {
	return Result{result: cloudeventssdk.NewResult(format, arguments)}
}

// NewResult creates and returns a result.
//
// ex) result := cloudevents.NewHTTPResult(http.StatusOK, "")
func NewHTTPResult(statusCode int, format string, arguments ...any) Result {
	return Result{result: cloudeventssdk.NewHTTPResult(statusCode, format, arguments)}
}

// Result is the result of event delivery.
type Result struct {
	result protocol.Result
}

// IsACK returns whether the recipient acknowledged the event.
//
// ex) isACK := result.IsACK()
func (r *Result) IsACK() bool {
	return cloudeventssdk.IsACK(r.result)
}

// IsNACK returns whether the recipient did not acknowledge the event.
//
// ex) isNACK := result.IsNACK()
func (r *Result) IsNACK() bool {
	return cloudeventssdk.IsNACK(r.result)
}

// IsUndelivered returns whether it was delivered or not.
//
// ex) isUndelivered := result.IsUndelivered()
func (r *Result) IsUndelivered() bool {
	return cloudeventssdk.IsUndelivered(r.result)
}

// GetHttpStatusCode returns the status code if the result is http.
//
// ex) statusCode, err := result.GetHttpStatusCode()
func (r *Result) GetHttpStatusCode() (int, error) {
	httpResult := new(cloudeventssdk_http.Result)

	if !cloudeventssdk.ResultAs(r.result, &httpResult) {
		return -1, errors.New("match failed")
	} else {
		return httpResult.StatusCode, nil
	}
}

// Error returns the error string.
//
// ex) errString := result.Error()
func (r *Result) Error() string {
	return r.result.Error()
}
