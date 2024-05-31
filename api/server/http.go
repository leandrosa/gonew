package server

import (
	"fmt"
	"log"
	"net/http"
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
	fmt.Println("Starting server")
	address := fmt.Sprintf("Addr: %s", server.Addr)
	fmt.Println(address)

	log.Println("main: running simple server on addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: couldn't start simple server: %v\n", err)
	}
}
