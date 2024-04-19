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

// RegisterHandler is add handler.
//
// ex) server.RegisterHandler("/xxx", http.MethodPost, handler)
func (this *Server) RegisterHandler(path, method string, handler http.Handler) {
	this.getRouter().Handle(path, handler).Methods(method)
}

// RegisterHandlerFunc is add handler function.
//
// ex) server.RegisterHandlerFunc("/xxx", http.MethodPost, handlerFunc)
func (this *Server) RegisterHandlerFunc(path, method string, handlerFunc http.HandlerFunc) {
	this.getRouter().HandleFunc(path, handlerFunc).Methods(method)
}

// RegisterPathPrefixHandler is add path prefix handler.
//
// ex) server.RegisterPathPrefixHandler("xxx", handler)
func (this *Server) RegisterPathPrefixHandler(prefix string, handler http.Handler) {
	this.getRouter().PathPrefix(prefix).Handler(handler)
}

// RegisterPathPrefixHandlerFunc is add path prefix handler function.
//
// ex) server.RegisterPathPrefixHandlerFunc("xxx", Func)
func (this *Server) RegisterPathPrefixHandlerFunc(prefix string, handlerFunc http.HandlerFunc) {
	this.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc)
}

// Start is start the server.
//
// ex) err := server.Start(":10000")
func (this *Server) Start(address string, listenAndServeFailureFunc func(err error), middlewareFunc ...mux.MiddlewareFunc) error {
	if middlewareFunc != nil {
		this.getRouter().Use(middlewareFunc...)
	} else {
		this.getRouter().Use(func(nextHandler http.Handler) http.Handler {
			return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
				nextHandler.ServeHTTP(responseWriter, request)
			})
		})
	}

	this.server = &http.Server{
		Addr:    address,
		Handler: this.getRouter()}

	go func() {
		if err := this.server.ListenAndServe(); err != nil && err != http.ErrServerClosed && listenAndServeFailureFunc != nil {
			listenAndServeFailureFunc(err)
		}
	}()

	return nil
}

// Stop is stop the server.
//
// ex) err := server.Stop(10)
func (this *Server) Stop(shutdownTimeout time.Duration) error {
	this.SetRouter(nil)

	if this.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	return this.server.Shutdown(ctx)
}

// SetRouter is set the router.
//
// ex) server.SetRouter(router)
func (this *Server) SetRouter(router *mux.Router) {
	this.router = router
}

func (this *Server) getRouter() *mux.Router {
	if this.router == nil {
		this.router = mux.NewRouter()
	}

	return this.router
}
