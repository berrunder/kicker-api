package models

// MockDbRepo - mocks DbRepo for tests
type MockDbRepo struct {
	Games map[int64]Game
}

// NewMockDb creates mockDbRepo instance
func NewMockDb() *MockDbRepo {
	games := make(map[int64]Game)

	return &MockDbRepo{games}
}

// FindGame by id, nil if not found
func (repo MockDbRepo) FindGame(id int64) (*Game, error) {
	game, has := repo.Games[id]
	if !has {
		return nil, nil
	}

	return &game, nil
}

// FindGames returns all games, paged by page and pageSize
func (repo MockDbRepo) FindGames(page int, pageSize int64) ([]Game, error) {
	var games []Game
	for _, v := range repo.Games {
		games = append(games, v)
	}

	// Sort by PlayedAt time
	GameBy(func(a, b Game) bool {
		return a.PlayedAt.After(b.PlayedAt)
	}).Sort(games)

	cnt := int64(len(games))

	from := minPositive(int64(page)*pageSize, cnt)

	to := minPositive((int64(page)+1)*pageSize, cnt)

	return games[from:to], nil
}

func minPositive(a, b int64) int64 {
	if a > b {
		if b > 0 {
			return b
		}
		return 0
	}

	if a > 0 {
		return a
	}
	return 0
}

// CountGames returns total count of games in database
func (repo MockDbRepo) CountGames() (int64, error) {
	return int64(len(repo.Games)), nil
}
