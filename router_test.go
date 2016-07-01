package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/berrunder/kicker-api/models"
)

func setupMockDb() *models.MockDbRepo {
	repo := models.NewMockDb()
	team1 := &models.Team{ID: 1, Name: "team1"}
	team2 := &models.Team{ID: 2, Name: "team2"}
	repo.Games[1] = models.Game{ID: 1, Team1: team1, Team2: team2, Team1ID: team1.ID, Team2ID: team2.ID, PlayedAt: time.Now(), Score1: 1, Score2: 2}
	repo.Games[2] = models.Game{ID: 2, Team1: team1, Team2: team2, Team1ID: team1.ID, Team2ID: team2.ID, PlayedAt: time.Now().AddDate(0, 0, -1), Score1: 2, Score2: 1}
	return repo
}

func TestEmptyGamesIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/score", nil)

	repo := models.NewMockDb()
	router := SetupRouter(&DbHandler{repo}, true)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v, %v expected", rec.Code, http.StatusOK)
	}

	expected := `{"data":[],"page":0,"total":0}`
	trimmed := strings.TrimSpace(rec.Body.String())

	if trimmed != expected {
		t.Errorf("\n...expected = '%v'\n...obtained = '%v'", expected, trimmed)
	}
}

func TestGamesIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/score", nil)

	repo := setupMockDb()
	router := SetupRouter(&DbHandler{repo}, true)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Non-expected status code %v, %v expected", rec.Code, http.StatusOK)
	}

	var games []models.Game
	games = append(games, repo.Games[1])
	games = append(games, repo.Games[2])

	expectedMap := make(map[string]interface{})
	expectedMap["data"] = games
	expectedMap["page"] = 0
	expectedMap["total"] = len(games)

	expected, err := json.Marshal(expectedMap)
	if err != nil {
		panic(err)
	}

	trimmed := strings.TrimSpace(rec.Body.String())

	if trimmed != string(expected) {
		t.Errorf("\n...expected = '%v'\n...obtained = '%v'", string(expected), trimmed)
	}
}
