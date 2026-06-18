package appstore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"appstore-reviews/internal/review"
)

type labeled struct {
	Label string `json:"label"`
}

type feedResponse struct {
	Feed struct {
		Entry []feedEntry `json:"entry"`
	} `json:"feed"`
}

type feedEntry struct {
	ID      labeled `json:"id"`
	Title   labeled `json:"title"`
	Content labeled `json:"content"`
	Updated labeled `json:"updated"`
	Rating  labeled `json:"im:rating"`
	Version labeled `json:"im:version"`
	Author  struct {
		Name labeled `json:"name"`
	} `json:"author"`
}

func toReview(e feedEntry) (review.Review, error) {
	score, err := strconv.Atoi(e.Rating.Label)
	if err != nil {
		return review.Review{}, fmt.Errorf("bad rating %q: %w", e.Rating.Label, err)
	}
	submitted, err := time.Parse(time.RFC3339, e.Updated.Label)
	if err != nil {
		return review.Review{}, fmt.Errorf("bad timestamp %q: %w", e.Updated.Label, err)
	}
	return review.Review{
		ID:          e.ID.Label,
		Author:      e.Author.Name.Label,
		Score:       score,
		Title:       e.Title.Label,
		Content:     e.Content.Label,
		SubmittedAt: submitted,
		AppVersion:  e.Version.Label,
	}, nil
}

func FetchReviews(appID string) ([]review.Review, error) {
	url := fmt.Sprintf(
		"https://itunes.apple.com/us/rss/customerreviews/id=%s/sortBy=mostRecent/page=1/json",
		appID,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", appID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for app %s", resp.StatusCode, appID)
	}

	var fr feedResponse
	if err := json.NewDecoder(resp.Body).Decode(&fr); err != nil {
		return nil, fmt.Errorf("decode error %s: %w", appID, err)
	}

	reviews := make([]review.Review, 0, len(fr.Feed.Entry))
	for _, e := range fr.Feed.Entry {
		r, err := toReview(e)
		if err != nil {
			return nil, fmt.Errorf("Parsing error %q: %w", e.ID.Label, err)
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}