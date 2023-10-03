package main

import (
	"context"
	"io"
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

	logPath = "log.txt"
)

func main() {
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open a file %s: %v", logPath, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mainClient := client.New[nbp.CurrencyResponse](nbp.NewConverter())

	nbpClient := nbp.NewCurrencyClient(nbpDomain)
	nbpReq, err := nbpClient.NewRequest(ctx, nbp.WithCurrency(currency.EUR), nbp.WithHistory(100))
	if err != nil {
		log.Fatalf("failed to create NBP request: %v", err)
	}

	monitorSvc := monitor.New(mainClient, nbpReq)
	requestsPipe := monitorSvc.Start(ctx, requestsNo, requestsInterval)

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	writer := processor.NewWriter[nbp.CurrencyResponse](multiWriter)

	sched := scheduler.NewScheduler()
	sched.Register(writer)
	sched.Process(ctx, requestsPipe)
}
