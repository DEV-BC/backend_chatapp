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

	"github.com/DEV-BC/backend_chatapp/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//run server in its own goroutine outside of main function
	go func() {
		log.Printf("Server is running on https://%s\n", cfg.HTTPServer.Address)
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
