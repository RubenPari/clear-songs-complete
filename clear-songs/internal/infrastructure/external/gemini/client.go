package gemini

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Redact Google API key in error strings (full request URL may appear in errors).
var googleAPIKeyQuery = regexp.MustCompile(`([?&])key=[A-Za-z0-9_-]+`)

func redactGoogleAPIKeyInString(msg string) string {
	return googleAPIKeyQuery.ReplaceAllString(msg, "${1}key=REDACTED")
}

// DefaultGeminiModel is used when GEMINI_MODEL is unset. gemini-2.0-flash is not
// available to new API users (404); see https://ai.google.dev/gemini-api/docs/models
const DefaultGeminiModel = "gemini-2.5-flash"

func geminiModelFromEnv() string {
	if m := strings.TrimSpace(os.Getenv("GEMINI_MODEL")); m != "" {
		return m
	}
	return DefaultGeminiModel
}

// GeminiRepository implements AIRepository using Google's Gemini API
type GeminiRepository struct {
	client *genai.Client
	model  string
}

// NewGeminiRepository creates a new Gemini repository
func NewGeminiRepository(ctx context.Context, apiKey string) (*GeminiRepository, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := geminiModelFromEnv()
	log.Printf("Gemini genre fallback: using model %q (override with GEMINI_MODEL)", model)

	return &GeminiRepository{
		client: client,
		model:  model,
	}, nil
}

// ResolveArtistGenre asks Gemini to determine the primary genre of an artist
func (r *GeminiRepository) ResolveArtistGenre(ctx context.Context, artistName string) (string, error) {
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
