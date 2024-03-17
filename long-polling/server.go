// Package long_polling provides long polling client and server implementations.
package long_polling

import (
	net_http "net/http"

	"github.com/gorilla/mux"
	"github.com/heaven-chp/common-library-go/http"
	"github.com/jcuga/golongpoll"
)

// ServerInfo is server information.
type ServerInfo struct {
	Address string
	Timeout int

	SubscriptionURI                string
	HandlerToRunBeforeSubscription func(w net_http.ResponseWriter, r *net_http.Request) bool

	PublishURI                string
	HandlerToRunBeforePublish func(w net_http.ResponseWriter, r *net_http.Request) bool
}

// FilePersistorInfo is file persistor information.
type FilePersistorInfo struct {
	Use                     bool
	FileName                string
	WriteBufferSize         int
	WriteFlushPeriodSeconds int
}

// Server is a struct that provides server related methods.
type Server struct {
	server http.Server

	longpollManager *golongpoll.LongpollManager
}

// Start is start the server.
//
// ex) err := server.Start(ServerInfo{...}, FilePersistorInfo{...}, nil)
func (this *Server) Start(serverInfo ServerInfo, filePersistorInfo FilePersistorInfo, listenAndServeFailureFunc func(err error)) error {
	option := golongpoll.Options{
		LoggingEnabled:            false,
		MaxLongpollTimeoutSeconds: serverInfo.Timeout,
		//MaxEventBufferSize: 250,
		//EventTimeToLiveSeconds:,
		//DeleteEventAfterFirstRetrieval:,
	}

	if filePersistorInfo.Use {
		filePersistor, err := golongpoll.NewFilePersistor(filePersistorInfo.FileName, filePersistorInfo.WriteBufferSize, filePersistorInfo.WriteFlushPeriodSeconds)
		if err != nil {
			return err
		}

		option.AddOn = filePersistor
	}

	longpollManager, err := golongpoll.StartLongpoll(option)
	if err != nil {
		return err
	}
	this.longpollManager = longpollManager

	router := mux.NewRouter()

	router.Use(func(nextHandler net_http.Handler) net_http.Handler {
		return net_http.HandlerFunc(func(responseWriter net_http.ResponseWriter, request *net_http.Request) {
			nextHandler.ServeHTTP(responseWriter, request)
		})
	})

	subscriptionHandler := func() func(net_http.ResponseWriter, *net_http.Request) {
		return func(w net_http.ResponseWriter, r *net_http.Request) {
			if serverInfo.HandlerToRunBeforeSubscription != nil &&
				serverInfo.HandlerToRunBeforeSubscription(w, r) == false {
				return
			}

			this.longpollManager.SubscriptionHandler(w, r)
		}
	}
	router.HandleFunc(serverInfo.SubscriptionURI, subscriptionHandler()).Methods(net_http.MethodGet)

	publishHandler := func() func(net_http.ResponseWriter, *net_http.Request) {
		return func(w net_http.ResponseWriter, r *net_http.Request) {
			if serverInfo.HandlerToRunBeforePublish != nil &&
				serverInfo.HandlerToRunBeforePublish(w, r) == false {
				return
			}

			this.longpollManager.PublishHandler(w, r)
		}
	}
	router.HandleFunc(serverInfo.PublishURI, publishHandler()).Methods(net_http.MethodPost)

	this.server.SetRouter(router)

	return this.server.Start(serverInfo.Address, listenAndServeFailureFunc)
}

// Stop is stop the server.
//
// ex) err := server.Stop(10)
func (this *Server) Stop(shutdownTimeout uint64) error {
	if this.longpollManager != nil {
		defer this.longpollManager.Shutdown()
	}

	err := this.server.Stop(shutdownTimeout)
	if err != nil {
		return err
	}

	return nil
}
