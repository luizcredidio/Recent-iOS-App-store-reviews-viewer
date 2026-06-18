package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"appstore-reviews/internal/api"
	"appstore-reviews/internal/review"
	"appstore-reviews/internal/store"
)

func newStore(t *testing.T) *store.FileStore {
	t.Helper()
	st, err := store.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestHandleReviews_WindowFilter(t *testing.T) {
	st := newStore(t)
	now := time.Now()
	_, _ = st.Save("app1", []review.Review{
		{ID: "recent", Author: "a", Score: 5, Content: "c", SubmittedAt: now.Add(-time.Hour)},
		{ID: "old", Author: "b", Score: 3, Content: "d", SubmittedAt: now.Add(-49 * time.Hour)},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reviews/app1", nil)
	api.New(st, 48*time.Hour).Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	var got []review.Review
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "recent" {
		t.Fatalf("want [recent], got %v", got)
	}
}

func TestHandleReviews_MissingID(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reviews/", nil)
	api.New(newStore(t), 48*time.Hour).Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}

func TestHandleReviews_CORSHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reviews/app1", nil)
	api.New(newStore(t), 48*time.Hour).Routes().ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("want CORS header *, got %q", got)
	}
}
