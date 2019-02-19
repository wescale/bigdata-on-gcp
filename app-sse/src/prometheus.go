package main

import "github.com/prometheus/client_golang/prometheus"

func promHistogramVec() *prometheus.HistogramVec {
	histogramMean := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "mean_in_sse",
			Help: "Time for pubish to pubsub in nanosecond",
		},
		[]string{"hashtag", "trade"},
	)

	prometheus.Register(histogramMean)

	return histogramMean
}
