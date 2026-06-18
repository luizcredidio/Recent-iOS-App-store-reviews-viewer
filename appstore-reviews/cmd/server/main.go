package main

import (
	"log"
	"appstore-reviews/internal/appstore"
	"appstore-reviews/internal/store"
)

func main() {
	const appID = "595068606"

	reviews, err := appstore.FetchReviews(appID)
	if err != nil {
		log.Fatalf("failed to fetch reviews: %v", err)
	}

	log.Printf("fetched %d reviews for app %s", len(reviews), appID)
	for _, r := range reviews {
		log.Printf("%s | %d★ | %s | %s",
			r.SubmittedAt.Format("2006-01-02 15:04"),
			r.Score,
			r.Author,
			r.Title,
		)
	}

	st, err := store.New()
	if err != nil {
		log.Fatalf("store init: %v", err)
	}

	added, err := st.Save(appID, reviews)
	if err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Printf("saved %d", added)
}