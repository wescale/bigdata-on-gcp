package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"

	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/opentracing/opentracing-go"
)

var (
	addr           = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	consumerKey    = flag.String("consumer-key", os.Getenv("CONSUMER_KEY"), "Twitter Consumer Key")
	consumerSecret = flag.String("consumer-secret", os.Getenv("CONSUMER_SECRET"), "Twitter Consumer Secret")
	accessToken    = flag.String("access-token", os.Getenv("ACCESS_TOKEN"), "Twitter Access Token")
	accessSecret   = flag.String("access-secret", os.Getenv("ACCESS_SECRET"), "Twitter Access Secret")
	hashtag        = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	topicname      = flag.String("topic-name", os.Getenv("TOPIC_NAME"), "Twitter hashtag")
	zipkinuri      = flag.String("zipkin-endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "Zipkin endpoint")
)

type server struct {
	ps              pubsub.PublisherClient
	publishTimeChan chan int64
	timeInjectors   *prometheus.HistogramVec
	countInjectors  *prometheus.CounterVec
	streamError     chan bool
}

type twitterClient struct {
	clt    *twitter.Client
	strm   *twitter.Stream
	demux  twitter.SwitchDemux
	Filter []string
}

func main() {

	flag.Parse()

	var s server
	var tc twitterClient

	///////////////////////////////// Zipkin Connection ////////////////////////////////
	collector, err := zipkin.NewHTTPCollector(*zipkinuri)
	if err != nil {
		log.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "injector-twitter")

	// Create our tracer.
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		log.Printf("unable to create Zipkin tracer: %+v\n", err)
		os.Exit(-1)
	}

	// Explicitly set our tracer to be the default tracer.
	opentracing.InitGlobalTracer(tracer)

	///////////////////////////////// Pubsub Connection ////////////////////////////////
	s.ps = s.connexionPublisher(tracer, "pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.publishTimeChan = make(chan int64)
	s.streamError = make(chan bool)

	///////////////////////////////// Twitter Connection ////////////////////////////////
	tc.clt = newTwitter(consumerKey, consumerSecret, accessToken, accessSecret)

	// Prometheus
	s.timeInjectors = promHistogramVec()
	s.countInjectors = promCounterVec()
	go func() {
		for {
			elapsed := <-s.publishTimeChan
			s.timeInjectors.WithLabelValues(*hashtag, os.Getenv("TOPIC_NAME")).Observe(float64(elapsed))
		}
	}()

	tc.demux = twitter.NewSwitchDemux()
	tc.demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet != nil {
			s.publishmessage(tweet, s.publishTimeChan)
			s.countInjectors.WithLabelValues(*hashtag).Add(1)
		} else {
			log.Printf("Tweet null")
		}
	}

	tc.demux.Warning = func(warning *twitter.StallWarning) {
		if warning != nil {
			log.Println("Warning")
			log.Println(warning)
			s.streamError <- true
		} else {
			log.Printf("warning null")
		}
	}

	tc.demux.StreamDisconnect = func(disconnect *twitter.StreamDisconnect) {
		if disconnect != nil {
			log.Println("StreamDisconnect")
			log.Println(disconnect)
			s.streamError <- true
		} else {
			log.Printf("disconnect null")
		}
	}

	tc.demux.Other = func(message interface{}) {
		if message != nil {
			log.Println("Other")
			log.Println(message)
			s.streamError <- true
		} else {
			log.Printf("message null")
		}
	}

	log.Println("Create filter")
	tc.strm = tc.filterTwitter(*hashtag)
	go tc.demux.HandleChan(tc.strm.Messages)

	// Receive messages until stopped or stream quits
	go s.reconnectStream(&tc)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Stopping Stream...")
		tc.strm.Stop()
		log.Println("Wait for 1 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	log.Println("launch server...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
