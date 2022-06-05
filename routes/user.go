package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Brawdunoir/dionysos-server/models"
	"github.com/gin-gonic/gin"
)

// Create a user in the aganro database
func CreateUser(c *gin.Context) {
	client = getClient()
	db = getDatabase("dionysos")
	SetupCollection()

	var user models.User
	c.BindJSON(&user)

	col, err := db.Collection(context.TODO(), cols[0])
	if err != nil {
		fmt.Println(err)
	}

	metas, errs, err := col.CreateDocuments(context.TODO(), []models.User{user})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		log.Fatalf("Failed to create documents: %v", err)
	} else if err := errs.FirstNonNil(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User not created"})
		log.Fatalf("Failed to create documents: first error: %v", err)
	}

	fmt.Printf("Created documents with keys '%s' in collection '%s' in database '%s'\n", strings.Join(metas.Keys(), ","), col.Name(), db.Name())

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
