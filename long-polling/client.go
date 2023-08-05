package long_polling

import (
	"fmt"
	net_http "net/http"
	net_url "net/url"

	"github.com/google/go-querystring/query"
	"github.com/heaven-chp/common-library-go/http"
	"github.com/heaven-chp/common-library-go/json"
)

type SubscriptionRequest struct {
	Category  string `url:"category"`
	Timeout   int    `url:"timeout"`
	SinceTime int64  `url:"since_time,omitempty"`
	LastID    string `url:"last_id,omitempty"`
}

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

type PublishRequest struct {
	Category string `json:"category"`
	Data     string `json:"data"`
}

func Subscription(url string, header map[string][]string, request SubscriptionRequest, username, password string) (SubscriptionResponse, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	values, err := query.Values(request)
	if err != nil {
		return SubscriptionResponse{}, err
	}
	u.RawQuery = values.Encode()

	response, err := http.Request(fmt.Sprintf("%v", u), net_http.MethodGet, header, "", request.Timeout, username, password)
	if err != nil {
		return SubscriptionResponse{}, err
	}

	subscriptionResponse := SubscriptionResponse{Header: response.Header, StatusCode: response.StatusCode}

	if response.StatusCode == net_http.StatusOK {
		err = json.ToStructFromString(response.Body, &subscriptionResponse)
		if err != nil {
			return SubscriptionResponse{}, err
		}
	}

	return subscriptionResponse, nil
}

func Publish(url string, timeout int, header map[string][]string, publishRequest PublishRequest, username, password string) (http.Response, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return http.Response{}, err
	}

	body, err := json.ToString(publishRequest)
	if err != nil {
		return http.Response{}, err
	}

	return http.Request(fmt.Sprintf("%v", u), net_http.MethodPost, header, body, timeout, username, password)
}
