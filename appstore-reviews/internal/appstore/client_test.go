package appstore

import (
	"testing"
	"time"
)

func TestToReview_Valid(t *testing.T) {
	e := feedEntry{
		ID:      labeled{"id-1"},
		Title:   labeled{"Great app"},
		Content: labeled{"Really good"},
		Updated: labeled{"2024-01-15T10:00:00Z"},
		Rating:  labeled{"5"},
		Version: labeled{"2.0"},
		Author:  struct{ Name labeled `json:"name"` }{Name: labeled{"Alice"}},
	}
	r, err := toReview(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID != "id-1" || r.Score != 5 || r.Author != "Alice" || r.AppVersion != "2.0" {
		t.Fatalf("unexpected review: %+v", r)
	}
	if r.SubmittedAt.IsZero() {
		t.Fatal("SubmittedAt should not be zero")
	}
}

func TestToReview_BadRating(t *testing.T) {
	e := feedEntry{
		Updated: labeled{time.Now().Format(time.RFC3339)},
		Rating:  labeled{"notanumber"},
	}
	if _, err := toReview(e); err == nil {
		t.Fatal("want error for non-numeric rating, got nil")
	}
}

func TestToReview_BadTimestamp(t *testing.T) {
	e := feedEntry{
		Updated: labeled{"not-a-date"},
		Rating:  labeled{"4"},
	}
	if _, err := toReview(e); err == nil {
		t.Fatal("want error for bad timestamp, got nil")
	}
}
