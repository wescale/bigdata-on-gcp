package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"


	language "cloud.google.com/go/language/apiv1"
)

var (
	addr    = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	topic   = flag.String("topic-name", os.Getenv("TOPIC_NAME"), "The topic listen.")
	subname = flag.String("sub-name", os.Getenv("SUB_NAME"), "the subscription write")
)

type server struct {
	pub         pubsub.PublisherClient
	sub         pubsub.SubscriberClient
	tweetStream chan func () (mastodon.Event, int64, string)
	msgStream   chan func () (libmetier.MessageSocial, int64)
	timeProm    *prometheus.HistogramVec
	language    *language.Client
	ctx         context.Context
}

func main() {

	flag.Parse()

	var s server

	s.ctx = context.Background()

	rand.Seed(time.Now().UnixNano())
	s.pub = connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.sub = connexionSubcriber("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")

	s.tweetStream = make(chan func()(mastodon.Event, int64, string))
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
