package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

func main(){
	m := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rc := recover(); rc != nil {
				ctx.WriteString(fmt.Sprintf("Something went wrong! %v", rc))
			}
		}()
		switch string(ctx.Path()) {
		case "/pay":
			pay(ctx)
		case "/metrics":
			handleMetrics(ctx)
		}
	}
	log.Print("Listening on port 8080...")
	fasthttp.ListenAndServe(":8080", m)
}

type PayRequest struct {
	Total float64
	RecordedCardKey string
	BankType string
}

func pay(ctx *fasthttp.RequestCtx) {
	request := &PayRequest{}
	err := json.Unmarshal(ctx.Request.Body(), request)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("Can't read body %v", err))
		return
	}

	start := time.Now()

	defer func() {
		// metric gauge decrease active payment requests.
		payRequestsActiveGauge.WithLabelValues(request.BankType).Dec()

		// metric summary for payment duration.
		paymentDurationSummary.WithLabelValues(request.BankType).Observe(time.Since(start).Seconds())
	}()

	// metric gauge increase active payment requests.
	payRequestsActiveGauge.WithLabelValues(request.BankType).Inc()

	if request.RecordedCardKey != "" {
		// metric counter for payments with recorded card count.
		paymentWithRecordedCardCounter.Inc()
	}

	switch request.BankType {
	case "A":
		payA(*request)
	case "B":
		payB(*request)
	}

	// metric histogram for payment completed with total price.
	paymentValueHistogram.Observe(request.Total)

	ctx.WriteString("OK")
}

func payA(request PayRequest) error {
	randomlyWait(50)
	return randomlyError(request.BankType)
}

func payB(request PayRequest) error {
	randomlyWait(120)
	return randomlyError(request.BankType)
}