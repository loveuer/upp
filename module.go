package upp

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/loveuer/upp/pkg/api"
	"github.com/loveuer/upp/pkg/cache"
	"github.com/loveuer/upp/pkg/db"
	"github.com/loveuer/upp/pkg/es"
	"github.com/loveuer/upp/pkg/interfaces"
)

type module func(u *upp)

func InitDB(uri string, models ...any) module {
	db, err := db.New(uri)
	if err != nil {
		log.Panic(err.Error())
	}

	if err = db.AutoMigrate(models...); err != nil {
		log.Panic(err.Error())
	}

	return func(u *upp) {
		u.db = db
	}
}

func InitCache(uri string) module {
	cache, err := cache.New(uri)
	if err != nil {
		log.Panic(err.Error())
	}

	return func(u *upp) {
		u.cache = cache
	}
}

func InitES(uri string) module {
	client, err := es.New(context.TODO(), uri)
	if err != nil {
		log.Panic(err.Error())
	}

	return func(u *upp) {
		u.es = client
	}
}

type ApiConfig struct {
	Address   string
	TLSConfig *tls.Config
}

func InitApi(api *api.App, cfgs ...ApiConfig) module {
	cfg := ApiConfig{}
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}

	if cfg.Address == "" {
		cfg.Address = "localhost:8080"
	}

	return func(u *upp) {
		api.Upp = u
		u.api = &uppApi{
			engine: api,
			config: cfg,
		}
	}
}

func InitTaskChan(ch <-chan func(upp interfaces.Upp) error) module {
	return func(u *upp) {
		if u.taskCh == nil {
			u.taskCh = make([]<-chan func(u interfaces.Upp) error, 0)
		}

		u.taskCh = append(u.taskCh, ch)
	}
}

// sync functions
// 添加 同步执行函数
func InitFn(fns ...func(interfaces.Upp)) module {
	return func(u *upp) {
		u.initFns._sync = append(u.initFns._sync, fns...)
	}
}

// async functions
// 添加 异步执行函数
func InitAsyncFn(fns ...func(interfaces.Upp)) module {
	return func(u *upp) {
		u.initFns._async = append(u.initFns._async, fns...)
	}
}
