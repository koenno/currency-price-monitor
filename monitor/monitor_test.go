package monitor

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/koenno/currency-price-monitor/monitor/mocks"
	"github.com/koenno/currency-price-monitor/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldStopWhenContextIsCancelled(t *testing.T) {
	// given
	requestsNumber := 1
	requestsInterval := time.Minute
	requesterMock := mocks.NewRequester[string](t)
	req, _ := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	sut := New[string](requesterMock, req)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	desc := newDescriptor("1")
	requesterMock.EXPECT().Process(mock.Anything).Return(desc, nil)

	// when
	output := sut.Start(ctx, uint(requestsNumber), requestsInterval)

	// then
	var descs []request.Descriptor[string]
	for d := range output {
		descs = append(descs, d)
	}
	assert.Equal(t, requestsNumber, len(descs))
}

func TestShouldSendReceivedDescriptorToOutputChannel(t *testing.T) {
	// given
	requestsNumber := 1
	requestsInterval := time.Minute
	requesterMock := mocks.NewRequester[string](t)
	req, _ := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	sut := New[string](requesterMock, req)
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(500 * time.Millisecond)
	go func() {
		<-timer.C
		cancel()
	}()

	desc := newDescriptor("1")
	requesterMock.EXPECT().Process(mock.Anything).Return(desc, nil).Times(requestsNumber)

	// when
	output := sut.Start(ctx, uint(requestsNumber), requestsInterval)

	// then
	var descs []request.Descriptor[string]
	for d := range output {
		descs = append(descs, d)
	}
	assert.Equal(t, requestsNumber, len(descs))
	assert.Equal(t, desc, descs[0])
}

func TestShouldUpdateMonitorWithTimeInterval(t *testing.T) {
	// given
	expectedAllRequestsNumber := 4
	requestsNumber := 2
	requestsInterval := 100 * time.Millisecond
	requesterMock := mocks.NewRequester[string](t)
	req, _ := http.NewRequest(http.MethodGet, "some.domain.com", nil)
	sut := New[string](requesterMock, req)
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(200 * time.Millisecond)
	go func() {
		<-timer.C
		cancel()
	}()

	desc := newDescriptor("1")
	requesterMock.EXPECT().Process(mock.Anything).Return(desc, nil)

	// when
	output := sut.Start(ctx, uint(requestsNumber), requestsInterval)

	// then
	var descs []request.Descriptor[string]
	for d := range output {
		descs = append(descs, d)
	}
	assert.GreaterOrEqual(t, len(descs), expectedAllRequestsNumber)
}

func newDescriptor(ID string) request.Descriptor[string] {
	return request.Descriptor[string]{
		ID:              ID,
		URL:             "http://some/domain.com",
		Time:            time.Time{},
		ValidStatusCode: true,
		JSON:            true,
		Valid:           true,
		Duration:        0,
		Payload:         "",
	}
}
