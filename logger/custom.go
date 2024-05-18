package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

// make sure it's idempotent
var once sync.Once

type (
	// ReplaceAttrFunc to provide a way to change key and value of slog.Attr
	ReplaceAttrFunc func(groups []string, a slog.Attr) slog.Attr

	// Handle to provide a way to add values from context to slog.Record
	HandleFunc func(ctx context.Context, rec slog.Record) (slog.Record, error)
)

type CustomLogger struct {
	writer      io.Writer
	replaceAttr ReplaceAttrFunc
	handle      HandleFunc
	addSource   bool
	level       string
	isJSON      bool
}

var logger *slog.Logger

func WithWriter(wr io.Writer) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.writer = wr
	}
}

func WithReplaceAttr(f ReplaceAttrFunc) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.replaceAttr = f
	}
}

func WithSource(b bool) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.addSource = b
	}
}

func WithHandle(h HandleFunc) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.handle = h
	}
}

func WithLevel(l string) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.level = l
	}
}

func WithJSON(b bool) func(CustomLogger) {
	return func(cl CustomLogger) {
		cl.isJSON = b
	}
}

func New(opts ...func(l CustomLogger)) *slog.Logger {
	once.Do(func() {
		handler := new(opts...)
		logger = slog.New(handler)
	})

	return logger
}

func new(opts ...func(l CustomLogger)) customHandler {
	l := CustomLogger{
		writer:      os.Stdout,
		replaceAttr: replaceAttr,
		isJSON:      true,
		addSource:   true,
		level:       "warn",
	}

	for _, opt := range opts {
		opt(l)
	}

	options := slog.HandlerOptions{
		AddSource:   l.addSource,
		Level:       getLevel(l.level),
		ReplaceAttr: l.replaceAttr,
	}

	if l.isJSON {
		return customHandler{
			Handler:    slog.NewJSONHandler(l.writer, &options),
			handleFunc: l.handle,
		}
	}

	return customHandler{
		Handler:    slog.NewTextHandler(l.writer, &options),
		handleFunc: l.handle,
	}
}

// customHandler for changing behaviour as per need
type customHandler struct {
	slog.Handler
	handleFunc HandleFunc
}

func (c customHandler) Handle(ctx context.Context, rec slog.Record) error {
	var err error

	if c.handleFunc != nil {
		rec, err = c.handleFunc(ctx, rec)
		if err != nil {
			return err
		}
	}

	return c.Handler.Handle(ctx, rec)
}
