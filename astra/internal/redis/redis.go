package redis

import (
	"context"
	"github.com/andibalo/ramein/astra/internal/config"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var instance *redis.Client

// Redis driver
type Redis struct {
	sync.RWMutex
	Client   *redis.Client
	Addr     string
	Password string
	DB       int
}

// NewRedisDriver create a new instance
func NewRedisDriver(cfg config.Config) *Redis {
	return &Redis{
		Addr:     cfg.RedisURL(),
		Password: cfg.RedisPassword(),
		DB:       cfg.RedisDB(),
	}
}

// Connect establish a redis connection
func (r *Redis) Connect(ctx context.Context) (bool, error) {
	r.Lock()
	defer r.Unlock()

	// Reuse redis connections
	if instance == nil {
		r.Client = redis.NewClient(&redis.Options{
			Addr:        r.Addr,
			Password:    r.Password,
			DB:          r.DB,
			PoolSize:    10,
			PoolTimeout: 30 * time.Second,
		})

		instance = r.Client

		_, err := r.Ping(ctx)

		if err != nil {
			return false, err
		}
	} else {
		r.Client = instance

		_, err := r.Ping(ctx)

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Ping checks the redis connection
func (r *Redis) Ping(ctx context.Context) (bool, error) {
	pong, err := r.Client.Ping(ctx).Result()

	if err != nil {
		return false, err
	}
	return pong == "PONG", nil
}

// Set sets a record
func (r *Redis) Set(ctx context.Context, key, value string, expiration time.Duration) (bool, error) {
	result := r.Client.Set(ctx, key, value, expiration)

	if result.Err() != nil {
		return false, result.Err()
	}

	return result.Val() == "OK", nil
}

// Get gets a record value
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	result := r.Client.Get(ctx, key)

	if result.Err() != nil {
		return "", result.Err()
	}

	return result.Val(), nil
}

// Exists deletes a record
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result := r.Client.Exists(ctx, key)

	if result.Err() != nil {
		return false, result.Err()
	}

	return result.Val() > 0, nil
}

// Del deletes a record
func (r *Redis) Del(ctx context.Context, key string) (int64, error) {
	result := r.Client.Del(ctx, key)

	if result.Err() != nil {
		return 0, result.Err()
	}

	return result.Val(), nil
}

// HGet gets a record from hash
func (r *Redis) HGet(ctx context.Context, key, field string) (string, error) {
	result := r.Client.HGet(ctx, key, field)

	if result.Err() != nil {
		return "", result.Err()
	}

	return result.Val(), nil
}

// HSet sets a record in hash
func (r *Redis) HSet(ctx context.Context, key, field, value string) (int64, error) {
	result := r.Client.HSet(ctx, key, field, value)

	if result.Err() != nil {
		return 0, result.Err()
	}

	return result.Val(), nil
}

// HExists checks if key exists on a hash
func (r *Redis) HExists(ctx context.Context, key, field string) (bool, error) {
	result := r.Client.HExists(ctx, key, field)

	if result.Err() != nil {
		return false, result.Err()
	}

	return result.Val(), nil
}

// HDel deletes a hash record
func (r *Redis) HDel(ctx context.Context, key, field string) (int64, error) {
	result := r.Client.HDel(ctx, key, field)

	if result.Err() != nil {
		return 0, result.Err()
	}

	return result.Val(), nil
}

// HLen count hash records
func (r *Redis) HLen(ctx context.Context, key string) (int64, error) {
	result := r.Client.HLen(ctx, key)

	if result.Err() != nil {
		return 0, result.Err()
	}

	return result.Val(), nil
}

// HTruncate deletes a hash
func (r *Redis) HTruncate(ctx context.Context, key string) (int64, error) {
	result := r.Client.Del(ctx, key)

	if result.Err() != nil {
		return 0, result.Err()
	}

	return result.Val(), nil
}

// HScan return an iterative obj for a hash
func (r *Redis) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.Client.HScan(ctx, key, cursor, match, count)
}
