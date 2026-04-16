package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGenreBatchJSON_plainArray(t *testing.T) {
	raw := `[{"key":"abc123","genre":"hip hop"},{"key":"def","genre":"classic rock"}]`
	items, err := parseGenreBatchJSON(raw)
	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "abc123", items[0].Key)
	assert.Equal(t, "hip hop", items[0].Genre)
}

func TestParseGenreBatchJSON_markdownFence(t *testing.T) {
	raw := "```json\n[{\"key\":\"x\",\"genre\":\"pop\"}]\n```"
	items, err := parseGenreBatchJSON(raw)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "x", items[0].Key)
	assert.Equal(t, "pop", items[0].Genre)
}

func TestParseGenreBatchJSON_prefixNoise(t *testing.T) {
	raw := "Here is the JSON:\n[{\"key\":\"k1\",\"genre\":\"edm\"}]"
	items, err := parseGenreBatchJSON(raw)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "k1", items[0].Key)
}
