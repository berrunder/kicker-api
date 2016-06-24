package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/berrunder/kicker-api/models"

	"gopkg.in/urfave/cli.v2"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var repo *DbRepo

func runServer(c *cli.Context) error {
	// open and connect at the same time, panicing on error
	db := sqlx.MustConnect("mysql", c.String("datasource"))
	repo = &DbRepo{DB: db}

	router := SetupRouter()

	port := ":" + c.String("port")
	log.Printf("Listening on %v...\n", port)
	router.Run(port)

	return nil
}

// SetupRouter for API
func SetupRouter() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(CORSMiddleware())

	v1 := router.Group("/v1")
	{
		v1.GET("/score", getScores)
		v1.GET("/score/:id", getScore)
		v1.OPTIONS("/score", optionsHandler)
		v1.OPTIONS("/score/:id", optionsHandler)
	}

	return router
}

const pageSize = 100

func getScores(c *gin.Context) {

	var scores []models.Game
	var err error

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

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  scores,
		"page":  page,
	})
}

func getScore(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, getAPIError("Bad request", err))
		return
	}

	s, err := repo.FindGame(id)

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
