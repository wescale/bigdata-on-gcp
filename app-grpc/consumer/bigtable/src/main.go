package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigtable"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	"context"

	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/opentracing/opentracing-go"
)

var (
	addr           = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	hashtag        = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	projectid      = flag.String("project-id", os.Getenv("PROJECT_ID"), "Twitter hashtag")
	instanceid     = flag.String("instance-id", os.Getenv("INSTANCE_ID"), "Twitter hashtag")
	tableid        = flag.String("table-id", os.Getenv("TABLE_ID"), "Twitter hashtag")
	subname        = flag.String("sub-name", os.Getenv("SUB_NAME"), "Twitter hashtag")
	secretpath     = flag.String("secret-path", os.Getenv("SECRET_PATH"), "Twitter hashtag")
	aggregasub     = flag.String("aggrega-sub", os.Getenv("SUB_AGGREGA"), "subscription for cloud scheduler")
	zipkinuri      = flag.String("zipkin-endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "Zipkin endpoint")
)

type server struct {
	sub pubsub.SubscriberClient
	bt bigtable.Client
	messages chan libmetier.MessageSocial
	timeProm *prometheus.HistogramVec
}

func main() {

	flag.Parse()
	var s server

	// Define globals
	ctx := context.Background()


	///////////////////////////////// Zipkin Connection ////////////////////////////////
	collector, err := zipkin.NewHTTPCollector(*zipkinuri)
	if err != nil {
		log.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "bigtable")

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

	log.Println("Get secret from: " + *secretpath)
	s.sub = s.connexionSubcriber(tracer, "pubsub.googleapis.com:443", *secretpath, "https://www.googleapis.com/auth/pubsub")
	s.bt = bigtableClient(ctx)
	s.messages = make(chan libmetier.MessageSocial)
	s.timeProm = promHistogramVec()

	log.Println("launch consume thread")
	go s.consumemessage()
	go s.consumeAggregatorMsg()

	log.Println("write in bigtable")
	go s.writeMessages(ctx)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
