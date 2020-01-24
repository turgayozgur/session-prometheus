package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/pay", pay)

	log.Print("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

type PayRequest struct {
	Total           float64
	RecordedCardKey string
	BankType        string
}

func pay(w http.ResponseWriter, r *http.Request) {
	request := &PayRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()

	defer func() {
		if rc := recover(); rc != nil {
			// metric counter increment failed payment requests.
			paymentRequestFailedCounter.WithLabelValues(request.BankType).Inc()
		}

		// metric gauge decrease active payment requests.
		payRequestsActiveGauge.WithLabelValues(request.BankType).Dec()

		// metric summary for payment duration.
		paymentDurationSummary.WithLabelValues(request.BankType).Observe(time.Since(start).Seconds())
	}()

	// metric gauge increase active payment requests.
	payRequestsActiveGauge.WithLabelValues(request.BankType).Inc()

	if request.RecordedCardKey != "" {
		// metric counter for payments with recorded card count.
		paymentWithRecordedCardCounter.WithLabelValues(request.BankType).Inc()
	}

	switch request.BankType {
	case "A":
		payA(*request)
	case "B":
		payB(*request)
	}

	// metric histogram for payment completed with total price.
	paymentValueHistogram.Observe(request.Total)

	fmt.Fprintf(w, "Success")
}

func payA(request PayRequest) {
	randomlyWait(50)
	randomlyError(request.BankType)
}

func payB(request PayRequest) {
	randomlyWait(120)
	randomlyError(request.BankType)
}
