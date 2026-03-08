package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server is running")
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//run server in its own goroutine outside of main function
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	sig := <-shutdownCh
	log.Printf("Shutdown signal received: %v", sig)
	//process of closing the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Server shutdown failed: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}

	signal.Stop(shutdownCh)
	close(shutdownCh)

	log.Println("Application exited cleanly")
}
