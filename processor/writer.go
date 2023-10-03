package processor

import (
	"context"
	"io"

	"github.com/koenno/currency-price-monitor/request"
)

type Writer struct {
	out io.Writer
}

func NewWriter[T any](out io.Writer) Writer {
	return Writer{
		out: out,
	}
}

func (w Writer) Process(ctx context.Context, desc request.Descriptor) error {
	_, err := desc.WriteTo(w.out)
	return err
}
