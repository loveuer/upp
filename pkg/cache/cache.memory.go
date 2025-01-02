package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitea.com/loveuer/gredis"
)

var _ Cache = (*_mem)(nil)

type _mem struct {
	client *gredis.Gredis
}

func (m *_mem) GetScan(ctx context.Context, key string) Scanner {
	return newScanner(m.Get(ctx, key))
}

func (m *_mem) GetExScan(ctx context.Context, key string, duration time.Duration) Scanner {
	return newScanner(m.GetEx(ctx, key, duration))
}

func (m *_mem) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := m.client.Get(key)
	if err != nil {
		if errors.Is(err, gredis.ErrKeyNotFound) {
			return nil, ErrorKeyNotFound
		}

		return nil, err
	}

	bs, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid value type=%T", v)
	}

	return bs, nil
}

func (m *_mem) Gets(ctx context.Context, keys ...string) ([][]byte, error) {
	bss := make([][]byte, 0, len(keys))
	for _, key := range keys {
		bs, err := m.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		bss = append(bss, bs)
	}

	return bss, nil
}

func (m *_mem) GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error) {
	v, err := m.client.GetEx(key, duration)
	if err != nil {
		if errors.Is(err, gredis.ErrKeyNotFound) {
			return nil, ErrorKeyNotFound
		}

		return nil, err
	}

	bs, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid value type=%T", v)
	}

	return bs, nil
}

func (m *_mem) Set(ctx context.Context, key string, value any) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}
	return m.client.Set(key, bs)
}

func (m *_mem) Sets(ctx context.Context, vm map[string]any) error {
	for k, v := range vm {
		if err := m.Set(ctx, k, v); err != nil {
			return err
		}
	}

	return nil
}

func (m *_mem) SetEx(ctx context.Context, key string, value any, duration time.Duration) error {
	bs, err := handleValue(value)
	if err != nil {
		return err
	}
	return m.client.SetEx(key, bs, duration)
}

func (m *_mem) Del(ctx context.Context, keys ...string) error {
	m.client.Delete(keys...)
	return nil
}

func (m *_mem) Close() error {
	m.client = nil

	return nil
}
