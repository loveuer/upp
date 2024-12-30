package upp

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/loveuer/upp/pkg/api"
	"github.com/loveuer/upp/pkg/cache"
	"github.com/loveuer/upp/pkg/interfaces"
	"github.com/loveuer/upp/pkg/tool"
	"gorm.io/gorm"
)

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
	initFns []func(interfaces.Upp)
	taskCh  []<-chan func(interfaces.Upp) error
}

func (u *upp) With(modules ...module) {
	for _, m := range modules {
		m(u)
	}
}

func (u *upp) Run(ctx context.Context) {
	u.RunSignal(ctx)
}

func (u *upp) RunInitFns(ctx context.Context) {
	for _, fn := range u.initFns {
		fn(u)
	}
}

func (u *upp) RunSignal(ctxs ...context.Context) {
	c := context.Background()
	if len(ctxs) > 0 {
		c = ctxs[0]
	}

	ctx, cancel := signal.NotifyContext(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	u.ctx = ctx

	if len(u.initFns) > 0 {
		u.RunInitFns(ctx)
	}

	if u.api != nil {
		u.StartAPI(ctx)
	}

	if len(u.taskCh) > 0 {
		u.StartTask(ctx)
	}

	<-ctx.Done()

	u.UseLogger().Warn(" UPP | quit by signal...")

	<-tool.Timeout(2).Done()
}

func (u *upp) StartAPI(ctx context.Context) {
	u.UseLogger().Info("UPP | run api at %s", u.api.config.Address)
	go u.api.engine.Run(u.api.config.Address)
	go func() {
		<-ctx.Done()
		u.api.engine.Shutdown(tool.Timeout(2))
	}()
}

func (u *upp) StartTask(ctx context.Context) {
	for _, _ch := range u.taskCh {
		go func(ch <-chan func(interfaces.Upp) error) {
			var err error
			for {
				select {
				case <-ctx.Done():
				case task, ok := <-ch:
					if !ok {
						return
					}

					if err = task(u); err != nil {
						u.UseLogger(ctx).Error(err.Error())
					}
				}
			}
		}(_ch)
	}
}

func New(configs ...Config) *upp {
	app := &upp{
		logger: upp_logger_pool,
	}

	return app
}
