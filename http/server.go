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
// ex)
//
//	server.RegisterHandler("/xxx", handler) //all method
//	server.RegisterHandler("/xxx", handler, http.MethodGet)
//	server.RegisterHandler("/xxx", handler, http.MethodGet, http.MethodPost)
func (s *Server) RegisterHandler(path string, handler http.Handler, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().Handle(path, handler)
	} else {
		s.getRouter().Handle(path, handler).Methods(methods...)
	}
}

// RegisterHandlerFunc is add handler function.
//
// ex)
//
//	server.RegisterHandlerFunc("/xxx", handlerFunc) //all method
//	server.RegisterHandlerFunc("/xxx", handlerFunc, http.MethodGet)
//	server.RegisterHandlerFunc("/xxx", handlerFunc, http.MethodGet, http.MethodPost)
func (s *Server) RegisterHandlerFunc(path string, handlerFunc http.HandlerFunc, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().HandleFunc(path, handlerFunc)
	} else {
		s.getRouter().HandleFunc(path, handlerFunc).Methods(methods...)
	}
}

// RegisterPathPrefixHandler is add path prefix handler.
//
// ex)
//
//	server.RegisterPathPrefixHandler("/xxx", handler) //all method
//	server.RegisterPathPrefixHandler("/xxx", handler, http.MethodGet)
//	server.RegisterPathPrefixHandler("/xxx", handler, http.MethodGet, http.MethodPost)
func (s *Server) RegisterPathPrefixHandler(prefix string, handler http.Handler, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().PathPrefix(prefix).Handler(handler)
	} else {
		s.getRouter().PathPrefix(prefix).Handler(handler).Methods(methods...)
	}
}

// RegisterPathPrefixHandlerFunc is add path prefix handler function.
//
// ex)
//
//	server.RegisterPathPrefixHandlerFunc("/xxx", handlerFunc) //all method
//	server.RegisterPathPrefixHandlerFunc("/xxx", handlerFunc, http.MethodGet)
//	server.RegisterPathPrefixHandlerFunc("/xxx", handlerFunc, http.MethodGet, http.MethodPost)
func (s *Server) RegisterPathPrefixHandlerFunc(prefix string, handlerFunc http.HandlerFunc, methods ...string) {
	if len(methods) == 0 {
		s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc)
	} else {
		s.getRouter().PathPrefix(prefix).HandlerFunc(handlerFunc).Methods(methods...)
	}
}

// Start is start the server.
//
// ex) err := server.Start(":10000")
func (s *Server) Start(address string, listenAndServeFailureFunc func(err error), middlewareFunc ...mux.MiddlewareFunc) error {
	if middlewareFunc != nil {
		s.getRouter().Use(middlewareFunc...)
	} else {
		s.getRouter().Use(func(nextHandler http.Handler) http.Handler {
			return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
				nextHandler.ServeHTTP(responseWriter, request)
			})
		})
	}

	s.server = &http.Server{
		Addr:    address,
		Handler: s.getRouter()}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed && listenAndServeFailureFunc != nil {
			listenAndServeFailureFunc(err)
		}
	}()

	return nil
}

// Stop is stop the server.
//
// ex) err := server.Stop(10)
func (s *Server) Stop(shutdownTimeout time.Duration) error {
	s.SetRouter(nil)

	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// SetRouter is set the router.
//
// ex) server.SetRouter(router)
func (s *Server) SetRouter(router *mux.Router) {
	s.router = router
}

func (s *Server) getRouter() *mux.Router {
	if s.router == nil {
		s.router = mux.NewRouter()
	}

	return s.router
}
