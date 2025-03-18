package server

import (
	"context"
	"errors"
	"fmt"
	"go-intro/internal/handler"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// Server is an object representing a server instance, required dependencies can be added here. s
type Server struct {
	Port    int
	Server  *http.Server
	Handler *handler.Handler
}

// New creates a new instance of a Mux.
// handlers could also be a client/consumer interface pattern.
func New(port int, handler *handler.Handler) (*Server, error) {
	mux := &Server{
		Port: port,
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		Handler: handler,
	}

	// Create a new router
	router := chi.NewRouter()

	// Add middlewares
	router.Use(requestLogging)   // This logs the request's context data.
	router.Use(CORSMiddleware()) // This adds cors

	// This ensures that if there's fatal error during running of an endpoint the server will recover instead of shut down.
	router.Use(middleware.Recoverer)

	// timeout on the request using context.Done()
	router.Use(middleware.Timeout(25 * time.Second))

	router.Post("/add", nil)
	router.Get("/todos", mux.Handler.GetTodos)
	// ping
	router.Get("/ping", func(writer http.ResponseWriter, request *http.Request) {
		if _, err := writer.Write([]byte("devaaa")); err != nil {
			logrus.WithField("error", err).Error("failed to write ping response")
		}
	})

	// Assign the router to the mux.
	mux.Server.Handler = router
	return mux, nil
}

// ListenAndServe create a new server and runs it in a separate go routine.
// On failure, it logs fatally and shuts down the service.
func (m *Server) ListenAndServe(logger *logrus.Entry) {
	go func() {
		if err := m.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Infof("server listening on port: %d", m.Port)
}

// Shutdown gracefully shuts the server down.
func (m *Server) Shutdown() error {
	return m.Server.Shutdown(context.Background())
}
