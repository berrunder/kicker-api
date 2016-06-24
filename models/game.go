package models

import "time"

// Game type
type Game struct {
	ID       int64 `json:"id"`
	Team1    *Team
	Team2    *Team
	Team1ID  int64     `db:"team1_id" json:"-"`
	Team2ID  int64     `db:"team2_id" json:"-"`
	Score1   int8      `json:"score1"`
	Score2   int8      `json:"score2"`
	PlayedAt time.Time `json:"playedAt" db:"played_at"`
}
