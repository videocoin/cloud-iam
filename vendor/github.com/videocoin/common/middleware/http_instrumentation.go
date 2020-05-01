package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InstrumentHandler(reg prometheus.Registerer, h http.Handler) http.Handler {
	cnt := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests by HTTP status code.",
		},
		[]string{"code"},
	)

	cnt.WithLabelValues("200")
	cnt.WithLabelValues("500")
	cnt.WithLabelValues("503")
	if err := reg.Register(cnt); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			cnt = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			panic(err)
		}
	}

	gge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of requests being served.",
	})
	if err := reg.Register(gge); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			gge = are.ExistingCollector.(prometheus.Gauge)
		} else {
			panic(err)
		}
	}

	return promhttp.InstrumentHandlerCounter(cnt, promhttp.InstrumentHandlerInFlight(gge, h))
}
