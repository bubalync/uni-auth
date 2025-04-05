package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ErrMap(errs map[string]string) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.AnyValue(errs),
	}
}
