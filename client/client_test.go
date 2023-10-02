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
	resp, err := Request[string](context.Background(), fakeServer.URL)

	// then
	assert.ErrorIs(t, err, ErrResponse)
	assert.Equal(t, ResponseDescriptor[string]{}, resp)
}

func TestShouldReturnErrorWhenContentTypeIsNotJSON(t *testing.T) {
	// given
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/xml")
	}))

	// when
	resp, err := Request[string](context.Background(), fakeServer.URL)

	// then
	assert.ErrorIs(t, err, ErrResponsePayload)
	assert.Equal(t, ResponseDescriptor[string]{}, resp)
}

func TestShouldReturnErrorWhenJSONIsInvalid(t *testing.T) {
	// given
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Write([]byte(`{ "invalidJson": true `))
	}))

	// when
	resp, err := Request[string](context.Background(), fakeServer.URL)

	// then
	assert.ErrorIs(t, err, ErrResponsePayload)
	assert.Equal(t, ResponseDescriptor[string]{}, resp)
}

func TestShouldReturnResponseDescriptorWhenNoError(t *testing.T) {
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
	resp, err := Request[payload](context.Background(), fakeServer.URL)

	// then
	assert.NoError(t, err)
	assert.NotZero(t, resp.ID)
	assert.NotZero(t, resp.Time)
	assert.Equal(t, fakeServer.URL, resp.URL)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, resp.JSON)
	assert.True(t, resp.Valid)
	assert.NotZero(t, resp.Duration)
	assert.Equal(t, expectedPayload, resp.Payload)
}
