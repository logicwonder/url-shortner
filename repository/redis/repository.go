package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/logicwonder/url-shortner/shortner"
)

type redisRepository struct {
	client  *redis.Client
	timeout time.Duration
}

func newRedisClient(redisURL string, redisTimeout int) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(redisTimeout)*time.Second)
	defer cancel()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping(ctx).Result()

	return client, err
}

func NewRedisRepository(redisURL string, redisTimeout int) (shortner.RedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL, redisTimeout)
	if err != nil {
		return nil, fmt.Errorf("repository.NewRedisRepository: %w", err)
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortner.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	key := r.generateKey(code)
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("repository.Redirect.Find: ", shortner.ErrRedirectNotFound)
		}
		return nil, fmt.Errorf("repository.Redirect.Find: ", err)
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("repository.Redirect.Find: ", err)
	}
	redirect := &shortner.Redirect{}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil

}

func (r *redisRepository) Store(redirect *shortner.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(ctx, key, data).Result()
	if err != nil {
		return fmt.Errorf("repository.Redirect.Store: %w", err)
	}
	return nil
}
