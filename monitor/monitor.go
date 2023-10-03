package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/koenno/currency-price-monitor/request"
	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Requester --case underscore --with-expecter
type Requester interface {
	Process(*http.Request) (request.Descriptor, error)
}

type Monitor struct {
	requester Requester
	request   *http.Request
}

func New(requester Requester, request *http.Request) Monitor {
	return Monitor{
		requester: requester,
		request:   request,
	}
}

func (m Monitor) Start(ctx context.Context, requestsNumber uint, interval time.Duration) <-chan request.Descriptor {
	output := make(chan request.Descriptor)
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

func (m Monitor) update(number uint, output chan<- request.Descriptor) {
	for i := 0; i < int(number); i++ {
		m.singleUpdate(output)
	}
}

func (m Monitor) singleUpdate(output chan<- request.Descriptor) {
	desc, err := m.requester.Process(m.request)
	if err != nil {
		slog.Error("monitor failed to process a request", "error", err)
	}
	output <- desc
}
