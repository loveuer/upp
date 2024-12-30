package upp

import (
	"context"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/loveuer/upp/pkg/cache"
	"github.com/loveuer/upp/pkg/interfaces"
	"github.com/loveuer/upp/pkg/tool"
	"gorm.io/gorm"
)

func (u *upp) UseCtx() context.Context {
	return u.ctx
}

func (u *upp) UseDB(ctx ...context.Context) *gorm.DB {
	var c context.Context

	if len(ctx) > 0 {
		c = ctx[0]
	} else {
		c = tool.Timeout(30)
	}

	tx := u.db.Session(&gorm.Session{
		Context: c,
	})

	if u.debug {
		tx = tx.Debug()
	}

	return tx
}

func (u *upp) UseCache() cache.Cache {
	return u.cache
}

func (u *upp) UseES() *elasticsearch.Client {
	return u.es
}

func (u *upp) UseLogger(ctxs ...context.Context) interfaces.Logger {
	logger := u.logger.Get().(*upp_logger)

	logger.ctx = u.UseCtx()
	if len(ctxs) > 0 {
		logger.ctx = ctxs[0]
	}

	return logger
}
