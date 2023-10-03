package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/koenno/currency-price-monitor/client"
	"github.com/koenno/currency-price-monitor/client/nbp"
	"github.com/koenno/currency-price-monitor/monitor"
	"github.com/koenno/currency-price-monitor/processor"
	"github.com/koenno/currency-price-monitor/scheduler"
	"golang.org/x/text/currency"
)

const (
	nbpDomain        = "api.nbp.pl"
	requestsNo       = 10
	requestsInterval = 5 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mainClient := client.New[nbp.CurrencyResponse]()

	nbpClient := nbp.NewCurrencyClient(nbpDomain)
	nbpReq, err := nbpClient.NewRequest(ctx, nbp.WithCurrency(currency.EUR), nbp.WithHistory(100))
	if err != nil {
		log.Fatalf("failed to create NBP request: %v", err)
	}

	monitorSvc := monitor.New[nbp.CurrencyResponse](mainClient, nbpReq)
	requestsPipe := monitorSvc.Start(ctx, requestsNo, requestsInterval)

	stdoutWriter := processor.NewWriter[nbp.CurrencyResponse](os.Stdout)

	sched := scheduler.NewScheduler[nbp.CurrencyResponse]()
	sched.Register(stdoutWriter)
	sched.Process(ctx, requestsPipe)
}
