package models

import "github.com/jmoiron/sqlx"

// Repo - kicker repository interface
type Repo interface {
	FindGame(id int64) (*Game, error)
	FindGames(page int, pageSize int64) ([]Game, error)
	CountGames() (int64, error)
}

// DbRepo - rerository struct for database
type DbRepo struct {
	*sqlx.DB
}

// NewDbRepo returns DbRepo instance
func NewDbRepo(dialect string, dsn string) (*DbRepo, error) {
	// open and connect at the same time, panicing on error
	db, err := sqlx.Connect(dialect, dsn)
	if err != nil {
		return nil, err
	}

	return &DbRepo{DB: db}, nil
}
