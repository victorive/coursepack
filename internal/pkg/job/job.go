package job

import (
	"context"
	"log/slog"
	"runtime/debug"
)

func FireAndForget(ctx context.Context, name string, fn func(context.Context) error) {
	bgCtx := context.WithoutCancel(ctx)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Default().Error("job panic", "job", name, "panic", r, "stack", string(debug.Stack()))
			}
		}()
		if err := fn(bgCtx); err != nil {
			slog.Default().Warn("job error", "job", name, "error", err)
		}
	}()
}
