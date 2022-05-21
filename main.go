package main

import (
	"net/http"

	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/users", routes.UserPost)
	router.GET("/ping", getping)

	router.Run(":8080")
}

func getping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "pong"})
}
