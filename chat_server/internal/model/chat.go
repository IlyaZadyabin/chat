package model

import "time"

type ChatCreate struct {
	Usernames []string
}

type Message struct {
	From      string
	Text      string
	Timestamp time.Time
}
