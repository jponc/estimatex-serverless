package types

import (
	"time"
)

type Room struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	EndedAt   time.Time `json:"ended_at"`
}

type Participant struct {
	RoomID     string    `json:"room_id"`
	Name       string    `json:"name"`
	IsAdmin    bool      `json:"is_admin"`
	LatestVote string    `json:"latest_vote"`
	CreatedAt  time.Time `json:"created_at"`
}

type ParticipantArr []Participant
