package cache

import (
	"fmt"
	"net/url"
	"strings"

	"gitea.com/loveuer/gredis"
	"github.com/go-redis/redis/v8"
	"github.com/loveuer/upp/pkg/tool"
)

func New(uri string) (Cache, error) {
	var (
		client Cache
		err    error
	)

	strs := strings.Split(uri, "::")

	switch strs[0] {
	case "memory":
		gc := gredis.NewGredis(1024 * 1024)
		client = &_mem{client: gc}
	case "lru":
		if client, err = newLRUCache(); err != nil {
			return nil, err
		}
	case "redis":
		var (
			ins *url.URL
			err error
		)

		if len(strs) != 2 {
			return nil, fmt.Errorf("cache.Init: invalid cache uri: %s", uri)
		}

		uri := strs[1]

		if !strings.Contains(uri, "://") {
			uri = fmt.Sprintf("redis://%s", uri)
		}

		if ins, err = url.Parse(uri); err != nil {
			return nil, fmt.Errorf("cache.Init: url parse cache uri: %s, err: %s", uri, err.Error())
		}

		addr := ins.Host
		username := ins.User.Username()
		password, _ := ins.User.Password()

		var rc *redis.Client
		rc = redis.NewClient(&redis.Options{
			Addr:     addr,
			Username: username,
			Password: password,
		})

		if err = rc.Ping(tool.Timeout(5)).Err(); err != nil {
			return nil, fmt.Errorf("cache.Init: redis ping err: %s", err.Error())
		}

		client = &_redis{client: rc}
	default:
		return nil, fmt.Errorf("cache type %s not support", strs[0])
	}

	return client, nil
}
