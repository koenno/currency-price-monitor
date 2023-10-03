package request

import (
	"fmt"
	"testing"
	"time"

	"github.com/koenno/currency-price-monitor/request/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShouldContinueWritingElementsDespiteOfError(t *testing.T) {
	// given
	input := make(chan Descriptor[string], 2)
	input <- newDescriptor("1")
	input <- newDescriptor("2")
	close(input)
	writerMock := mocks.NewWriter(t)
	callCnt := 0
	writerMock.EXPECT().Write(mock.Anything).RunAndReturn(func(b []byte) (int, error) {
		defer func() {
			callCnt++
		}()
		if callCnt == 0 {
			return 0, fmt.Errorf("failed to write")
		}
		return len(b), nil
	})

	// when
	Out(input, writerMock)

	// then
	writerMock.AssertExpectations(t)
	assert.Equal(t, 2, callCnt)
}

func newDescriptor(ID string) Descriptor[string] {
	return Descriptor[string]{
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
