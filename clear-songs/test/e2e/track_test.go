package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/RubenPari/clear-songs/internal/application/shared/dto"
	"github.com/RubenPari/clear-songs/internal/infrastructure/di"
	httptransport "github.com/RubenPari/clear-songs/internal/infrastructure/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnv() {
	os.Setenv("CLIENT_ID", "mock_id")
	os.Setenv("CLIENT_SECRET", "mock_secret")
	os.Setenv("REDIRECT_URL", "http://127.0.0.1:3000/callback")
}

func TestTrackAPI_E2E(t *testing.T) {
	setupTestEnv()
	
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Create router
	router := gin.Default()
	
	container, err := di.NewContainer()
	if err != nil {
		t.Skipf("Skipping E2E test: DI container failed: %v", err)
		return
	}
	
	httptransport.SetUpRoutes(router, container)

	t.Run("GET /track/summary - Unauthenticated", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/track/summary", nil)
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		
		var response dto.APIResponse[any]
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		
		require.NotNil(t, response.Error, "Response error should not be nil")
		assert.Equal(t, "UNAUTHORIZED", response.Error.Code)
	})

	// Note: In a real E2E we would mock the session/auth middleware to inject a fake user
	// to test the internal logic (like validation errors). For now we verify it's blocked.
}
