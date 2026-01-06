package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type HandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type Handler struct {
	slog.Handler
	l *log.Logger
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	if r.NumAttrs() > 0 {
		fields := make(map[string]any, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			fields[a.Key] = a.Value.Any()

			return true
		})

		b, err := json.MarshalIndent(fields, "", "\t")
		if err != nil {
			return err
		}

		h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))
	} else {
		h.l.Println(timeStr, level, msg)
	}

	return nil
}

func NewHandler(
	out io.Writer,
	opts HandlerOptions,
) *Handler {
	h := &Handler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
