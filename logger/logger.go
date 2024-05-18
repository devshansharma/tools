package logger

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/mdobak/go-xerrors"
)

// GetLogger will return logger instance
func GetLogger() *slog.Logger {
	if logger == nil {
		panic("initialize logger first")
	}

	return logger
}

func getLevel(s string) slog.Level {
	s = strings.ToUpper(s)

	switch s {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

// WithShortFileNameAndErrorTrace for Short File Name and Error Trace
func WithShortFileNameAndErrorTrace(groups []string, a slog.Attr) slog.Attr {
	a = WithShortFileName(groups, a)
	return WithErrorTrace(groups, a)
}

// WithShortFileName to replace attr source.File with short name
func WithShortFileName(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, ok := a.Value.Any().(*slog.Source)
		if ok {
			source.File = filepath.Base(source.File)
		}
	}

	return a
}

// WithErrorTrace for showing error trace in logs
func WithErrorTrace(groups []string, a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindAny:
		switch v := a.Value.Any().(type) {
		case error:
			a.Value = fmtErr(v)
		}
	}

	return a
}

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

// fmtErr returns a slog.Value with keys `msg` and `trace`. If the error
// does not implement interface { StackTrace() errors.StackTrace }, the `trace`
// key is omitted.
func fmtErr(err error) slog.Value {
	var groupValues []slog.Attr

	groupValues = append(groupValues, slog.String("msg", err.Error()))

	frames := marshalStack(err)

	if frames != nil {
		groupValues = append(groupValues,
			slog.Any("trace", frames),
		)
	}

	return slog.GroupValue(groupValues...)
}

// marshalStack extracts stack frames from the error
func marshalStack(err error) []stackFrame {
	trace := xerrors.StackTrace(err)

	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()

	s := make([]stackFrame, len(frames))

	for i, v := range frames {
		f := stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(v.File)),
				filepath.Base(v.File),
			),
			Func: filepath.Base(v.Function),
			Line: v.Line,
		}

		s[i] = f
	}

	return s
}
