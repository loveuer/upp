package upp

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/loveuer/upp/pkg/api"
	"github.com/loveuer/upp/pkg/log"
)

type upp_logger struct {
	ctx context.Context
}

var upp_logger_pool = &sync.Pool{
	New: func() any {
		return &upp_logger{}
	},
}

func (ul *upp_logger) Debug(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Debug(traceId+" | "+msg, data...)

	upp_logger_pool.Put(ul)
}

func (ul *upp_logger) Info(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Info(traceId+" | "+msg, data...)

	upp_logger_pool.Put(ul)
}

func (ul *upp_logger) Warn(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Warn(traceId+" | "+msg, data...)

	upp_logger_pool.Put(ul)
}

func (ul *upp_logger) Error(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Error(traceId+" | "+msg, data...)

	upp_logger_pool.Put(ul)
}

func (ul *upp_logger) Panic(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Panic(traceId+" | "+msg, data...)
}

func (ul *upp_logger) Fatal(msg string, data ...any) {
	traceId, ok := ul.ctx.Value(api.TraceKey).(string)
	if !ok {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	log.Fatal(traceId+" | "+msg, data...)

	upp_logger_pool.Put(ul)
}
