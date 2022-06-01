package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	neo4jUri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}

	driver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))

	router := gin.Default()

	router.GET("/ping", getping)
	router.Run(":8080")
}

func getping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "pong"})
}

func driver(target string, token neo4j.AuthToken) neo4j.Driver {
	result, err := neo4j.NewDriver(target, token)
	if err != nil {
		panic(err)
	}
	return result
}
