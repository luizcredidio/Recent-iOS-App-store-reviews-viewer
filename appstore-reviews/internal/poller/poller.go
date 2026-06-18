package poller

import (
	"context"
	"log"
	"time"

	"appstore-reviews/internal/appstore"
	"appstore-reviews/internal/store"
)

type Poller struct {
	store    *store.FileStore
	appIDs   []string
	interval time.Duration
}

func New(st *store.FileStore, appIDs []string, interval time.Duration) *Poller {
	return &Poller{
		store:    st,
		appIDs:   appIDs,
		interval: interval,
	}
}

func (p *Poller) Run(ctx context.Context) {
	p.pollAll()

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("poller shutting down")
			return
		case <-ticker.C:
			p.pollAll()
		}
	}
}

func (p *Poller) pollAll() {
	for _, appID := range p.appIDs {
		reviews, err := appstore.FetchReviews(appID)
		if err != nil {
			log.Printf("poller: fetch failed for app %s: %v", appID, err)
			continue
		}

		added, err := p.store.Save(appID, reviews)
		if err != nil {
			log.Printf("failted to save for app %s in poller: %v", appID, err)
			continue
		}

		log.Printf("poller: app %s — fetched %d, %d new", appID, len(reviews), added)
	}
}