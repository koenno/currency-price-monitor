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
)

var (
	ErrSendRequest     = errors.New("failed to send request")
	ErrResponse        = errors.New("response failure")
	ErrResponsePayload = errors.New("erroneus response payload")

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

//go:generate mockery --name=Converter --case underscore --with-expecter
type Converter[T any] interface {
	Convert(from T) request.Currency
}

type Client[T any] struct {
	converter Converter[T]
}

func New[T any](converter Converter[T]) Client[T] {
	return Client[T]{
		converter: converter,
	}
}

func (c Client[T]) Process(req *http.Request) (request.Descriptor, error) {
	desc := request.Descriptor{
		ID:   uuid.NewString(),
		URL:  req.URL.String(),
		Time: time.Now(),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return desc, fmt.Errorf("%w: %v", ErrSendRequest, err)
	}
	desc.Duration = time.Since(desc.Time)

	defer resp.Body.Close()
	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return desc, fmt.Errorf("%w: unable to read body: %v", ErrResponse, err)
	}

	desc.ValidStatusCode = resp.StatusCode == http.StatusOK
	if resp.StatusCode != http.StatusOK {
		return desc, fmt.Errorf("%w: status code %d; body %s", ErrResponse, resp.StatusCode, string(payloadBytes))
	}

	desc.JSON = strings.Contains(resp.Header.Get("content-type"), "application/json")
	if !desc.JSON {
		return desc, fmt.Errorf("%w: unsupported %s", ErrResponsePayload, resp.Header.Get("content-type"))
	}

	desc.Valid = json.Valid(payloadBytes)
	if !desc.Valid {
		return desc, fmt.Errorf("%w: invalid json", ErrResponsePayload)
	}

	var payload T
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return desc, fmt.Errorf("%w: unable to decode body to json %v", ErrResponse, err)
	}

	desc.Payload = c.converter.Convert(payload)

	return desc, nil
}
