package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/anurag-rajawat/rest-api/pkg/types"
	"github.com/anurag-rajawat/rest-api/pkg/utils"
)

func SignUp(ctx *gin.Context) {
	var user types.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save password before saving user to generate bearer token
	password := user.Password

	newUser, err := user.Create(utils.GetDb())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// After signup login user as well
	token, err := user.GetToken(utils.GetDb(), user.Email, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Set("Authorization", "Bearer"+token)

	ctx.JSON(http.StatusCreated, gin.H{
		"user": newUser,
	})
}

func SignIn(ctx *gin.Context) {
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	if email == "" || password == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide credentials",
		})
		return
	}

	var user types.User
	token, err := user.GetToken(utils.GetDb(), email, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Set("Authorization", "Bearer"+token)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully logged in",
		"token":   token,
	})
}
