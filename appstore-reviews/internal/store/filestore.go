package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"appstore-reviews/internal/review"
)

type FileStore struct {
	dir string

	mu      sync.RWMutex
	reviews map[string]map[string]review.Review
}

func New() (*FileStore, error) {
	return NewFileStore("data")
}

func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data folder: %w", err)
	}

	fs := &FileStore{
		dir:     dir,
		reviews: make(map[string]map[string]review.Review),
	}

	if err := fs.load(); err != nil {
		return nil, fmt.Errorf("failed to load existing reviews: %w", err)
	}
	return fs, nil
}

func (fs *FileStore) load() error {
	entries, err := os.ReadDir(fs.dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		appID := e.Name()[:len(e.Name())-len(".json")]

		data, err := os.ReadFile(filepath.Join(fs.dir, e.Name()))
		if err != nil {
			return err
		}

		var list []review.Review
		if err := json.Unmarshal(data, &list); err != nil {
			return fmt.Errorf("parse %s: %w", e.Name(), err)
		}

		byID := make(map[string]review.Review, len(list))
		for _, r := range list {
			byID[r.ID] = r
		}
		fs.reviews[appID] = byID
	}
	return nil
}

func (fs *FileStore) Save(appID string, incoming []review.Review) (int, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	byID, ok := fs.reviews[appID]
	if !ok {
		byID = make(map[string]review.Review)
		fs.reviews[appID] = byID
	}

	added := 0
	for _, r := range incoming {
		if _, exists := byID[r.ID]; !exists {
			byID[r.ID] = r
			added++
		}
	}

	if added == 0 {
		return 0, nil
	}

	if err := fs.persistLocked(appID); err != nil {
		return 0, err
	}

	return added, nil
}

func (fs *FileStore) Get(appID string) []review.Review {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	byID := fs.reviews[appID]
	list := make([]review.Review, 0, len(byID))
	for _, r := range byID {
		list = append(list, r)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].SubmittedAt.After(list[j].SubmittedAt)
	})
	return list
}

func (fs *FileStore) persistLocked(appID string) error {
	list := make([]review.Review, 0, len(fs.reviews[appID]))
	for _, r := range fs.reviews[appID] {
		list = append(list, r)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	finalPath := filepath.Join(fs.dir, appID+".json")
	tmpPath := finalPath + ".tmp"

	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, finalPath)
}