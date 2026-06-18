package main

import (
	"log"
	"appstore-reviews/internal/appstore"
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
}