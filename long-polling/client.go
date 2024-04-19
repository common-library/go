// Package long_polling provides long polling client and server implementations.
package long_polling

import (
	"fmt"
	net_http "net/http"
	net_url "net/url"
	"time"

	"github.com/common-library/go/http"
	"github.com/common-library/go/json"
	"github.com/google/go-querystring/query"
)

// SubscriptionRequest is subscription request information.
type SubscriptionRequest struct {
	Category       string `url:"category"`
	TimeoutSeconds int    `url:"timeout"`
	SinceTime      int64  `url:"since_time,omitempty"`
	LastID         string `url:"last_id,omitempty"`
}

// SubscriptionResponse is subscription response information.
type SubscriptionResponse struct {
	Header     net_http.Header
	StatusCode int
	Events     []struct {
		Timestamp int64  `json:"timestamp"`
		Category  string `json:"category"`
		ID        string `json:"id"`
		Data      string `json:"data"`
	} `json:"events"`
}

// PublishRequest is publish request information.
type PublishRequest struct {
	Category string `json:"category"`
	Data     string `json:"data"`
}

// Subscription is subscribes to event.
//
// ex) response, err := long_polling.Subscription("http://127.0.0.1:10000/subscription", nil, request, "", "", nil)
func Subscription(url string, header map[string][]string, request SubscriptionRequest, username, password string, transport *net_http.Transport) (SubscriptionResponse, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	values, err := query.Values(request)
	if err != nil {
		return SubscriptionResponse{}, err
	}
	u.RawQuery = values.Encode()

	response, err := http.Request(fmt.Sprintf("%v", u), net_http.MethodGet, header, "", time.Duration(request.TimeoutSeconds)*time.Second, username, password, transport)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	subscriptionResponse := SubscriptionResponse{Header: response.Header, StatusCode: response.StatusCode}

	if response.StatusCode == net_http.StatusOK {
		if result, err := json.ConvertFromString[SubscriptionResponse](response.Body); err != nil {
			return SubscriptionResponse{}, err
		} else {
			subscriptionResponse.Events = result.Events
		}
	}

	return subscriptionResponse, nil
}

// Publish is publish an event.
//
// ex) response, err := long_polling.Publish("http://127.0.0.1:10000/publish", 10, nil, request, "", "", nil)
func Publish(url string, timeout time.Duration, header map[string][]string, publishRequest PublishRequest, username, password string, transport *net_http.Transport) (http.Response, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return http.Response{}, err
	}

	body, err := json.ToString(publishRequest)
	if err != nil {
		return http.Response{}, err
	}

	return http.Request(fmt.Sprintf("%v", u), net_http.MethodPost, header, body, timeout, username, password, transport)
}
