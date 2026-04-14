package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

// NewRedisCacheRepository creates a new Redis cache repository.
// If REDIS_URL is set (e.g. redis:// or rediss:// from Fly/Upstash), it is used via redis.ParseURL.
// Otherwise REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB are used (defaults: 127.0.0.1:6379).
func NewRedisCacheRepository() (*RedisCacheRepository, error) {
	ctx := context.Background()

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		rdb, err := newClientFromRedisURL(ctx, redisURL)
		if err != nil {
			return nil, err
		}
		return &RedisCacheRepository{
			client: rdb,
			ctx:    ctx,
		}, nil
	}

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
	if err := pingRedisWithRetry(ctx, rdb); err != nil {
		_ = rdb.Close()
		log.Printf("WARNING: Redis connection failed: %v", err)
		return nil, err
	}
	log.Println("Connected to Redis for caching")
	return &RedisCacheRepository{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func pingRedisWithRetry(ctx context.Context, rdb *redis.Client) error {
	const attempts = 3
	var lastErr error
	for i := 0; i < attempts; i++ {
		if i > 0 {
			time.Sleep(time.Duration(100*i) * time.Millisecond)
		}
		lastErr = rdb.Ping(ctx).Err()
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}

func newClientFromRedisURL(ctx context.Context, redisURL string) (*redis.Client, error) {
	dial := func(url string) (*redis.Client, error) {
		opt, err := redis.ParseURL(url)
		if err != nil {
			return nil, fmt.Errorf("parse REDIS_URL: %w", err)
		}
		rdb := redis.NewClient(opt)
		if err := pingRedisWithRetry(ctx, rdb); err != nil {
			_ = rdb.Close()
			return nil, err
		}
		return rdb, nil
	}

	rdb, err := dial(redisURL)
	if err == nil {
		log.Println("Connected to Redis for caching")
		return rdb, nil
	}

	// redis:// against a TLS-only endpoint (e.g. Upstash) often fails with EOF.
	if strings.HasPrefix(redisURL, "redis://") {
		alt := "rediss://" + strings.TrimPrefix(redisURL, "redis://")
		log.Printf("WARNING: Redis REDIS_URL ping failed (%v); retrying with TLS (rediss://)", err)
		rdb2, err2 := dial(alt)
		if err2 == nil {
			log.Println("Connected to Redis for caching")
			return rdb2, nil
		}
		return nil, fmt.Errorf("redis: %w; rediss fallback: %w", err, err2)
	}

	log.Printf("WARNING: Redis connection failed: %v", err)
	return nil, err
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
