package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/berrunder/kicker-api/models"

	"github.com/urfave/cli/v2"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const configFile = "dbconfig.yml"
const defaultEnv = "development"
const defaultDialect = "mysql"
const pageSize = 100

// DbHandler stores data repository dependency
type DbHandler struct {
	repo models.Repo
}

func runServer(c *cli.Context) error {
	var dialect string
	var dsn string

	dialect = c.String("dialect")
	dsn = c.String("datasource")

	env := os.Getenv("KICKER_ENV")
	if env == "" {
		env = defaultEnv
	}

	config, err := GetConfig(configFile, env)
	if err != nil {
		log.Printf("Error reading config: %v", err)
	} else if config != nil {
		if dialect == defaultDialect {
			dialect = config.Dialect
		}
		if dsn == "" {
			dsn = config.DataSource
		}
	}

	// open and connect at the same time, panicing on error
	repo, err := models.NewDbRepo(dialect, dsn)
	if err != nil {
		log.Fatalf("Error connecting to database:\n\t%v", err)
	}

	router := SetupRouter(&DbHandler{repo}, true)

	port := ":" + c.String("port")
	log.Printf("Listening on %v...\n", port)
	router.Run(port)

	return nil
}

// SetupRouter for API
func SetupRouter(handler *DbHandler, useLogger bool) *gin.Engine {
	router := gin.New()

	// Global middleware
	if useLogger {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/score", handler.getScores)
		v1.GET("/score/:id", handler.getScore)
		v1.OPTIONS("/score", optionsHandler)
		v1.OPTIONS("/score/:id", optionsHandler)
	}

	return router
}

func (h *DbHandler) getScores(c *gin.Context) {

	var scores []models.Game
	var err error
	repo := h.repo

	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}

	total, err := repo.CountGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, getInternalError(err))
		return
	}

	scores, err = repo.FindGames(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, getInternalError(err))
		return
	}

	if scores == nil {
		scores = make([]models.Game, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  scores,
		"page":  page,
	})
}

func (h *DbHandler) getScore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, getAPIError("Bad request", err))
		return
	}

	s, err := h.repo.FindGame(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, getInternalError(err))
		return
	}

	if s == nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "Not found"})
		return
	}

	c.JSON(http.StatusOK, s)
}

func optionsHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
	c.Next()
}

// CORSMiddleware for gin
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	}
}
