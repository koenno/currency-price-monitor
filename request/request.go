package request

import (
	"fmt"
	"io"
	"time"
)

type Currency struct {
	Name  string
	Rates []Rates
}

type Rates struct {
	Date  time.Time
	Value float64
}

type Descriptor struct {
	ID              string
	URL             string
	Time            time.Time
	ValidStatusCode bool
	JSON            bool
	Valid           bool
	Duration        time.Duration
	Payload         Currency
}

func (d Descriptor) WriteTo(w io.Writer) (int64, error) {
	str := fmt.Sprintf("request id=%v url=%v time=%v validStatusCode=%v json=%v validJson=%v duration=%v\n",
		d.ID, d.URL, d.Time, d.ValidStatusCode, d.JSON, d.Valid, d.Duration)
	n, err := io.WriteString(w, str)
	return int64(n), err
}
