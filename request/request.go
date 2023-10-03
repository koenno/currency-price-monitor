package request

import (
	"fmt"
	"io"
	"time"
)

type Currency struct {
	Name  string
	Rates []Rate
}

type Rate struct {
	Date  time.Time
	Value float64
}

func (r Rate) WriteTo(w io.Writer) (int64, error) {
	str := fmt.Sprintf("currency rate date=%v price=%v\n", r.Date, r.Value)
	n, err := io.WriteString(w, str)
	return int64(n), err
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
