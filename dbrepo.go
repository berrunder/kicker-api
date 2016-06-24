package main

import (
	"database/sql"

	"github.com/berrunder/kicker-api/models"
	"github.com/jmoiron/sqlx"
)

// DbRepo - rerository struct for database
type DbRepo struct {
	DB *sqlx.DB
}

// FindGame by id, nil if not found
func (repo DbRepo) FindGame(id int64) (*models.Game, error) {
	g := &models.Game{}
	g.Team1 = &models.Team{}
	g.Team2 = &models.Team{}

	row := repo.DB.QueryRow(`SELECT g.id, g.score1, g.score2, g.played_at, t1.id, t1.name, t2.id, t2.name
	  FROM game g 
	  JOIN team t1 on t1.id = g.team1_id
	  JOIN team t2 on t2.id = g.team2_id
	  WHERE g.id = ?`,
		id)

	err := row.Scan(&g.ID, &g.Score1, &g.Score2, &g.PlayedAt, &g.Team1.ID, &g.Team1.Name, &g.Team2.ID, &g.Team2.Name)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return g, nil
}

// FindGames returns all games, paged by page and pageSize
func (repo DbRepo) FindGames(page int, pageSize int64) ([]models.Game, error) {
	var scores []models.Game

	rows, err := repo.DB.Query(`SELECT g.id, g.score1, g.score2, g.played_at, t1.id, t1.name, t2.id, t2.name
	  FROM game g 
	  JOIN team t1 on t1.id = g.team1_id
	  JOIN team t2 on t2.id = g.team2_id
	  ORDER BY g.played_at DESC LIMIT ?, ?`, page, pageSize)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := models.Game{}
		s.Team1 = &models.Team{}
		s.Team2 = &models.Team{}
		if err := rows.Scan(&s.ID, &s.Score1, &s.Score2, &s.PlayedAt, &s.Team1.ID, &s.Team1.Name, &s.Team2.ID, &s.Team2.Name); err != nil {
			return nil, err
		}

		scores = append(scores, s)
	}

	return scores, nil
}

// CountGames returns total count of games in database
func (repo DbRepo) CountGames() (int64, error) {
	var total int64

	err := repo.DB.Get(&total, "SELECT COUNT(*) FROM game")

	if err != nil {
		return 0, err
	}

	return total, nil
}
