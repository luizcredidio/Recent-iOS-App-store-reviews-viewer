package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"net/http"
	"errors"

	"appstore-reviews/internal/poller"
	"appstore-reviews/internal/store"
	"appstore-reviews/internal/api"
)

func main() {
	appIDs := []string{"595068606"}
	interval := 30 * time.Second
	reviewWindow := 720 * time.Hour

	st, err := store.New()
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	p := poller.New(st, appIDs, interval)
	go p.Run(ctx)


	h := api.New(st, reviewWindow)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h.Routes(),
	}

	go func() {
	log.Printf("listening on %s", ":8080")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
	}()

	<-ctx.Done()

	log.Println("shut down")
}