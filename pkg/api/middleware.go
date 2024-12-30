package api

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/resp"
	"github.com/loveuer/upp/pkg/tool"
)

func NewRecover(enableStackTrace bool) HandlerFunc {
	return func(c *Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				if enableStackTrace {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\n%s\n", r, debug.Stack()))
				} else {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\n", r))
				}

				_ = c.Status(500).SendString(fmt.Sprint(r))
			}
		}()

		return c.Next()
	}
}

func NewLogger() HandlerFunc {
	return func(c *Ctx) error {
		var (
			now   = time.Now()
			logFn func(msg string, data ...any)
			ip    = c.IP()
		)

		traceId := c.Context().Value(nf.TraceKey)
		c.Locals(nf.TraceKey, traceId)

		err := c.Next()

		c.Writer.Header().Set(nf.TraceKey, fmt.Sprint(traceId))

		status, _ := strconv.Atoi(c.Writer.Header().Get(resp.RealStatusHeader))
		duration := time.Since(now)

		msg := fmt.Sprintf("%s | %15s | %d[%3d] | %s | %6s | %s", traceId, ip, c.StatusCode, status, tool.HumanDuration(duration.Nanoseconds()), c.Method(), c.Path())

		switch {
		case status >= 500:
			logFn = log.Error
		case status >= 400:
			logFn = log.Warn
		default:
			logFn = log.Info
		}

		logFn(msg)

		return err
	}
}
