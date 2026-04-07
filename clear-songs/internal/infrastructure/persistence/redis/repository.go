package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/redis/go-redis/v9"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	defaultTTL = 5 * time.Minute
	tokenTTL   = 24 * time.Hour
)

// RedisCacheRepository implements the CacheRepository interface using Redis
type RedisCacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCacheRepository creates a new Redis cache repository
func NewRedisCacheRepository() (*RedisCacheRepository, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "6379"
	}
	
	db := 0
	if dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &db)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	
	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis connection failed: %v", err)
		log.Println("WARNING: Application will continue without Redis caching")
		return nil, err
	}

	log.Println("Connected to Redis for caching")
	
	return &RedisCacheRepository{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// SetToken stores the OAuth token in cache
func (r *RedisCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error {
	if token == nil {
		return r.ClearToken(ctx)
	}
	return r.Set(ctx, "spotify_token", token, tokenTTL)
}

// GetToken retrieves the OAuth token from cache
func (r *RedisCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error) {
	var token oauth2.Token
	found, err := r.Get(ctx, "spotify_token", &token)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &token, nil
}

// ClearToken removes the token from cache
func (r *RedisCacheRepository) ClearToken(ctx context.Context) error {
	return r.Delete(ctx, "spotify_token")
}

// GetUserTracks retrieves cached user tracks
func (r *RedisCacheRepository) GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	var tracks []spotifyAPI.SavedTrack
	found, err := r.Get(ctx, "userTracks", &tracks)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return tracks, nil
}

// SetUserTracks stores user tracks in cache
func (r *RedisCacheRepository) SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error {
	return r.Set(ctx, "userTracks", tracks, ttl)
}

// InvalidateUserTracks removes user tracks from cache
func (r *RedisCacheRepository) InvalidateUserTracks(ctx context.Context) error {
	return r.Delete(ctx, "userTracks")
}

// GetPlaylistTracks retrieves cached playlist tracks
func (r *RedisCacheRepository) GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	key := "tracksPlaylist" + playlistID.String()
	var tracks []spotifyAPI.PlaylistTrack
	found, err := r.Get(ctx, key, &tracks)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return tracks, nil
}

// SetPlaylistTracks stores playlist tracks in cache
func (r *RedisCacheRepository) SetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	key := "tracksPlaylist" + playlistID.String()
	return r.Set(ctx, key, tracks, ttl)
}

// InvalidatePlaylistTracks removes playlist tracks from cache
func (r *RedisCacheRepository) InvalidatePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error {
	key := "tracksPlaylist" + playlistID.String()
	return r.Delete(ctx, key)
}

// Get retrieves a value from cache
func (r *RedisCacheRepository) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	if r.client == nil {
		return false, nil
	}
	
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	
	if err := json.Unmarshal(val, target); err != nil {
		return false, err
	}
	
	return true, nil
}

// Set stores a value in cache
func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if r.client == nil {
		return nil // Silently fail if Redis is not available
	}
	
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a key from cache
func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Del(ctx, key).Err()
}

// Ensure RedisCacheRepository implements CacheRepository interface
var _ shared.CacheRepository = (*RedisCacheRepository)(nil)
