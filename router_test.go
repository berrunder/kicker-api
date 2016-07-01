package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/berrunder/kicker-api/models"
	"github.com/gin-gonic/gin"
)

func setupMockDb() *models.MockDbRepo {
	repo := models.NewMockDb()
	team1 := &models.Team{ID: 1, Name: "team1"}
	team2 := &models.Team{ID: 2, Name: "team2"}
	repo.Games[1] = models.Game{ID: 1, Team1: team1, Team2: team2, Team1ID: team1.ID, Team2ID: team2.ID, PlayedAt: time.Now(), Score1: 1, Score2: 2}
	repo.Games[2] = models.Game{ID: 2, Team1: team1, Team2: team2, Team1ID: team1.ID, Team2ID: team2.ID, PlayedAt: time.Now().AddDate(0, 0, -1), Score1: 2, Score2: 1}
	return repo
}

var useLogger = flag.Bool("log", false, "Enable logger middleware in Gin for test")

func getRouter(repo models.Repo) *gin.Engine {
	return SetupRouter(&DbHandler{repo}, *useLogger)
}

func TestEmptyGamesIndex(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/score", nil)

	repo := models.NewMockDb()
	router := getRouter(repo)

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
	router := getRouter(repo)

	router.ServeHTTP(rec, req)

	expectedCode := http.StatusOK
	if rec.Code != expectedCode {
		t.Fatalf("Non-expected status code %v, %v expected", rec.Code, expectedCode)
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

func TestGameNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/score/1", nil)

	repo := models.NewMockDb()
	router := getRouter(repo)

	router.ServeHTTP(rec, req)

	expectedCode := http.StatusNotFound
	if rec.Code != expectedCode {
		t.Fatalf("Non-expected status code %v, %v expected", rec.Code, expectedCode)
	}

	expected := []byte("{\"status\":\"Not found\"}\n")
	obtained := rec.Body.Bytes()

	if !bytes.Equal(expected, obtained) {
		t.Errorf("\n...expected = '%q'\n...obtained = '%q'", expected, obtained)
	}
}

func TestGetGame(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/score/1", nil)

	repo := setupMockDb()
	router := getRouter(repo)

	router.ServeHTTP(rec, req)

	expectedCode := http.StatusOK

	if rec.Code != expectedCode {
		t.Fatalf("Non-expected status code %v, %v expected", rec.Code, expectedCode)
	}

	expected, err := json.Marshal(repo.Games[1])
	if err != nil {
		panic(err)
	}
	expected = append(expected, []byte("\n")[0])

	obtained := rec.Body.Bytes()

	if !bytes.Equal(expected, obtained) {
		t.Errorf("\n...expected = '%q'\n...obtained = '%q'", expected, obtained)
	}
}

func TestOptions(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/v1/score/1", nil)

	repo := models.NewMockDb()
	router := getRouter(repo)

	router.ServeHTTP(rec, req)

	expectedMethods := "DELETE,POST,PUT"
	if !strings.Contains(strings.Join(rec.HeaderMap["Access-Control-Allow-Methods"], ""), expectedMethods) {
		t.Errorf("\n...expected to contain = '%q'\n...obtained = '%q'", expectedMethods, rec.HeaderMap["Access-Control-Allow-Methods"])
	}

	expectedOrigin := "*"
	if !strings.Contains(strings.Join(rec.HeaderMap["Access-Control-Allow-Origin"], ""), expectedOrigin) {
		t.Errorf("\n...expected to contain = '%q'\n...obtained = '%q'", expectedOrigin, rec.HeaderMap["Access-Control-Allow-Origin"])
	}

	expectedAllowHeaders := "Content-Type"
	if !strings.Contains(strings.Join(rec.HeaderMap["Access-Control-Allow-Headers"], ""), expectedAllowHeaders) {
		t.Errorf("\n...expected to contain '%q'\n...obtained = '%q'", expectedAllowHeaders, rec.HeaderMap["Access-Control-Allow-Headers"])
	}
}
