package upp

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/loveuer/upp/pkg/interfaces"
	"github.com/loveuer/upp/pkg/tool"
)

func (u *upp) startAPI(ctx context.Context) {
	address := env.ListenHttp
	if address == "" {
		address = u.api.config.Address
	}

	fmt.Printf("Upp | api listen at %s\n", address)
	go u.api.engine.Run(address)
	go func() {
		<-ctx.Done()
		u.api.engine.Shutdown(tool.Timeout(2))
	}()
}

func (u *upp) startTask(ctx context.Context) {
	fmt.Printf("Upp | start task channel[%02d]", len(u.taskCh))
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

func (u *upp) runInitFns(ctx context.Context) {
	for _, fn := range u.initFns._sync {
		fn(u)
	}
}

func (u *upp) startInitFns(ctx context.Context) {
	for _, fn := range u.initFns._async {
		go fn(u)
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

	print(Banner)

	if len(u.initFns._sync) > 0 {
		u.runInitFns(ctx)
	}

	if len(u.initFns._async) > 0 {
		u.startInitFns(ctx)
	}

	if u.api != nil {
		u.startAPI(ctx)
	}

	if len(u.taskCh) > 0 {
		u.startTask(ctx)
	}

	<-ctx.Done()

	u.UseLogger().Warn(" Upp | quit by signal...")
	if u.cache != nil {
		u.cache.Close()
	}

	<-tool.Timeout(2).Done()
}
