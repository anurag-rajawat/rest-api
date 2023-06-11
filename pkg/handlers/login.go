package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/anurag-rajawat/rest-api/pkg/types"
)

func SignUpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRequest types.User
		if err := ctx.BindJSON(&userRequest); err != nil {
			log.Warn(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		if userRequest.UserName == "" || userRequest.Email == "" || userRequest.Password == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "please provide required details",
			})
			return
		}

		newUser, err := userRequest.Create(db)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"user": newUser,
		})
	}
}

func SignInHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userRequest types.User
		if err := ctx.BindJSON(&userRequest); err != nil {
			log.Warn(err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})
			return
		}

		email := userRequest.Email
		password := userRequest.Password

		if email == "" || password == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "please provide valid credentials",
			})
			return
		}

		user, err := userRequest.FindByEmail(db, email)
		if err != nil {
			log.Error(err)
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = user.CheckPassword(user.Password, password)
		if err != nil {
			log.Warn(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		token, err := user.GetJwtToken()
		if err != nil {
			log.Error(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.SetCookie("Authorization", token, 24*3600, ctx.Request.URL.Path, "", false, true)
		ctx.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"type":         "bearer",
		})
	}
}
