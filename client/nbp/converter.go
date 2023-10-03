package nbp

import (
	"time"

	"log"

	"github.com/koenno/currency-price-monitor/request"
)

type Converter struct {
}

func NewConverter() Converter {
	return Converter{}
}

func (c Converter) Convert(from CurrencyResponse) request.Currency {
	result := request.Currency{
		Name:  from.Code,
		Rates: []request.Rate{},
	}
	for _, fromRate := range from.Rates {
		date, err := time.Parse(time.DateOnly, fromRate.EffectiveDate)
		if err != nil {
			log.Printf("failed to convert date %s: %v", fromRate.EffectiveDate, err)
			continue
		}
		result.Rates = append(result.Rates, request.Rate{
			Date:  date,
			Value: fromRate.Mid,
		})
	}
	return result
}
