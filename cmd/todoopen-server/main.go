package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justEstif/todo-open/internal/app"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("TODOOPEN_SERVER_ADDR"); v != "" {
		addr = v
	}

	srv := app.NewServer(addr)
	log.Printf("todo.open server starting on %s", addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}
	log.Printf("todo.open server stopped")
}
