package di

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/RubenPari/clear-songs/internal/application/playlist"
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/shared/constants"
	"github.com/RubenPari/clear-songs/internal/infrastructure/config"
	"github.com/RubenPari/clear-songs/internal/infrastructure/external/gemini"
	"github.com/RubenPari/clear-songs/internal/infrastructure/external/spotify"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/redis"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// Container holds all application dependencies
type Container struct {
	// Repositories (as interfaces)
	SpotifyRepo  shared.SpotifyRepository
	CacheRepo    shared.CacheRepository
	DatabaseRepo shared.DatabaseRepository
	AIRepo       shared.AIRepository

	// OAuth Config
	OAuthConfig *oauth2.Config

	// Auth Use Cases
	LoginUseCase    *auth.LoginUseCase
	CallbackUseCase *auth.CallbackUseCase
	LogoutUseCase   *auth.LogoutUseCase
	IsAuthUseCase   *auth.IsAuthUseCase

	// Track Use Cases
	GetTrackSummaryUseCase *track.GetTrackSummaryUseCase
	DeleteTracksByArtistUseCase *track.DeleteTracksByArtistUseCase
	DeleteTracksByRangeUseCase  *track.DeleteTracksByRangeUseCase
	DeleteTrackUseCase          *track.DeleteTrackUseCase
	GetTracksByArtistUseCase    *track.GetTracksByArtistUseCase
	InvalidateLibraryCacheUseCase *track.InvalidateLibraryCacheUseCase

	// Playlist Use Cases
	GetUserPlaylistsUseCase         *playlist.GetUserPlaylistsUseCase
	DeletePlaylistTracksUseCase     *playlist.DeletePlaylistTracksUseCase
	DeletePlaylistAndLibraryUseCase *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// Creates container.
func NewContainer() (*Container, error) {
	// Initialize OAuth config
	oauthConfig, err := GetOAuth2Config()
	if err != nil {
		return nil, err
	}

	// Initialize Spotify repository
	spotifyRepo := spotify.NewSpotifyRepository(oauthConfig.ClientID, oauthConfig.ClientSecret, oauthConfig.RedirectURL, constants.Scopes)

	// Redis is required: OAuth token and caching depend on it.
	redisCache, err := redis.NewRedisCacheRepository()
	if err != nil {
		return nil, fmt.Errorf("redis required: %w", err)
	}
	cacheRepo := redisCache

	// Initialize database repository (may be nil if database not available)
	databaseRepo := postgres.NewPostgresRepository(postgres.Db)

	// Initialize AI repository (for genre resolution fallback)
	var aiRepo shared.AIRepository
	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey != "" {
		geminiRepo, err := gemini.NewGeminiRepository(context.Background(), geminiKey)
		if err != nil {
			zap.L().Warn("Gemini initialization failed", zap.Error(err))
			aiRepo = gemini.NewNoOpAIRepository()
		} else {
			aiRepo = geminiRepo
		}
	} else {
		zap.L().Warn("GEMINI_API_KEY not set, AI genre resolution disabled")
		aiRepo = gemini.NewNoOpAIRepository()
	}

	// Initialize auth use cases
	loginUseCase := auth.NewLoginUseCase(oauthConfig)
	callbackUseCase := auth.NewCallbackUseCase(oauthConfig, spotifyRepo, cacheRepo)
	logoutUseCase := auth.NewLogoutUseCase(spotifyRepo, cacheRepo)
	isAuthUseCase := auth.NewIsAuthUseCase(spotifyRepo)

	// Initialize track use cases
	getTrackSummaryUseCase := track.NewGetTrackSummaryUseCase(spotifyRepo, cacheRepo, aiRepo)
	deleteTracksByArtistUseCase := track.NewDeleteTracksByArtistUseCase(spotifyRepo, cacheRepo)
	getTracksByArtistUseCase := track.NewGetTracksByArtistUseCase(spotifyRepo, cacheRepo)
	deleteTrackUseCase := track.NewDeleteTrackUseCase(spotifyRepo, cacheRepo, databaseRepo)
	deleteTracksByRangeUseCase := track.NewDeleteTracksByRangeUseCase(spotifyRepo, cacheRepo)
	invalidateLibraryCacheUseCase := track.NewInvalidateLibraryCacheUseCase(cacheRepo)

	// Initialize playlist use cases
	getUserPlaylistsUseCase := playlist.NewGetUserPlaylistsUseCase(spotifyRepo, cacheRepo)
	deletePlaylistTracksUseCase := playlist.NewDeletePlaylistTracksUseCase(spotifyRepo, cacheRepo)
	deletePlaylistAndLibraryUseCase := playlist.NewDeletePlaylistAndLibraryTracksUseCase(
		spotifyRepo,
		cacheRepo,
		databaseRepo,
	)

	container := &Container{
		SpotifyRepo:                spotifyRepo,
		CacheRepo:                  cacheRepo,
		DatabaseRepo:               databaseRepo,
		AIRepo:                     aiRepo,
		OAuthConfig:                oauthConfig,
		LoginUseCase:               loginUseCase,
		CallbackUseCase:            callbackUseCase,
		LogoutUseCase:              logoutUseCase,
		IsAuthUseCase:              isAuthUseCase,
		GetTrackSummaryUseCase:     getTrackSummaryUseCase,
		DeleteTracksByArtistUseCase: deleteTracksByArtistUseCase,
		DeleteTracksByRangeUseCase:  deleteTracksByRangeUseCase,
		DeleteTrackUseCase:          deleteTrackUseCase,
		GetTracksByArtistUseCase:    getTracksByArtistUseCase,
		InvalidateLibraryCacheUseCase: invalidateLibraryCacheUseCase,
		GetUserPlaylistsUseCase:         getUserPlaylistsUseCase,
		DeletePlaylistTracksUseCase:     deletePlaylistTracksUseCase,
		DeletePlaylistAndLibraryUseCase: deletePlaylistAndLibraryUseCase,
	}

	return container, nil
}

// Fetches oauth2 config.
func GetOAuth2Config() (*oauth2.Config, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := config.GetRedirectURL()

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, errors.New("missing required environment variables: CLIENT_ID, CLIENT_SECRET, REDIRECT_URL or REDIRECT_URI")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       constants.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  spotifyAPI.AuthURL,
			TokenURL: spotifyAPI.TokenURL,
		},
	}, nil
}
