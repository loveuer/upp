package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type _redis struct {
	client *redis.Client
}

func (r *_redis) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorKeyNotFound
		}

		return nil, err
	}

	return []byte(result), nil
}

func (r *_redis) Gets(ctx context.Context, keys ...string) ([][]byte, error) {
	result, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorKeyNotFound
		}

		return nil, err
	}

	return lo.Map(
		result,
		func(item any, index int) []byte {
			return []byte(cast.ToString(item))
		},
	), nil
}

func (r *_redis) GetScan(ctx context.Context, key string) Scanner {
	return newScanner(r.Get(ctx, key))
}

func (r *_redis) GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error) {
	result, err := r.client.GetEx(ctx, key, duration).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrorKeyNotFound
		}

		return nil, err
	}

	return []byte(result), nil
}

func (r *_redis) GetExScan(ctx context.Context, key string, duration time.Duration) Scanner {
	return newScanner(r.GetEx(ctx, key, duration))
}

func (r *_redis) Set(ctx context.Context, key string, value any) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	_, err = r.client.Set(ctx, key, bs, redis.KeepTTL).Result()
	return err
}

func (r *_redis) Sets(ctx context.Context, values map[string]any) error {
	vm := make(map[string]any)
	for k, v := range values {
		bs, err := handleValue(v)
		if err != nil {
			return err
		}

		vm[k] = bs
	}

	return r.client.MSet(ctx, vm).Err()
}

func (r *_redis) SetEx(ctx context.Context, key string, value any, duration time.Duration) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}

	_, err = r.client.SetEX(ctx, key, bs, duration).Result()

	return err
}

func (r *_redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
