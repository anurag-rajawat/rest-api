package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/anurag-rajawat/rest-api/pkg/types"
	"github.com/anurag-rajawat/rest-api/pkg/utils"
)

func GetUsers(ctx *gin.Context) {
	var user types.User
	users, err := user.FindAll(utils.GetDb())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := validateId(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user types.User
	usr, err := user.FindById(utils.GetDb(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": usr,
	})
}

func UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := validateId(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var user types.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedUser, err := user.Update(utils.GetDb(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"user": updatedUser,
	})
}

func DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := validateId(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var user types.User
	if err = user.Delete(utils.GetDb(), id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func validateId(idStr string) (uint64, error) {
	if idStr == "" {
		return 0, errors.New("please provide ID")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID")
	}
	return id, nil
}
