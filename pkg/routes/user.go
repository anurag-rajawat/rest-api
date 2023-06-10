package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/anurag-rajawat/rest-api/pkg/handlers"
)

// RegisterRoutes creates all the required routes
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.Handle(http.MethodGet, "/", handlers.HomeHandler())
	v1 := router.Group("/v1")

	v1.GET("/", handlers.HomeHandler())
	v1.POST("/signup", handlers.SignUpHandler(db))
	v1.POST("/signin", handlers.SignInHandler(db))
	v1.GET("/users", handlers.GetUsersHandler(db))
	v1.GET("/users/:id", handlers.GetUserHandler(db))
	v1.PUT("/users/:id", handlers.UpdateUserHandler(db))
	v1.DELETE("/users/:id", handlers.DeleteUserHandler(db))
}
