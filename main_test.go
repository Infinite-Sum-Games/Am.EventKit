package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Thanus-Kumaar/anokha-2025-backend/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testRouter *gin.Engine

// Setup a test suite
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	testRouter = tests.SetupTestRouter()
	exitCode := m.Run() // Run all tests
	os.Exit(exitCode)
}

func TestCORSHeaders(t *testing.T) {
	t.Run("AllowedOrigin::exact_match", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://example.com")

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("AllowedOrigin::wildcard_subdomain_match", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://sub.example.com")

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.NotEqual(t, "http://sub.example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.NotEqual(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("DisallowedOrigin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://malicious.com")

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("NoOriginHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("AllowedOrigin::POST_method", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.Header.Set("Origin", "http://example.com")

		testRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})
}

func TestCORSOptionsRequest(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)

	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "GET,POST,DELETE,PUT,OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "X-Csrf-Token,Origin,Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "43200", w.Header().Get("Access-Control-Max-Age"))
}
