package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(ctx *gin.Context) {
	if ctx.Request.URL.String() == "/" {
		ctx.JSON(http.StatusPermanentRedirect, gin.H{
			"message": "Namaste World! 🙏",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Namaste World! 🙏",
	})
}
