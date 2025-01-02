package interfaces

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/loveuer/upp/pkg/cache"
	"gorm.io/gorm"
)

type Upp interface {
	Debug() bool
	UseCtx() context.Context
	UseDB(ctx ...context.Context) *gorm.DB
	UseCache() cache.Cache
	UseES() *elasticsearch.Client
	UseLogger(ctxs ...context.Context) Logger
}
