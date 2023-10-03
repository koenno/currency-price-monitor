package nbp

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnProperURLWithDefaultValues(t *testing.T) {
	// given
	domain := "something.com"
	client := NewCurrencyClient(domain)

	// when
	req, err := client.NewRequest(context.Background())

	// then
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, req.Method)
	assert.NoError(t, err)
	assert.Equal(t, domain, req.URL.Host)
	assert.Equal(t, "http", req.URL.Scheme)
	assert.Equal(t, "/api/exchangerates/rates/a/eur/last/1", req.URL.Path)
	query, err := url.ParseQuery(req.URL.RawQuery)
	assert.NoError(t, err)
	assert.Equal(t, string(FormatJSON), query.Get("format"))
}
