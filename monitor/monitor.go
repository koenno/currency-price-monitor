package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/koenno/currency-price-monitor/request"
	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Requester --case underscore --with-expecter
type Requester[T any] interface {
	Process(req *http.Request) (request.Descriptor[T], error)
}

type Monitor[T any] struct {
	requester Requester[T]
	request   *http.Request
}

func New[T any](requester Requester[T], request *http.Request) Monitor[T] {
	return Monitor[T]{
		requester: requester,
		request:   request,
	}
}

func (m Monitor[T]) Start(ctx context.Context, requestsNumber uint, interval time.Duration) <-chan request.Descriptor[T] {
	output := make(chan request.Descriptor[T])
	go func() {
		defer close(output)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		m.update(requestsNumber, output)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.update(requestsNumber, output)
			}
		}
	}()
	return output
}

func (m Monitor[T]) update(number uint, output chan<- request.Descriptor[T]) {
	for i := 0; i < int(number); i++ {
		m.singleUpdate(output)
	}
}

func (m Monitor[T]) singleUpdate(output chan<- request.Descriptor[T]) {
	desc, err := m.requester.Process(m.request)
	if err != nil {
		slog.Error("monitor failed to process a request", "error", err)
	}
	output <- desc
}
