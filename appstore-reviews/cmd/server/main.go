package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"appstore-reviews/internal/poller"
	"appstore-reviews/internal/store"
	"appstore-reviews/internal/api"
)

func main() {
	appIDs := []string{"595068606"}
	interval := 30 * time.Second

	st, err := store.New()
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	h := api.New(st, 720 * time.Hour)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h.Routes()))

	p := poller.New(st, appIDs, interval)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("starting poller: apps=%v interval=%s", appIDs, interval)
	p.Run(ctx)

	log.Println("poller stopped")
}