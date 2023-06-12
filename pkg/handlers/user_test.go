package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/anurag-rajawat/rest-api/pkg/handlers"
	"github.com/anurag-rajawat/rest-api/pkg/types"
)

func TestGetUsersHandler(t *testing.T) {
	t.Run("without initialized DB should return error", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)

		// When
		handlers.GetUsersHandler(Db)(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, 1, len(ctx.Errors.Errors()))
		assert.Error(t, ctx.Errors.Last())
	})

	t.Run("without users should return no users", func(t *testing.T) {
		// Given
		err := InitDb()
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		expected := map[string][]types.User{
			"users": {},
		}

		// When
		handlers.GetUsersHandler(Db)(ctx)

		// Then
		var got map[string][]types.User
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusOK, w.Code)

		t.Cleanup(CleanDb)
	})

	t.Run("with a user should return that user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		expected := map[string][]types.User{
			"users": {user},
		}

		// When
		handlers.GetUsersHandler(Db)(ctx)

		// Then
		var got map[string][]types.User
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusOK, w.Code)

		t.Cleanup(CleanDb)
	})
}

func TestGetUserHandler(t *testing.T) {
	t.Run("with invalid ID should return no user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "invalid",
		}}
		expected := map[string]string{
			"error": "invalid ID",
		}

		// When
		handlers.GetUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("with valid ID of a user should return valid user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "1",
		}}
		expected := map[string]types.User{
			"user": user,
		}

		// When
		handlers.GetUserHandler(Db)(ctx)

		// Then
		var got map[string]types.User
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})

	t.Run("with valid ID but not of a user should return no user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "99",
		}}
		expected := map[string]string{
			"error": "user not found",
		}

		// When
		handlers.GetUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})
}

func TestUpdateUserHandler(t *testing.T) {
	t.Run("with invalid ID should not update user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "invalid",
		}}
		data := types.User{
			UserName: "updatedname",
			Email:    "updatedemail@gmail.com",
			Password: "updatedpassword",
		}
		dataBytes, err := json.Marshal(&data)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(dataBytes))
		expected := map[string]string{
			"error": "invalid ID",
		}

		// When
		handlers.UpdateUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("with valid ID of a user but invalid update data should not update user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "1",
		}}
		ctx.Request.Body = nil
		expected := map[string]string{
			"error": "invalid request",
		}

		// When
		handlers.UpdateUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})

	t.Run("with valid ID but not of a user should not update user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "99",
		}}
		data := types.User{
			UserName: "updatedname",
			Email:    "updatedemail@gmail.com",
			Password: "updatedpassword",
		}
		dataBytes, err := json.Marshal(&data)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(dataBytes))
		expected := map[string]string{
			"error": "user not found",
		}

		// When
		handlers.UpdateUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})

	t.Run("with valid ID of a user and update data should update user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "1",
		}}
		data := types.User{
			UserName: "updatedname",
			Email:    "updatedemail@gmail.com",
			Password: "updatedpassword",
		}
		dataBytes, err := json.Marshal(&data)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(dataBytes))

		// When
		handlers.UpdateUserHandler(Db)(ctx)

		// Then
		var got map[string]types.User
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		t.Cleanup(CleanDb)
	})
}

func TestDeleteUserHandler(t *testing.T) {
	t.Run("with invalid ID should not delete user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "invalid",
		}}
		expected := map[string]string{
			"error": "invalid ID",
		}

		// When
		handlers.DeleteUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("with valid ID but not of a user should not delete user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "99",
		}}
		expected := map[string]string{
			"error": "user not found",
		}

		// When
		handlers.DeleteUserHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})

	t.Run("with valid ID of a user should delete user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		_ = SeedOneUser()
		ctx.Params = []gin.Param{{
			Key:   "id",
			Value: "1",
		}}

		// When
		handlers.DeleteUserHandler(Db)(ctx)

		// Then
		assert.Equal(t, http.StatusNoContent, ctx.Writer.Status())

		t.Cleanup(CleanDb)
	})
}
