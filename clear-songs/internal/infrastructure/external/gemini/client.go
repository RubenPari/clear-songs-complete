package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
)

// Redact Google API key in error strings (full request URL may appear in errors).
var googleAPIKeyQuery = regexp.MustCompile(`([?&])key=[A-Za-z0-9_-]+`)

// Redact google apikey in string.
func redactGoogleAPIKeyInString(msg string) string {
	return googleAPIKeyQuery.ReplaceAllString(msg, "${1}key=REDACTED")
}

// DefaultGeminiModel is used when GEMINI_MODEL is unset. gemini-2.0-flash is not
// available to new API users (404); see https://ai.google.dev/gemini-api/docs/models
const DefaultGeminiModel = "gemini-2.5-flash"

// Gemini model from env.
func geminiModelFromEnv() string {
	if m := strings.TrimSpace(os.Getenv("GEMINI_MODEL")); m != "" {
		return m
	}
	return DefaultGeminiModel
}

// Gemini genre batch size.
func geminiGenreBatchSize() int {
	const defaultN = 24
	s := strings.TrimSpace(os.Getenv("GEMINI_GENRE_BATCH_SIZE"))
	if s == "" {
		return defaultN
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 100 {
		return defaultN
	}
	return n
}

// Gemini genre batch parallel.
func geminiGenreBatchParallel() int64 {
	const defaultP = 2
	s := strings.TrimSpace(os.Getenv("GEMINI_GENRE_BATCH_PARALLEL"))
	if s == "" {
		return defaultP
	}
	p, err := strconv.Atoi(s)
	if err != nil || p < 1 || p > 16 {
		return defaultP
	}
	return int64(p)
}

// Gemini request timeout sec.
func geminiRequestTimeoutSec() int {
	const defaultSec = 25
	s := strings.TrimSpace(os.Getenv("GEMINI_REQUEST_TIMEOUT_SEC"))
	if s == "" {
		return defaultSec
	}
	sec, err := strconv.Atoi(s)
	if err != nil || sec < 5 || sec > 120 {
		return defaultSec
	}
	return sec
}

// GeminiRepository implements AIRepository using Google's Gemini API
type GeminiRepository struct {
	client *genai.Client
	model  string
}

// Creates gemini repository.
func NewGeminiRepository(ctx context.Context, apiKey string) (*GeminiRepository, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := geminiModelFromEnv()
	zap.L().Info("Gemini genre fallback configured", zap.String("model", model))

	return &GeminiRepository{
		client: client,
		model:  model,
	}, nil
}

// Resolves artist genres.
func (r *GeminiRepository) ResolveArtistGenres(ctx context.Context, lookups []shared.AIArtistLookup) (map[string]string, error) {
	if len(lookups) == 0 {
		return map[string]string{}, nil
	}

	sec := geminiRequestTimeoutSec()
	chunkSize := geminiGenreBatchSize()
	parallel := geminiGenreBatchParallel()

	chunks := chunkLookups(lookups, chunkSize)
	merged := make(map[string]string, len(lookups))
	var mu sync.Mutex

	sem := semaphore.NewWeighted(parallel)
	g, gctx := errgroup.WithContext(ctx)

	for _, chunk := range chunks {
		chunk := chunk
		g.Go(func() error {
			if err := sem.Acquire(gctx, 1); err != nil {
				return err
			}
			defer sem.Release(1)

			cctx, cancel := context.WithTimeout(gctx, time.Duration(sec)*time.Second)
			defer cancel()

			part := r.resolveChunkWithFallback(cctx, chunk)
			mu.Lock()
			for k, v := range part {
				merged[k] = v
			}
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return merged, err
	}
	return merged, nil
}

// Chunk lookups.
func chunkLookups(lookups []shared.AIArtistLookup, size int) [][]shared.AIArtistLookup {
	if size < 1 {
		size = 24
	}
	var out [][]shared.AIArtistLookup
	for i := 0; i < len(lookups); i += size {
		j := i + size
		if j > len(lookups) {
			j = len(lookups)
		}
		out = append(out, lookups[i:j:j])
	}
	return out
}

// Resolves chunk with fallback.
func (r *GeminiRepository) resolveChunkWithFallback(ctx context.Context, chunk []shared.AIArtistLookup) map[string]string {
	if len(chunk) == 1 {
		return r.resolveOne(ctx, chunk[0])
	}

	out, err := r.resolveChunkBatch(ctx, chunk)
	if err != nil {
		zap.L().Warn("Gemini batch failed, trying single calls",
			zap.Int("artist_count", len(chunk)),
			zap.Error(err),
		)
		return r.resolveChunkSingles(ctx, chunk)
	}
	if out == nil {
		out = make(map[string]string)
	}
	for _, l := range chunk {
		g, ok := out[l.Key]
		if !ok || strings.TrimSpace(g) == "" {
			g2, e := r.resolveArtistGenreSingle(ctx, l.Name)
			if e != nil {
				zap.L().Warn("Gemini batch incomplete item", zap.String("key", l.Key), zap.Error(e))
				out[l.Key] = ""
				continue
			}
			out[l.Key] = g2
		}
	}
	return out
}

// Resolves one.
func (r *GeminiRepository) resolveOne(ctx context.Context, l shared.AIArtistLookup) map[string]string {
	m := make(map[string]string, 1)
	g, err := r.resolveArtistGenreSingle(ctx, l.Name)
	if err != nil {
		zap.L().Warn("Gemini resolve artist genre failed",
			zap.String("key", l.Key),
			zap.String("artist", l.Name),
			zap.Error(err),
		)
		m[l.Key] = ""
		return m
	}
	m[l.Key] = g
	return m
}

// Resolves chunk singles.
func (r *GeminiRepository) resolveChunkSingles(ctx context.Context, chunk []shared.AIArtistLookup) map[string]string {
	out := make(map[string]string, len(chunk))
	for _, l := range chunk {
		g, err := r.resolveArtistGenreSingle(ctx, l.Name)
		if err != nil {
			zap.L().Warn("Gemini single fallback failed", zap.String("key", l.Key), zap.Error(err))
			out[l.Key] = ""
			continue
		}
		out[l.Key] = g
	}
	return out
}

type genreBatchItem struct {
	Key   string `json:"key"`
	Genre string `json:"genre"`
}

// Resolves chunk batch.
func (r *GeminiRepository) resolveChunkBatch(ctx context.Context, chunk []shared.AIArtistLookup) (map[string]string, error) {
	var sb strings.Builder
	for _, l := range chunk {
		k := strings.ReplaceAll(l.Key, "\t", " ")
		n := strings.ReplaceAll(l.Name, "\t", " ")
		sb.WriteString(k)
		sb.WriteByte('\t')
		sb.WriteString(n)
		sb.WriteByte('\n')
	}

	prompt := fmt.Sprintf(`Reply with ONLY a JSON array (no markdown fences, no commentary). Each element must be an object with string fields "key" and "genre".
Use the exact "key" values from the input. Infer one primary music genre per artist (short phrase).

Input lines (key TAB display name):
%s`, sb.String())

	model := r.client.GenerativeModel(r.model)
	model.SetTemperature(0)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("Gemini batch API: %s", redactGoogleAPIKeyInString(err.Error()))
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty batch response from Gemini")
	}

	raw := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	items, err := parseGenreBatchJSON(raw)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(items))
	for _, it := range items {
		k := strings.TrimSpace(it.Key)
		g := strings.TrimSpace(strings.ToLower(it.Genre))
		if k != "" {
			out[k] = g
		}
	}
	return out, nil
}

// Parse genre batch json.
func parseGenreBatchJSON(raw string) ([]genreBatchItem, error) {
	s := strings.TrimSpace(raw)
	s = stripMarkdownJSONFence(s)
	idx := strings.Index(s, "[")
	if idx < 0 {
		return nil, fmt.Errorf("no JSON array in Gemini batch response")
	}
	end := strings.LastIndex(s, "]")
	if end <= idx {
		return nil, fmt.Errorf("invalid JSON array in Gemini batch response")
	}
	s = s[idx : end+1]

	var items []genreBatchItem
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, fmt.Errorf("parse batch JSON: %w", err)
	}
	return items, nil
}

var markdownJSONFence = regexp.MustCompile("(?s)^```(?:json)?\\s*")

// Strip markdown jsonfence.
func stripMarkdownJSONFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = markdownJSONFence.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "```"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return s
}

// Resolves artist genre single.
func (r *GeminiRepository) resolveArtistGenreSingle(ctx context.Context, artistName string) (string, error) {
	model := r.client.GenerativeModel(r.model)
	model.SetTemperature(0)

	prompt := fmt.Sprintf(
		`What is the primary music genre of the artist "%s"? Reply with only the genre name, nothing else.`,
		artistName,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Gemini API error for artist %s: %s", artistName, redactGoogleAPIKeyInString(err.Error()))
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini for artist %s", artistName)
	}

	result := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	return strings.TrimSpace(strings.ToLower(result)), nil
}

// Ensure GeminiRepository implements AIRepository
var _ shared.AIRepository = (*GeminiRepository)(nil)
