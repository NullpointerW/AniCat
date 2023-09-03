package log

import (
	"context"
	"log/slog"
	"runtime"
)

type EnhanceHandler struct {
	slog.Handler
	calldepth int
}

func (eh *EnhanceHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		pc  uintptr
		pcs [1]uintptr
		// skip enhanceLogger wrapped function =>depth+1
		// skip [runtime.Callers, this function, this function's caller] in slog.log =>depth+3
		// skip runtime.Callers in *EnhanceHandler.Handle =>depth+1
		// so depth=5
		depth = 5
	)

	runtime.Callers(depth+eh.calldepth, pcs[:])
	pc = pcs[0]
	r.PC = pc
	return eh.Handler.Handle(ctx, r)
}

func NewEnhanceHandler(h slog.Handler, calldepth int) slog.Handler {
	return &EnhanceHandler{h, calldepth}
}
