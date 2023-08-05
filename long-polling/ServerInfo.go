package long_polling

import "net/http"

type ServerInfo struct {
	Address string
	Timeout int

	SubscriptionURI                string
	HandlerToRunBeforeSubscription func(w http.ResponseWriter, r *http.Request) bool

	PublishURI                string
	HandlerToRunBeforePublish func(w http.ResponseWriter, r *http.Request) bool
}
