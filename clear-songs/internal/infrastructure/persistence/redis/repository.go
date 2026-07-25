// Package redis implements the CacheRepository interface using Redis as the
// backing store. It provides JSON-based get/set/delete operations, OAuth token
// management, user-tracks caching, and playlist-tracks caching. When Redis is
// unavailable a NoOp fallback is used instead.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RubenPari/clear-songs/internal/domain/shared"
	"github.com/redis/go-redis/v9"
	spotifyAPI "github.com/zmb3/spotify"
	"go.uber.org/zap"
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

// NewRedisCacheRepository creates a Redis client from environment variables.
// It first checks for REDIS_URL (a full connection string), then falls back to
// individual REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, and REDIS_DB variables.
// Returns an error if the Redis server cannot be reached after retries.
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
		zap.L().Warn("Redis connection failed", zap.Error(err))
		return nil, err
	}
	zap.L().Info("connected to Redis for caching")
	return &RedisCacheRepository{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// pingRedisWithRetry attempts to ping Redis up to 3 times with linear backoff
// (100ms, 200ms, 300ms). Returns the last error if all attempts fail.
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

// newClientFromRedisURL creates a Redis client from a full connection URL.
// If the initial connection with a redis:// scheme fails, it automatically
// retries with rediss:// (TLS) to support providers like Upstash that require
// TLS even when the URL does not explicitly specify it.
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
		zap.L().Info("connected to Redis for caching")
		return rdb, nil
	}

	// redis:// against a TLS-only endpoint (e.g. Upstash) often fails with EOF.
	if strings.HasPrefix(redisURL, "redis://") {
		alt := "rediss://" + strings.TrimPrefix(redisURL, "redis://")
		zap.L().Warn("Redis REDIS_URL ping failed, retrying with TLS", zap.Error(err))
		rdb2, err2 := dial(alt)
		if err2 == nil {
			zap.L().Info("connected to Redis for caching")
			return rdb2, nil
		}
		return nil, fmt.Errorf("redis: %w; rediss fallback: %w", err, err2)
	}

	zap.L().Warn("Redis connection failed", zap.Error(err))
	return nil, err
}

// SetToken stores the OAuth2 token in Redis with a 24-hour TTL. Passing nil
// clears the token instead.
func (r *RedisCacheRepository) SetToken(ctx context.Context, token *oauth2.Token) error {
	if token == nil {
		return r.ClearToken(ctx)
	}
	return r.Set(ctx, shared.CacheKeyTokenPrefix, token, tokenTTL)
}

// GetToken retrieves the cached OAuth2 token from Redis. Returns (nil, nil)
// when no token is stored.
func (r *RedisCacheRepository) GetToken(ctx context.Context) (*oauth2.Token, error) {
	var token oauth2.Token
	found, err := r.Get(ctx, shared.CacheKeyTokenPrefix, &token)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &token, nil
}

// ClearToken removes the cached OAuth2 token from Redis.
func (r *RedisCacheRepository) ClearToken(ctx context.Context) error {
	return r.Delete(ctx, shared.CacheKeyTokenPrefix)
}

// GetUserTracks retrieves the cached saved tracks from Redis. Returns (nil, nil)
// when no tracks are cached.
func (r *RedisCacheRepository) GetUserTracks(ctx context.Context) ([]spotifyAPI.SavedTrack, error) {
	var tracks []spotifyAPI.SavedTrack
	found, err := r.Get(ctx, shared.CacheKeyUserTracks, &tracks)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return tracks, nil
}

// SetUserTracks stores the user's saved tracks in Redis with the given TTL.
func (r *RedisCacheRepository) SetUserTracks(ctx context.Context, tracks []spotifyAPI.SavedTrack, ttl time.Duration) error {
	return r.Set(ctx, shared.CacheKeyUserTracks, tracks, ttl)
}

// InvalidateUserTracks removes the cached user tracks and all track-summary
// cache entries (matched via SCAN pattern) from Redis.
func (r *RedisCacheRepository) InvalidateUserTracks(ctx context.Context) error {
	if err := r.Delete(ctx, shared.CacheKeyUserTracks); err != nil {
		return err
	}
	return r.deleteKeysByPattern(ctx, shared.CacheKeyTrackSummaryFmt+"*")
}

// deleteKeysByPattern uses SCAN to iteratively find and delete all keys matching
// the given pattern. It processes keys in batches of 100 to avoid blocking Redis.
func (r *RedisCacheRepository) deleteKeysByPattern(ctx context.Context, pattern string) error {
	if r.client == nil {
		return nil
	}
	var cursor uint64
	for {
		keys, next, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

// GetPlaylistTracks retrieves cached playlist tracks from Redis. Returns
// (nil, nil) when no tracks are cached for the given playlist.
func (r *RedisCacheRepository) GetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	key := playlistTracksCacheKey(playlistID)
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

// SetPlaylistTracks stores playlist tracks in Redis with the given TTL.
func (r *RedisCacheRepository) SetPlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID, tracks []spotifyAPI.PlaylistTrack, ttl time.Duration) error {
	key := playlistTracksCacheKey(playlistID)
	return r.Set(ctx, key, tracks, ttl)
}

// InvalidatePlaylistTracks removes the cached tracks for a specific playlist.
func (r *RedisCacheRepository) InvalidatePlaylistTracks(ctx context.Context, playlistID spotifyAPI.ID) error {
	key := playlistTracksCacheKey(playlistID)
	return r.Delete(ctx, key)
}

// playlistTracksCacheKey builds the Redis key for a specific playlist's tracks.
func playlistTracksCacheKey(playlistID spotifyAPI.ID) string {
	return fmt.Sprintf(shared.CacheKeyPlaylistTracksFmt, playlistID.String())
}

// Get retrieves a JSON-encoded value from Redis and unmarshals it into target.
// Returns (false, nil) on cache miss. Returns (false, nil) when the client is nil.
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

// Set marshals value to JSON and stores it in Redis with the given TTL.
// Silently returns nil when the client is nil.
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

// Delete removes a single key from Redis. Silently returns nil when the client is nil.
func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Del(ctx, key).Err()
}

// Ensure RedisCacheRepository implements CacheRepository interface
var _ shared.CacheRepository = (*RedisCacheRepository)(nil)
