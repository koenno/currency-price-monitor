package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/koenno/currency-price-monitor/request"
	"golang.org/x/exp/slog"
)

var (
	ErrSendRequest     = errors.New("failed to send request")
	ErrResponse        = errors.New("response failure")
	ErrResponsePayload = errors.New("erroneus response payload")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func Request[T any](ctx context.Context, URL string) (request.Descriptor[T], error) {
	desc := request.Descriptor[T]{
		ID:   uuid.NewString(),
		URL:  URL,
		Time: time.Now(),
	}
	slog.Info("request", "id", desc.ID, "url", URL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return desc, fmt.Errorf("request failure: unable to create a request: %v", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return desc, fmt.Errorf("%w: %v", ErrSendRequest, err)
	}
	desc.Duration = time.Since(desc.Time)
	slog.Info("request", "id", desc.ID, "duration", desc.Duration)

	defer resp.Body.Close()
	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return desc, fmt.Errorf("%w: unable to read body: %v", ErrResponse, err)
	}

	desc.ValidStatusCode = resp.StatusCode == http.StatusOK
	slog.Info("request", "id", desc.ID, "validStatusCode", desc.ValidStatusCode)
	if resp.StatusCode != http.StatusOK {
		return desc, fmt.Errorf("%w: status code %d", ErrResponse, resp.StatusCode)
	}

	desc.JSON = resp.Header.Get("content-type") == "application/json"
	slog.Info("request", "id", desc.ID, "validContentType", desc.JSON)
	if !desc.JSON {
		return desc, fmt.Errorf("%w: unsupported %s", ErrResponsePayload, resp.Header.Get("content-type"))
	}

	desc.Valid = json.Valid(payloadBytes)
	slog.Info("request", "id", desc.ID, "validJson", desc.Valid)
	if !desc.Valid {
		return desc, fmt.Errorf("%w: invalid json", ErrResponsePayload)
	}

	err = json.Unmarshal(payloadBytes, &desc.Payload)
	if err != nil {
		return desc, fmt.Errorf("%w: unable to decode body to json %v", ErrResponse, err)
	}

	return desc, nil
}
