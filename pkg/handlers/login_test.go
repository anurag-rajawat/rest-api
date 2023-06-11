package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/anurag-rajawat/rest-api/pkg/handlers"
	"github.com/anurag-rajawat/rest-api/pkg/types"
)

func TestSignInHandler(t *testing.T) {
	t.Run("with valid credentials should sign-in user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		// save original password otherwise it will check hashed password
		user.Password = "passwd"
		body, err := json.Marshal(&user)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// When
		handlers.SignInHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Len(t, got, 2)
		assert.Len(t, w.Result().Cookies(), 1)
		assert.Contains(t, w.Result().Cookies()[0].Name, "Authorization")

		t.Cleanup(CleanDb)
	})

	t.Run("with no credentials should not sign in user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		userRequest := types.User{}
		body, err := json.Marshal(&userRequest)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		expected := map[string]string{
			"error": "please provide valid credentials",
		}

		// When
		handlers.SignInHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})

	t.Run("with invalid request should not sign in user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		expected := map[string]string{
			"error": "invalid request",
		}

		// When
		handlers.SignInHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("with incorrect email should not sign in user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		user.Email = "incorrect@gmail.com"
		// save original password otherwise it will check hashed password
		user.Password = "passwd"
		body, err := json.Marshal(&user)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		expected := map[string]string{
			"error": "user not found",
		}

		// When
		handlers.SignInHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})

	t.Run("with incorrect password should not sign in user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		user.Password = "incorrect"
		body, err := json.Marshal(&user)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		expected := map[string]string{
			"error": "invalid credentials",
		}

		// When
		handlers.SignInHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, expected, got)

		t.Cleanup(CleanDb)
	})
}

func TestSignUpHandler(t *testing.T) {
	t.Run("with invalid request should not signup user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		expected := map[string]string{
			"error": "invalid request",
		}

		// When
		handlers.SignUpHandler(Db)(ctx)

		// Then
		var got map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("with no credentials should not signup user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		userRequest := types.User{}
		body, err := json.Marshal(&userRequest)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		expected := map[string]string{
			"error": "please provide required details",
		}

		// When
		handlers.SignUpHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, expected, got)
	})

	t.Run("with valid request data should signup user", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		err := InitDb()
		if err != nil {
			t.Error(err)
		}
		userRequest := types.User{
			UserName: "user1",
			Email:    "user1@gmail.com",
			Password: "user1",
		}
		body, err := json.Marshal(&userRequest)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// When
		handlers.SignUpHandler(Db)(ctx)

		// Then
		var got map[string]types.User
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Len(t, got, 1)
		t.Cleanup(CleanDb)
	})

	t.Run("duplicate user should not signup", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		ctx := GetTestGinContext(w)
		user := SeedOneUser()
		user.ID = 2
		user.Password = "passwd"
		body, err := json.Marshal(&user)
		if err != nil {
			t.Error(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// When
		handlers.SignUpHandler(Db)(ctx)

		// Then
		var got map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &got)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, got["error"], "ERROR: duplicate key")
		log.Info(got)

		t.Cleanup(CleanDb)
	})
}
