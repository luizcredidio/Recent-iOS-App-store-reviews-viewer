package store

import (
	"testing"
	"time"

	"appstore-reviews/internal/review"
)

func rev(id string, at time.Time) review.Review {
	return review.Review{ID: id, Author: "a", Score: 5, Content: "c", SubmittedAt: at}
}

func TestSaveDedup(t *testing.T) {
	st, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	r := rev("r1", time.Now())

	if n, err := st.Save("app1", []review.Review{r}); err != nil || n != 1 {
		t.Fatalf("first save: want n=1 err=nil, got n=%d err=%v", n, err)
	}
	if n, err := st.Save("app1", []review.Review{r}); err != nil || n != 0 {
		t.Fatalf("duplicate save: want n=0 err=nil, got n=%d err=%v", n, err)
	}
	if got := st.Get("app1"); len(got) != 1 {
		t.Fatalf("want 1 review in store, got %d", len(got))
	}
}

func TestGetNewestFirst(t *testing.T) {
	st, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	_, _ = st.Save("app1", []review.Review{
		rev("old", now.Add(-2*time.Hour)),
		rev("new", now),
		rev("mid", now.Add(-time.Hour)),
	})

	got := st.Get("app1")
	want := []string{"new", "mid", "old"}
	for i, r := range got {
		if r.ID != want[i] {
			t.Fatalf("position %d: want %q got %q", i, want[i], r.ID)
		}
	}
}

func TestPersistAndReload(t *testing.T) {
	dir := t.TempDir()

	st1, err := NewFileStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st1.Save("app1", []review.Review{rev("r1", time.Now())}); err != nil {
		t.Fatal(err)
	}

	st2, err := NewFileStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	got := st2.Get("app1")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("after reload: want [r1], got %v", got)
	}
}
