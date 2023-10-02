package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrSendRequest     = errors.New("failed to send request")
	ErrResponse        = errors.New("response failure")
	ErrResponsePayload = errors.New("erroneus response payload")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type ResponseDescriptor[T any] struct {
	URL        string
	StatusCode int
	JSON       bool
	Valid      bool
	Duration   time.Duration
	Payload    T
}

func Request[T any](ctx context.Context, URL string) (ResponseDescriptor[T], error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return ResponseDescriptor[T]{}, fmt.Errorf("request failure: unable to create a request: %v", err)
	}

	startTime := time.Now()
	resp, err := httpClient.Do(req)
	if err != nil {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: %v", ErrSendRequest, err)
	}
	reqDuration := time.Since(startTime)

	defer resp.Body.Close()
	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: unable to read body: %v", ErrResponse, err)
	}

	if resp.StatusCode != http.StatusOK {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: status code %d", ErrResponse, resp.StatusCode)
	}

	if resp.Header.Get("content-type") != "application/json" {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: unsupported %s", ErrResponsePayload, resp.Header.Get("content-type"))
	}

	validJSON := json.Valid(payloadBytes)
	if !validJSON {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: invalid json", ErrResponsePayload)
	}

	var data T
	err = json.Unmarshal(payloadBytes, &data)
	if err != nil {
		return ResponseDescriptor[T]{}, fmt.Errorf("%w: unable to decode body to json %v", ErrResponse, err)
	}

	return ResponseDescriptor[T]{
		URL:        URL,
		StatusCode: resp.StatusCode,
		JSON:       true,
		Valid:      validJSON,
		Duration:   reqDuration,
		Payload:    data,
	}, nil
}
