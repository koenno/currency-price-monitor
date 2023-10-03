package scheduler

import (
	"context"
	"errors"
	"sync"

	"github.com/koenno/currency-price-monitor/request"
	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Processor --case underscore --with-expecter
type Processor interface {
	Process(context.Context, request.Descriptor) error
}

type Scheduler struct {
	processors []Processor
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (r *Scheduler) Register(processor Processor) {
	r.processors = append(r.processors, processor)
}

func (s *Scheduler) Process(ctx context.Context, input <-chan request.Descriptor) {
	for desc := range input {
		err := s.processSingle(ctx, desc)
		if err != nil {
			slog.Error("failure while processing descriptor", "id", desc.ID, "error", err)
		}
	}
}

func (s *Scheduler) processSingle(ctx context.Context, desc request.Descriptor) error {
	var (
		wg      sync.WaitGroup
		errsMtx sync.Mutex
		errs    []error
	)
	wg.Add(len(s.processors))
	for _, p := range s.processors {
		go func(p Processor) {
			defer wg.Done()
			defer errsMtx.Unlock()
			err := p.Process(ctx, desc)
			errsMtx.Lock()
			errs = append(errs, err)
		}(p)
	}
	wg.Wait()
	return errors.Join(errs...)
}
