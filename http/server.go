// Package http provides http client and server implementations.
package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is a struct that provides server related methods.
type Server struct {
	server *http.Server
	router *mux.Router
}

// AddHandler is add handler.
//
// ex) server.AddHandler("/v1/test", http.MethodPost, handler)
func (this *Server) AddHandler(path, method string, handler func(http.ResponseWriter, *http.Request)) {
	if this.router == nil {
		this.router = mux.NewRouter()
	}

	this.router.HandleFunc(path, handler).Methods(method)
}

// AddPathPrefixHandler is add path prefix handler.
//
// ex) server.AddPathPrefixHandler("/swagger/", httpSwagger.WrapHandler)
func (this *Server) AddPathPrefixHandler(prefix string, handler http.Handler) {
	if this.router == nil {
		this.router = mux.NewRouter()
	}

	this.router.PathPrefix(prefix).Handler(handler)
}

// Start is start the server.
//
// ex) err := server.Start("127.0.0.1")
func (this *Server) Start(address string, listenAndServeFailureFunc func(err error), middlewareFunc ...mux.MiddlewareFunc) error {
	if this.router == nil {
		this.router = mux.NewRouter()
	}

	if middlewareFunc != nil {
		this.router.Use(middlewareFunc...)
	} else {
		this.router.Use(func(nextHandler http.Handler) http.Handler {
			return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
				nextHandler.ServeHTTP(responseWriter, request)
			})
		})
	}

	this.server = &http.Server{
		Addr:    address,
		Handler: this.router}

	go func() {
		err := this.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed && listenAndServeFailureFunc != nil {
			listenAndServeFailureFunc(err)
		}
	}()

	return nil
}

// Stop is stop the server.
//
// ex) err := server.Stop(10)
func (this *Server) Stop(shutdownTimeout uint64) error {
	this.router = nil
	if this.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	return this.server.Shutdown(ctx)
}

// SetRouter is set the router.
//
// ex) server.SetRouter(router)
func (this *Server) SetRouter(router *mux.Router) {
	this.router = router
}
