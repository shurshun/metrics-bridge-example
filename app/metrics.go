package app

import (
	"net/http"
	"net/http/pprof"
	//	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func initMetrics(addr string) {
	r := http.NewServeMux()

	r.HandleFunc("/_liveness", livenessHandler)
	r.HandleFunc("/_readiness", readinessHandler)

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(addr, r)
}
