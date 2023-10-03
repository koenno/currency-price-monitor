package request

import (
	"io"

	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Writer --case underscore --with-expecter
type Writer interface {
	io.Writer
}

func Out[T any](input <-chan Descriptor[T], output Writer) {
	for desc := range input {
		_, err := desc.WriteTo(output)
		if err != nil {
			slog.Error("output failure", "id", desc.ID, "error", err)
		}
	}
}
