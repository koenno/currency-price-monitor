package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/koenno/currency-price-monitor/client/mocks"
	"github.com/koenno/currency-price-monitor/request"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenResponseStatusCodeIsNotSuccessful(t *testing.T) {
	// given
	converterMock := mocks.NewConverter[string](t)
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New[string](converterMock)

	// when
	desc, err := sut.Process(req)

	// then
	assert.ErrorIs(t, err, ErrResponse)
	assert.NotZero(t, desc.ID)
	assert.NotZero(t, desc.Time)
	assert.Equal(t, fakeServer.URL, desc.URL)
	assert.False(t, desc.ValidStatusCode)
	assert.False(t, desc.JSON)
	assert.False(t, desc.Valid)
	assert.NotZero(t, desc.Duration)
	assert.Zero(t, desc.Payload)
}

func TestShouldReturnErrorWhenContentTypeIsNotJSON(t *testing.T) {
	// given
	converterMock := mocks.NewConverter[string](t)
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/xml")
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New[string](converterMock)

	// when
	desc, err := sut.Process(req)

	// then
	assert.ErrorIs(t, err, ErrResponsePayload)
	assert.NotZero(t, desc.ID)
	assert.NotZero(t, desc.Time)
	assert.Equal(t, fakeServer.URL, desc.URL)
	assert.True(t, desc.ValidStatusCode)
	assert.False(t, desc.JSON)
	assert.False(t, desc.Valid)
	assert.NotZero(t, desc.Duration)
	assert.Zero(t, desc.Payload)
}

func TestShouldReturnErrorWhenJSONIsInvalid(t *testing.T) {
	// given
	converterMock := mocks.NewConverter[string](t)
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Write([]byte(`{ "invalidJson": true `))
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New[string](converterMock)

	// when
	desc, err := sut.Process(req)

	// then
	assert.ErrorIs(t, err, ErrResponsePayload)
	assert.NotZero(t, desc.ID)
	assert.NotZero(t, desc.Time)
	assert.Equal(t, fakeServer.URL, desc.URL)
	assert.True(t, desc.ValidStatusCode)
	assert.True(t, desc.JSON)
	assert.False(t, desc.Valid)
	assert.NotZero(t, desc.Duration)
	assert.Zero(t, desc.Payload)
}

func TestShouldFillAllDescriptorDataWhenNoError(t *testing.T) {
	// given
	type payload struct {
		Name   string
		Number int
	}
	expectedPayload := payload{
		Name:   "test",
		Number: 234,
	}
	converterMock := mocks.NewConverter[payload](t)
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(expectedPayload)
	}))
	req, _ := http.NewRequest(http.MethodGet, fakeServer.URL, nil)
	sut := New[payload](converterMock)

	expectedCurrency := request.Currency{
		Name: "eur",
		Rates: []request.Rate{
			{
				Date:  newDate("2023-10-03"),
				Value: 4.75,
			},
		},
	}
	converterMock.EXPECT().Convert(expectedPayload).Return(expectedCurrency).Once()

	// when
	desc, err := sut.Process(req)

	// then
	assert.NoError(t, err)
	assert.NotZero(t, desc.ID)
	assert.NotZero(t, desc.Time)
	assert.Equal(t, fakeServer.URL, desc.URL)
	assert.True(t, desc.ValidStatusCode)
	assert.True(t, desc.JSON)
	assert.True(t, desc.Valid)
	assert.NotZero(t, desc.Duration)
	assert.Equal(t, expectedCurrency, desc.Payload)
}

func newDate(date string) time.Time {
	t, _ := time.Parse(time.DateOnly, date)
	return t
}
