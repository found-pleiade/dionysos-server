package main

import (
	"net/http"

	"github.com/Brawdunoir/dionysos-server/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/version", getVersion)
	router.POST("/users/", routes.CreateUser)
	router.GET("/users/:id", routes.GetUser)
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "v0.1.0"})
}
