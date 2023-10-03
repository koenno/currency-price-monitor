package nbp

import (
	"testing"
	"time"

	"github.com/koenno/currency-price-monitor/request"
	"github.com/stretchr/testify/assert"
)

func TestShouldConvertNBPData(t *testing.T) {
	// given
	toConvert := CurrencyResponse{
		Table:    "A",
		Currency: "dolar amerykański",
		Code:     "USD",
		Rates: []Rates{
			{
				No:            "092/A/NBP/2023",
				EffectiveDate: "2023-05-15",
				Mid:           4.1490,
			},
			{
				No:            "093/A/NBP/2023",
				EffectiveDate: "2023-05-16",
				Mid:           4.1228,
			},
		},
	}
	expectedRates := []request.Rates{
		{
			Date:  newDate(toConvert.Rates[0].EffectiveDate),
			Value: toConvert.Rates[0].Mid,
		},
		{
			Date:  newDate(toConvert.Rates[1].EffectiveDate),
			Value: toConvert.Rates[1].Mid,
		},
	}
	sut := NewConverter()

	// when
	converted := sut.Convert(toConvert)

	// then
	assert.Equal(t, toConvert.Code, converted.Name)
	assert.ElementsMatch(t, expectedRates, converted.Rates)
}

func newDate(date string) time.Time {
	t, _ := time.Parse(time.DateOnly, date)
	return t
}
