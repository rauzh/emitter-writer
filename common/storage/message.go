package storage

import "time"

type Message struct {
	IsDiscontinued bool      `json:"discounted"`
	ID             uint      `json:"id"`
	ASIN           string    `json:"ASIN"`
	Title          string    `json:"title"`
	Group          string    `json:"group"`
	Salesrank      uint64    `json:"salesrank"`
	Similarsnum    uint      `json:"similarsnum"`
	Similars       []string  `json:"similars"`
	Time           time.Time `json:"time"`
}
