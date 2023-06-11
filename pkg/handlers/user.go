package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/anurag-rajawat/rest-api/pkg/types"
)

func GetUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user types.User
		users, err := user.FindAll(db)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}

func GetUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := validateId(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var user types.User
		usr, err := user.FindById(db, id)
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
}

func UpdateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := validateId(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
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

		updatedUser, err := user.Update(db, id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"user": updatedUser,
		})
	}
}

func DeleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := validateId(idStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var user types.User
		if err = user.Delete(db, id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func validateId(idStr string) (uint64, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid ID")
	}
	return id, nil
}
