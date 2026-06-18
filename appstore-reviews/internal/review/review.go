package review

import "time"

type Review struct {
	ID          string    `json:"id"`
	Author      string    `json:"author"`
	Score       int       `json:"score"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	SubmittedAt time.Time `json:"submittedAt"`
	AppVersion  string    `json:"appVersion"`
}