package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func main() {
	_, err := getNeoDriver()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.GET("/ping", getping)
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "pong"})
}

func getNeoDriver() (neo4j.Driver, error) {
	neo4jURI, found := os.LookupEnv("NEO4J_URI")
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

	driver, err := neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	return driver, err
}
