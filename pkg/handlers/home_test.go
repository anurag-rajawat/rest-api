package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anurag-rajawat/rest-api/pkg/handlers"
)

func TestHome(t *testing.T) {
	tests := []struct {
		name       string
		url        url.URL
		statusCode int
	}{
		{
			name:       "RootEndPoint",
			url:        url.URL{Path: "/"},
			statusCode: http.StatusPermanentRedirect,
		},
		{
			name:       "V1RootEndPoint",
			url:        url.URL{Path: "/v1"},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			w := httptest.NewRecorder()
			ctx := GetTestGinContext(w)
			ctx.Request.URL = &tt.url
			expected := map[string]string{
				"message": "Namaste World! 🙏",
			}

			// When
			handlers.HomeHandler()(ctx)

			// Then
			var got map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, expected, got)
		})
	}
}
