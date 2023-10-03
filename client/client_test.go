package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenResponseStatusCodeIsNotSuccessful(t *testing.T) {
	// given
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	// when
	desc, err := Request[string](context.Background(), fakeServer.URL)

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
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/xml")
	}))

	// when
	desc, err := Request[string](context.Background(), fakeServer.URL)

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
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Write([]byte(`{ "invalidJson": true `))
	}))

	// when
	desc, err := Request[string](context.Background(), fakeServer.URL)

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
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(expectedPayload)
	}))

	// when
	desc, err := Request[payload](context.Background(), fakeServer.URL)

	// then
	assert.NoError(t, err)
	assert.NotZero(t, desc.ID)
	assert.NotZero(t, desc.Time)
	assert.Equal(t, fakeServer.URL, desc.URL)
	assert.True(t, desc.ValidStatusCode)
	assert.True(t, desc.JSON)
	assert.True(t, desc.Valid)
	assert.NotZero(t, desc.Duration)
	assert.Equal(t, expectedPayload, desc.Payload)
}
