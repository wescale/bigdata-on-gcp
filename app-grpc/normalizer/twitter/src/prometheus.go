package main

import "github.com/prometheus/client_golang/prometheus"

func getPromTime() *prometheus.HistogramVec {

	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "time_in_normalizer",
		Help: "Time for normalizer project in nanosecond",
	}, []string{"normalizer"})

	prometheus.Register(histogram)

	return histogram
}
