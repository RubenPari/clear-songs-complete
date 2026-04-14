package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/internal/application/auth"
	"github.com/RubenPari/clear-songs/internal/application/playlist"
	"github.com/RubenPari/clear-songs/internal/application/track"
	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/RubenPari/clear-songs/internal/domain/shared/constants"
	"github.com/RubenPari/clear-songs/internal/infrastructure/external/gemini"
	"github.com/RubenPari/clear-songs/internal/infrastructure/external/spotify"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/postgres"
	"github.com/RubenPari/clear-songs/internal/infrastructure/persistence/redis"
	spotifyAPI "github.com/zmb3/spotify"
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
	LoginUC    *auth.LoginUseCase
	CallbackUC *auth.CallbackUseCase
	LogoutUC   *auth.LogoutUseCase
	IsAuthUC   *auth.IsAuthUseCase

	// Track Use Cases
	GetTrackSummaryUseCase *track.GetTrackSummaryUseCase
	DeleteTracksByArtistUC *track.DeleteTracksByArtistUseCase
	DeleteTracksByRangeUC  *track.DeleteTracksByRangeUseCase
	DeleteTrackUC          *track.DeleteTrackUseCase
	GetTracksByArtistUC    *track.GetTracksByArtistUseCase

	// Playlist Use Cases
	GetUserPlaylistsUC         *playlist.GetUserPlaylistsUseCase
	DeletePlaylistTracksUC     *playlist.DeletePlaylistTracksUseCase
	DeletePlaylistAndLibraryUC *playlist.DeletePlaylistAndLibraryTracksUseCase
}

// NewContainer creates and initializes a new dependency injection container
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
			log.Printf("WARNING: Gemini initialization failed: %v", err)
			aiRepo = gemini.NewNoOpAIRepository()
		} else {
			aiRepo = geminiRepo
		}
	} else {
		log.Println("WARNING: GEMINI_API_KEY not set, AI genre resolution disabled")
		aiRepo = gemini.NewNoOpAIRepository()
	}

	// Initialize auth use cases
	loginUC := auth.NewLoginUseCase(oauthConfig)
	callbackUC := auth.NewCallbackUseCase(oauthConfig, spotifyRepo, cacheRepo)
	logoutUC := auth.NewLogoutUseCase(spotifyRepo, cacheRepo)
	isAuthUC := auth.NewIsAuthUseCase(spotifyRepo)

	// Initialize track use cases
	getTrackSummaryUseCase := track.NewGetTrackSummaryUseCase(spotifyRepo, cacheRepo, aiRepo)
	deleteTracksByArtistUC := track.NewDeleteTracksByArtistUseCase(spotifyRepo, cacheRepo)
	getTracksByArtistUC := track.NewGetTracksByArtistUseCase(spotifyRepo, cacheRepo)
	deleteTrackUC := track.NewDeleteTrackUseCase(spotifyRepo, cacheRepo, databaseRepo)
	deleteTracksByRangeUC := track.NewDeleteTracksByRangeUseCase(
		spotifyRepo,
		cacheRepo,
		getTrackSummaryUseCase,
		deleteTracksByArtistUC,
	)

	// Initialize playlist use cases
	getUserPlaylistsUC := playlist.NewGetUserPlaylistsUseCase(spotifyRepo, cacheRepo)
	deletePlaylistTracksUC := playlist.NewDeletePlaylistTracksUseCase(spotifyRepo, cacheRepo)
	deletePlaylistAndLibraryUC := playlist.NewDeletePlaylistAndLibraryTracksUseCase(
		spotifyRepo,
		cacheRepo,
		databaseRepo,
		deletePlaylistTracksUC,
	)

	container := &Container{
		SpotifyRepo:                spotifyRepo,
		CacheRepo:                  cacheRepo,
		DatabaseRepo:               databaseRepo,
		AIRepo:                     aiRepo,
		OAuthConfig:                oauthConfig,
		LoginUC:                    loginUC,
		CallbackUC:                 callbackUC,
		LogoutUC:                   logoutUC,
		IsAuthUC:                   isAuthUC,
		GetTrackSummaryUseCase:     getTrackSummaryUseCase,
		DeleteTracksByArtistUC:     deleteTracksByArtistUC,
		DeleteTracksByRangeUC:      deleteTracksByRangeUC,
		DeleteTrackUC:              deleteTrackUC,
		GetTracksByArtistUC:        getTracksByArtistUC,
		GetUserPlaylistsUC:         getUserPlaylistsUC,
		DeletePlaylistTracksUC:     deletePlaylistTracksUC,
		DeletePlaylistAndLibraryUC: deletePlaylistAndLibraryUC,
	}

	return container, nil
}

// GetOAuth2Config returns OAuth2 configuration from environment variables
func GetOAuth2Config() (*oauth2.Config, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("REDIRECT_URL")
	if redirectURI == "" {
		redirectURI = os.Getenv("REDIRECT_URI")
	}

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, errors.New("missing required environment variables: CLIENT_ID, CLIENT_SECRET, REDIRECT_URL")
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
