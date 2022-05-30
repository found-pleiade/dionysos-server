package routes

// var validate = validator.New()

// var userCollection *mongo.Collection = OpenCollection(Client, "users")

// UserCreate adds a user from the JSON received in the request body.
// func UserCreate(c *gin.Context) {

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	var user models.User

// 	if err := c.BindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}

// 	validationErr := validate.Struct(user)
// 	if validationErr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
// 		fmt.Println(validationErr)
// 		return
// 	}
// 	user.ID = primitive.NewObjectID()

// 	result, insertErr := userCollection.InsertOne(ctx, user)
// 	if insertErr != nil {
// 		msg := fmt.Sprintf("user item was not created")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 		fmt.Println(insertErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }

// // UserRemove removes a user from the JSON received in the request body.
// func UserRemove(c *gin.Context) {

// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()

// 	var user models.User

// 	if err := c.BindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}

// 	validationErr := validate.Struct(user)
// 	if validationErr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
// 		fmt.Println(validationErr)
// 		return
// 	}
// 	user.ID = primitive.NewObjectID()

// 	result, insertErr := userCollection.InsertOne(ctx, user)
// 	if insertErr != nil {
// 		msg := fmt.Sprintf("user item was not created")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 		fmt.Println(insertErr)
// 		return
// 	}

// 	c.JSON(http.StatusOK, result)
// }
