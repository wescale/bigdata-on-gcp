package main

import "github.com/prometheus/client_golang/prometheus"

func promHistogramVec() *prometheus.HistogramVec {
	histogramMean := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "mean_in_injector",
			Help: "Time for pubish to pubsub in nanosecond",
		},
		[]string{"hashtag", "trade"},
	)

	prometheus.Register(histogramMean)

	return histogramMean
}

func promCounterVec() *prometheus.CounterVec {
	messagesCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_injected",
			Help: "How many messages injected, partitioned by hashtag",
		},
		[]string{"hashtag"},
	)

	prometheus.Register(messagesCounter)

	return messagesCounter
}
