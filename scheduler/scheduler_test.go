package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/koenno/currency-price-monitor/request"
	"github.com/koenno/currency-price-monitor/scheduler/mocks"
	"github.com/stretchr/testify/mock"
)

func TestShouldCallAllRegisteredProcessors(t *testing.T) {
	// given
	procMock1 := mocks.NewProcessor(t)
	procMock2 := mocks.NewProcessor(t)
	input := make(chan request.Descriptor, 2)
	descs := []request.Descriptor{newDescriptor("1"), newDescriptor("2")}
	input <- descs[0]
	input <- descs[1]
	close(input)

	sut := NewScheduler()
	sut.Register(procMock1)
	sut.Register(procMock2)

	procMock1.EXPECT().Process(mock.Anything, descs[0]).Return(nil).Once()
	procMock1.EXPECT().Process(mock.Anything, descs[1]).Return(nil).Once()
	procMock2.EXPECT().Process(mock.Anything, descs[0]).Return(nil).Once()
	procMock2.EXPECT().Process(mock.Anything, descs[1]).Return(nil).Once()

	// when
	sut.Process(context.Background(), input)

	// then
	procMock1.AssertExpectations(t)
	procMock2.AssertExpectations(t)
}

func TestShouldContinueProcessingWhenSomeProcessorsFail(t *testing.T) {
	// given
	procMock1 := mocks.NewProcessor(t)
	procMock2 := mocks.NewProcessor(t)
	input := make(chan request.Descriptor, 2)
	descs := []request.Descriptor{newDescriptor("1"), newDescriptor("2")}
	input <- descs[0]
	input <- descs[1]
	close(input)

	sut := NewScheduler()
	sut.Register(procMock1)
	sut.Register(procMock2)

	procMock1.EXPECT().Process(mock.Anything, descs[0]).Return(errors.New("failure 1")).Once()
	procMock1.EXPECT().Process(mock.Anything, descs[1]).Return(nil).Once()
	procMock2.EXPECT().Process(mock.Anything, descs[0]).Return(errors.New("failure 2")).Once()
	procMock2.EXPECT().Process(mock.Anything, descs[1]).Return(nil).Once()

	// when
	sut.Process(context.Background(), input)

	// then
	procMock1.AssertExpectations(t)
	procMock2.AssertExpectations(t)
}

func newDescriptor(ID string) request.Descriptor {
	return request.Descriptor{
		ID:              ID,
		URL:             "http://some/domain.com",
		Time:            time.Time{},
		ValidStatusCode: true,
		JSON:            true,
		Valid:           true,
		Duration:        0,
		Payload:         request.Currency{},
	}
}
