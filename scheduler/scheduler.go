package scheduler

import (
	"context"
	"errors"
	"sync"

	"github.com/koenno/currency-price-monitor/request"
	"golang.org/x/exp/slog"
)

//go:generate mockery --name=Processor --case underscore --with-expecter
type Processor[T any] interface {
	Process(context.Context, request.Descriptor[T]) error
}

type Scheduler[T any] struct {
	processors []Processor[T]
}

func NewScheduler[T any]() *Scheduler[T] {
	return &Scheduler[T]{}
}

func (r *Scheduler[T]) Register(processor Processor[T]) {
	r.processors = append(r.processors, processor)
}

func (s *Scheduler[T]) Process(ctx context.Context, input <-chan request.Descriptor[T]) {
	for desc := range input {
		err := s.processSingle(ctx, desc)
		if err != nil {
			slog.Error("failure while processing descriptor", "id", desc.ID, "error", err)
		}
	}
}

func (s *Scheduler[T]) processSingle(ctx context.Context, desc request.Descriptor[T]) error {
	var (
		wg      sync.WaitGroup
		errsMtx sync.Mutex
		errs    []error
	)
	wg.Add(len(s.processors))
	for _, p := range s.processors {
		go func(p Processor[T]) {
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
