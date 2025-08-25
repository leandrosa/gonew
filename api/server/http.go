package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServerConfig struct {
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func NewServerConfig(port string) ServerConfig {
	return ServerConfig{
		port:         port,
		readTimeout:  15 * time.Second,
		writeTimeout: 15 * time.Second,
		idleTimeout:  15 * time.Second,
	}
}

func (s ServerConfig) GetAddress() string {
	return fmt.Sprintf(":%s", s.port)
}

func NewServer(config ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         config.GetAddress(),
		ReadTimeout:  config.readTimeout,
		WriteTimeout: config.writeTimeout,
		IdleTimeout:  config.idleTimeout,
	}
}

func AddHandler(route string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(route, handler)
}

func AddHandlerMessage(route string, message string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, message)
	})
}

func RunServer(server *http.Server) {
	startServer(server)
}

// startServer implementing graceful shutdown
// reference: https://dev.to/yanev/a-deep-dive-into-graceful-shutdown-in-go-484a
func startServer(server *http.Server) {
	serverError := make(chan error, 1)

	go func() {
		log.Printf("Server is running on %s", getFullAddress(server))
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return
	}

	log.Println("Server exited properly")
}

func getFullAddress(server *http.Server) string {
	addr := server.Addr

	// If addr starts with ":", assume localhost
	if addr[0] == ':' {
		return "http://localhost" + addr
	}
	return "http://" + addr
}
