package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	RequestsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "total_requests_count",
	}, []string{})
	RequestsHitCache = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_cache_hit_count",
	})
	RequestsErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_error_count",
	}, []string{})

	SpiderTaskCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "spider_task_count",
	})
	SpiderErrorTaskCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "spider_error_task_count",
	})
)

func RegisterResponser() {
	prometheus.MustRegister(RequestsCount)
	prometheus.MustRegister(RequestsHitCache)
	prometheus.MustRegister(RequestsErrorCount)
}

func RegisterSpider() {
	prometheus.MustRegister(SpiderTaskCount)
	prometheus.MustRegister(SpiderErrorTaskCount)
}

func RegisterStorer() {

}

func RunPrometheus(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
