package server

import "log/slog"

type Closer interface {
	Close() error
}

func CloseDataBase(log *slog.Logger, closer Closer) error {
	err := closer.Close()
	if err != nil {
		log.Error("failed to close database", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	return nil
}
