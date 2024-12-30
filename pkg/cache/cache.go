package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Gets(ctx context.Context, keys ...string) ([][]byte, error)
	GetScan(ctx context.Context, key string) Scanner
	GetEx(ctx context.Context, key string, duration time.Duration) ([]byte, error)
	GetExScan(ctx context.Context, key string, duration time.Duration) Scanner
	// Set value 会被序列化, 优先使用 MarshalBinary 方法, 没有则执行 json.Marshal
	Set(ctx context.Context, key string, value any) error
	Sets(ctx context.Context, vm map[string]any) error
	// SetEx value 会被序列化, 优先使用 MarshalBinary 方法, 没有则执行 json.Marshal
	SetEx(ctx context.Context, key string, value any, duration time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type Scanner interface {
	Scan(model any) error
}

type encoded_value interface {
	MarshalBinary() ([]byte, error)
}

type decoded_value interface {
	UnmarshalBinary(bs []byte) error
}

const (
	Prefix = "upp:"
)

var ErrorKeyNotFound = errors.New("key not found")

func handleValue(value any) ([]byte, error) {
	var (
		bs  []byte
		err error
	)

	switch value.(type) {
	case []byte:
		return value.([]byte), nil
	}

	if imp, ok := value.(encoded_value); ok {
		bs, err = imp.MarshalBinary()
	} else {
		bs, err = json.Marshal(value)
	}

	return bs, err
}
