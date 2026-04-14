package track

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRangeQuery(t *testing.T) {
	t.Parallel()

	t.Run("nil request", func(t *testing.T) {
		min, max, errMsg := ValidateRangeQuery(nil)
		assert.Equal(t, 0, min)
		assert.Equal(t, 0, max)
		assert.Equal(t, "", errMsg)
	})

	t.Run("empty", func(t *testing.T) {
		req := &RangeRequest{}
		min, max, errMsg := ValidateRangeQuery(req)
		assert.Equal(t, 0, min)
		assert.Equal(t, 0, max)
		assert.Equal(t, "", errMsg)
	})

	t.Run("min only", func(t *testing.T) {
		m := 5
		req := &RangeRequest{Min: &m}
		min, max, errMsg := ValidateRangeQuery(req)
		assert.Equal(t, 5, min)
		assert.Equal(t, 0, max)
		assert.Equal(t, "", errMsg)
	})

	t.Run("min and genre without max is valid", func(t *testing.T) {
		m := 5
		req := &RangeRequest{Min: &m, Genre: "Rock"}
		min, max, errMsg := ValidateRangeQuery(req)
		assert.Equal(t, 5, min)
		assert.Equal(t, 0, max)
		assert.Equal(t, "", errMsg)
	})

	t.Run("max only", func(t *testing.T) {
		m := 10
		req := &RangeRequest{Max: &m}
		min, max, errMsg := ValidateRangeQuery(req)
		assert.Equal(t, 0, min)
		assert.Equal(t, 10, max)
		assert.Equal(t, "", errMsg)
	})

	t.Run("min greater than max", func(t *testing.T) {
		a, b := 20, 10
		req := &RangeRequest{Min: &a, Max: &b}
		_, _, errMsg := ValidateRangeQuery(req)
		assert.NotEmpty(t, errMsg)
	})

	t.Run("min equals max with both positive", func(t *testing.T) {
		v := 7
		req := &RangeRequest{Min: &v, Max: &v}
		min, max, errMsg := ValidateRangeQuery(req)
		assert.Equal(t, 7, min)
		assert.Equal(t, 7, max)
		assert.Equal(t, "", errMsg)
	})
}

func TestRangeRequest_BindQuery_MinWithoutMax_WithGenre(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/track/summary?min=5&genre=Rock", nil)

	var req RangeRequest
	err := c.ShouldBindQuery(&req)
	require.NoError(t, err)
	require.NotNil(t, req.Min)
	assert.Equal(t, 5, *req.Min)
	assert.Nil(t, req.Max)
	assert.Equal(t, "Rock", req.Genre)

	min, max, errMsg := ValidateRangeQuery(&req)
	assert.Equal(t, "", errMsg)
	assert.Equal(t, 5, min)
	assert.Equal(t, 0, max)
}
