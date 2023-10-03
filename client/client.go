package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

type Client[T any] struct {
}

func New[T any]() Client[T] {
	return Client[T]{}
}

func (c Client[T]) Process(req *http.Request) (request.Descriptor[T], error) {
	desc := request.Descriptor[T]{
		ID:   uuid.NewString(),
		URL:  req.URL.String(),
		Time: time.Now(),
	}
	slog.Info("request", "id", desc.ID, "url", desc.URL)

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
		return desc, fmt.Errorf("%w: status code %d; body %s", ErrResponse, resp.StatusCode, string(payloadBytes))
	}

	desc.JSON = strings.Contains(resp.Header.Get("content-type"), "application/json")
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
