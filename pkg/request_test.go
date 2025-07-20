package pkg

import (
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/Thanus-Kumaar/anokha-2025-backend/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testRouter *gin.Engine

func TestMain(m *testing.M) {
	testRouter = tests.SetupTestRouter()
	exitCode := m.Run()
	os.Exit(exitCode)
}

type MockGinContext struct {
	*gin.Context
	mock.Mock
}

func TestTagRequestWithId(t *testing.T) {
	t.Run("Test::set_request_id_and_call_next", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(w, testRouter)
		mockContext := &MockGinContext{Context: c}

		TagRequestWithId(mockContext.Context)
		reqID, exists := mockContext.Get("request_id")
		assert.True(t, exists)
		assert.NotEmpty(t, reqID)

		reqIDStr, ok := reqID.(string)
		assert.True(t, ok)
		ksuidRegex := regexp.MustCompile(`^[0-9A-Za-z]{27}$`)
		assert.True(t, ksuidRegex.MatchString(reqIDStr))
	})

	t.Run("Test::each_call_generates_unique_ksuid", func(t *testing.T) {
		w1 := httptest.NewRecorder()
		c1 := gin.CreateTestContextOnly(w1, testRouter)
		mockContext1 := &MockGinContext{Context: c1}
		TagRequestWithId(mockContext1.Context)
		id1, _ := mockContext1.Get("request_id")

		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		mockContext2 := &MockGinContext{Context: c2}
		TagRequestWithId(mockContext2.Context)
		id2, _ := mockContext2.Get("request_id")

		assert.NotEqual(t, id1, id2)
	})
}

func TestGrabRequestId(t *testing.T) {
	t.Run("Test::retrieve_existing_request_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(w, testRouter)

		expectedID := "test-ksuid-1234567890ABCDEF"
		c.Set("request_id", expectedID)

		retrievedID := GrabRequestId(c)
		assert.Equal(t, expectedID, retrievedID)
	})

	t.Run("Test::return_missing_id_when_not_set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(w, testRouter)

		retrievedID := GrabRequestId(c)
		assert.Equal(t, "missing-id", retrievedID)
	})

	t.Run("Test::non_string_request_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(w, testRouter)

		c.Set("request_id", 12345)

		retrievedID := GrabRequestId(c)
		assert.Equal(t, "12345", retrievedID)
	})
}
