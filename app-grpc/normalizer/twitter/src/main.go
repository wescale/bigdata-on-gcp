package main

import (
	"strconv"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	language "cloud.google.com/go/language/apiv1"

	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/opentracing/opentracing-go"
)

var (
	addr      = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	topic     = flag.String("topic-name", os.Getenv("TOPIC_NAME"), "The topic listen.")
	subname   = flag.String("sub-name", os.Getenv("SUB_NAME"), "the subscription write")
	els       = flag.String("enable-language", os.Getenv("LANGUAGE"), "enable the language-sentiment call")
	zipkinuri = flag.String("zipkin-endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "Zipkin endpoint")
)

type server struct {
	pub         pubsub.PublisherClient
	sub         pubsub.SubscriberClient
	tweetStream chan func () (twitter.Tweet, int64, string)
	msgStream   chan func () (libmetier.MessageSocial, int64)
	timeProm    *prometheus.HistogramVec
	language    *language.Client
	ctx         context.Context
	el bool
}

func main() {

	flag.Parse()

	var s server
	var err error

	s.el, err = strconv.ParseBool(*els)
	if err != nil {
		log.Println(err)
		s.el = false
	}

	s.ctx = context.Background()

	///////////////////////////////// Zipkin Connection ////////////////////////////////
	collector, err := zipkin.NewHTTPCollector(*zipkinuri)
	if err != nil {
		log.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "normalizer-twitter")

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


	rand.Seed(time.Now().UnixNano())
	s.pub = connexionPublisher(tracer, "pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.sub = connexionSubcriber(tracer, "pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	s.tweetStream = make(chan func()(twitter.Tweet, int64, string))
	s.msgStream = make(chan func()(libmetier.MessageSocial, int64))

	s.timeProm = getPromTime()
	s.language = newNaturalLanguage(s.ctx)

	log.Println("launch converter thread")
	go s.convert()

	log.Println("launch consume thread")
	go s.consumemessage()

	log.Println("launch send thread")
	go s.sendMessage()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Wait for 1 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
