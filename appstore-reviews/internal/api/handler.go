package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"appstore-reviews/internal/review"
	"appstore-reviews/internal/store"
)

type Handler struct {
	store  *store.FileStore
	window time.Duration
}

func New(st *store.FileStore, window time.Duration) *Handler {
	return &Handler{store: st, window: window}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/reviews/", h.handleReviews)
	return withCORS(mux)
}

func (h *Handler) handleReviews(w http.ResponseWriter, r *http.Request) {
	appID := strings.TrimPrefix(r.URL.Path, "/reviews/")
	if appID == "" || strings.Contains(appID, "/") {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	
	all := h.store.Get(appID)
	cutoff := time.Now().Add(-h.window)

	recent := make([]review.Review, 0, len(all))
	for _, rev := range all {
		if rev.SubmittedAt.After(cutoff) {
			recent = append(recent, rev)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recent); err != nil {
		log.Printf("backend: encode failed for app %s: %v", appID, err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}