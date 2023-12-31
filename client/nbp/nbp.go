package nbp

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/text/currency"
)

type CurrencyClient struct {
	domain string
}

func NewCurrencyClient(domain string) CurrencyClient {
	return CurrencyClient{
		domain: domain,
	}
}

type Format string

const (
	FormatJSON Format = "json"
)

type options struct {
	historyInDays uint
	currencyUnit  currency.Unit
	format        Format
}

func defaultOptions() *options {
	return &options{
		historyInDays: 1,
		currencyUnit:  currency.EUR,
		format:        FormatJSON,
	}
}

type RequestOption func(*options)

func WithHistory(days uint) RequestOption {
	return func(o *options) {
		o.historyInDays = days
	}
}

func WithCurrency(unit currency.Unit) RequestOption {
	return func(o *options) {
		o.currencyUnit = unit
	}
}

func WithFormat(format Format) RequestOption {
	return func(o *options) {
		o.format = format
	}
}

// http://api.nbp.pl/api/exchangerates/rates/a/eur/last/100/?format=json
func (c CurrencyClient) NewRequest(ctx context.Context, opts ...RequestOption) (*http.Request, error) {
	cfg := defaultOptions()
	for _, o := range opts {
		o(cfg)
	}

	endpoint := "api/exchangerates/rates/a"
	rawURL := fmt.Sprintf("http://%s/%s/%s/last/%d",
		c.domain, endpoint, strings.ToLower(cfg.currencyUnit.String()), cfg.historyInDays)
	URL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create url: %v", err)
	}

	q := URL.Query()
	q.Set("format", string(cfg.format))
	URL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a request: %v", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "currency-price-monitor")

	return req, nil
}

type CurrencyResponse struct {
	Table    string  `json:"table"`
	Currency string  `json:"currency"`
	Code     string  `json:"code"`
	Rates    []Rates `json:"rates"`
}

type Rates struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}
