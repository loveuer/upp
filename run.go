package upp

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/loveuer/upp/pkg/interfaces"
	"github.com/loveuer/upp/pkg/tool"
)

func (u *upp) StartAPI(ctx context.Context) {
	address := env.ListenHttp
	if address == "" {
		address = u.api.config.Address
	}

	u.UseLogger().Info("UPP | run api at %s", address)
	go u.api.engine.Run(address)
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

	ctx, cancel := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
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
	if u.cache != nil {
		u.cache.Close()
	}

	<-tool.Timeout(2).Done()
}
