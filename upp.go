package upp

import (
	"context"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/loveuer/upp/pkg/api"
	"github.com/loveuer/upp/pkg/cache"
	"github.com/loveuer/upp/pkg/interfaces"
	"github.com/loveuer/upp/pkg/log"
	"gorm.io/gorm"
)

const Banner = `
  __  __        
 / / / /__  ___ 
/ /_/ / _ \/ _ \
\____/ .__/ .__/
    /_/  /_/    

`

type uppApi struct {
	engine *api.App
	config ApiConfig
}

type upp struct {
	debug   bool
	ctx     context.Context
	logger  *sync.Pool
	db      *gorm.DB
	cache   cache.Cache
	es      *elasticsearch.Client
	api     *uppApi
	initFns struct {
		_sync  []func(interfaces.Upp)
		_async []func(interfaces.Upp)
	}
	taskCh []<-chan func(interfaces.Upp) error
}

func (u *upp) With(modules ...module) {
	for _, m := range modules {
		m(u)
	}
}

func New(configs ...Config) *upp {
	config := Config{}

	if len(configs) > 0 {
		config = configs[0]
	}

	app := &upp{
		logger: upp_logger_pool,
		initFns: struct {
			_sync  []func(interfaces.Upp)
			_async []func(interfaces.Upp)
		}{
			_sync:  make([]func(interfaces.Upp), 0),
			_async: make([]func(interfaces.Upp), 0),
		},
	}

	if config.Debug || env.Debug {
		log.SetLogLevel(log.LogLevelDebug)
		app.debug = true
	}

	return app
}
