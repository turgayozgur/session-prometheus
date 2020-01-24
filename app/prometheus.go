package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	payRequestsActiveGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "session_prom_pay_requests_active",
			Help: "Number of pay requests processing.",
		}, []string{"bank_type"})
	paymentWithRecordedCardCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_prom_payment_with_recorded_card_total",
			Help: "Number of payment completed with recorded card.",
		}, []string{"bank_type"})
	paymentDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "session_prom_payment_duration_seconds",
			Help:       "Summary of payment duration seconds over last 10 minutes.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"bank_type"})
	paymentValueHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "session_prom_payment_value_tl",
			Help:    "Histogram of received payment values (in TL)",
			Buckets: []float64{20, 100, 200, 350, 500, 1000},
		})
	paymentRequestFailedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "session_prom_payment_failed_total",
			Help: "Number of payment failed.",
		}, []string{"bank_type"})
)

func init() {
	r := prometheus.NewRegistry()

	r.MustRegister(payRequestsActiveGauge)
	r.MustRegister(paymentWithRecordedCardCounter)
	r.MustRegister(paymentDurationSummary)
	r.MustRegister(paymentValueHistogram)
	r.MustRegister(paymentRequestFailedCounter)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}
