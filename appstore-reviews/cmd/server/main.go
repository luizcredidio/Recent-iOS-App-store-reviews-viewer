package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"net/http"
	"errors"

	"appstore-reviews/internal/config"
	"appstore-reviews/internal/poller"
	"appstore-reviews/internal/store"
	"appstore-reviews/internal/api"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	st, err := store.New()
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	p := poller.New(st, cfg.AppIDs, cfg.PollInterval())
	go p.Run(ctx)

	h := api.New(st, cfg.ReviewWindow())
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: h.Routes(),
	}

	go func() {
		log.Printf("listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("shut down")
}