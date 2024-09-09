package main

import (
	"context"
	"log"
	"net/http"
	"oauthive/api"
	"oauthive/db"
	"oauthive/domain/authenticator"
	"os"
	"os/signal"
	"time"
)

func startServer(handler http.Handler, addr string) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Printf("Server is starting on %s...\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	return server
}

func gracefulShutdown(server *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	log.Println("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	log.Println("Server shut down gracefully.")
}

func main() {
	database := db.Init("database.db")
	authenticator := authenticator.NewAuthenticator()

	mux := api.NewMux(database, authenticator)
	server := startServer(mux, ":3333")

	gracefulShutdown(server, 5*time.Second)
}
