//nolint:typecheck
package routes

import (
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// SetupRouter sets up the router
func SetupRouter(router *gin.Engine) *gin.Engine {

	binding.Validator = new(defaultValidator)

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
	}))

	userRouter := router.Group("/users")
	{
		userRouter.POST("/", CreateUser)
		userRouter.GET("/:id", GetUser)
		userRouter.PATCH("/:id", UpdateUser)
		userRouter.DELETE("/:id", DeleteUser)
	}

	roomRouter := router.Group("/rooms")
	{
		roomRouter.POST("/", CreateRoom)
		roomRouter.GET("/:id", GetRoom)
		roomRouter.PATCH("/:id", UpdateRoom)
		roomRouter.DELETE("/:id", DeleteRoom)
	}

	router.GET("/version", func(c *gin.Context) {
		var version string
		if version = os.Getenv("VERSION"); version == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Version not set"})
		}
		c.String(http.StatusOK, version)
	})

	return router
}

/* Functions to pair playground validator with gin own validator using default binding tag */
type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
