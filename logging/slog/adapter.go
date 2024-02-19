//go:build go1.21
// +build go1.21

/*
Package goaslog contains an adapter that makes it possible to configure goa so it uses [log/slog]
as logger backend.
Usage:

	handler := slog.NewJSONHandler(os.Stderr, nil)
	// Initialize logger handler using [log/slog] package
	service.WithLogger(goaslog.New(handler))
	// ... Proceed with configuring and starting the goa service

	// In handlers:
	goaslog.Entry(ctx).Info("foo", "bar")
*/
package goaslog

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/shogo82148/goa-v1"
)

var _ goa.LogAdapter = (*adapter)(nil)
var _ goa.ContextLogAdapter = (*adapter)(nil)

// adapter is the slog goa logger adapter.
type adapter struct {
	handler slog.Handler
}

// New wraps a [log/slog.Handler] into a goa logger.
func New(handler slog.Handler) goa.LogAdapter {
	return &adapter{handler: handler}
}

// Info logs messages using [log/slog].
func (a *adapter) Info(msg string, data ...any) {
	a.log(context.Background(), slog.LevelInfo, msg, data...)
}

// InfoContext logs messages using [log/slog].
func (a *adapter) InfoContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelInfo, msg, data...)
}

// Warn logs message using [log/slog].
func (a *adapter) Warn(msg string, data ...any) {
	a.log(context.Background(), slog.LevelWarn, msg, data...)
}

// WarnContext logs message using [log/slog].
func (a *adapter) WarnContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelWarn, msg, data...)
}

// Error logs errors using [log/slog].
func (a *adapter) Error(msg string, data ...any) {
	a.log(context.Background(), slog.LevelError, msg, data...)
}

// ErrorContext logs errors using [log/slog].
func (a *adapter) ErrorContext(ctx context.Context, msg string, data ...any) {
	a.log(ctx, slog.LevelError, msg, data...)
}

// New creates a new logger given a context.
func (a *adapter) New(data ...any) goa.LogAdapter {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "", 0)
	r.Add(data...)

	attrs := make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	h := a.handler.WithAttrs(attrs)
	return &adapter{handler: h}
}

func (a *adapter) log(ctx context.Context, level slog.Level, msg string, data ...any) {
	if !a.handler.Enabled(ctx, level) {
		return
	}

	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this functions, this functions caller, the caller of the adapter]
	runtime.Callers(4, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(data...)
	a.handler.Handle(ctx, r)
}
