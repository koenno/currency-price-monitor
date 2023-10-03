package processor

import (
	"context"
	"io"

	"github.com/koenno/currency-price-monitor/request"
)

type Writer[T any] struct {
	out io.Writer
}

func NewWriter[T any](out io.Writer) Writer[T] {
	return Writer[T]{
		out: out,
	}
}

func (w Writer[T]) Process(ctx context.Context, desc request.Descriptor[T]) error {
	_, err := desc.WriteTo(w.out)
	return err
}
