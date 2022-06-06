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
	router.PATCH("/users/:id", routes.UpdateUser)
	router.DELETE("/users/:id", routes.DeleteUser)
<<<<<<< HEAD
	router.POST("/rooms/", routes.CreateRoom)
	router.GET("/rooms/:id", routes.GetRoom)
	router.PATCH("/rooms/:id", routes.UpdateRoom)
	router.DELETE("/rooms/:id", routes.DeleteRoom)
=======
>>>>>>> dd476d1 (Add methods to update and delete user)
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"response": "v0.1.0"})
}
