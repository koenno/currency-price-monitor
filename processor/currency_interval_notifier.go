package processor

import (
	"context"
	"io"

	"github.com/koenno/currency-price-monitor/request"
)

const (
	multiplier = 10000
)

type ClosedInterval struct {
	A float64
	B float64
}

type CurrencyIntervalWriter struct {
	out      io.Writer
	interval ClosedInterval
}

func NewCurrencyIntervalNotifier(out io.Writer, interval ClosedInterval) CurrencyIntervalWriter {
	return CurrencyIntervalWriter{
		out:      out,
		interval: interval,
	}
}

func (n CurrencyIntervalWriter) Process(ctx context.Context, desc request.Descriptor) error {
	for i := 0; i < len(desc.Payload.Rates); i++ {
		if desc.Payload.Rates[i].Value < n.interval.A || desc.Payload.Rates[i].Value > n.interval.B {
			desc.Payload.Rates[i].WriteTo(n.out)
		}
	}
	return nil
}
