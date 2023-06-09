package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anurag-rajawat/rest-api/pkg/handlers"
)

// RegisterRoutes creates all the required routes
func RegisterRoutes(router *gin.Engine) {
	router.Handle(http.MethodGet, "/", handlers.Home)
	v1 := router.Group("/v1")

	v1.GET("/", handlers.Home)
	v1.POST("/signup", handlers.SignUp)
	v1.POST("/signin", handlers.SignIn)
	v1.GET("/users", handlers.GetUsers)
	v1.GET("/users/:id", handlers.GetUser)
	v1.PUT("/users/:id", handlers.UpdateUser)
	v1.DELETE("/users/:id", handlers.DeleteUser)
}
