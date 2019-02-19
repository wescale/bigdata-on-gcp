package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"

	"context"
)

var (
	addr      = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	mserv     = flag.String("mastodon-server", os.Getenv("MASTODON_SERVER"), "Matodon auth.")
	mcid      = flag.String("mastodon-client-id", os.Getenv("MASTODON_CLIENT_ID"), "Matodon auth.")
	mcsct     = flag.String("mastodon-client-sct", os.Getenv("MASTODON_CLIENT_SECRET"), "Matodon auth.")
	mlog      = flag.String("mastodon-login", os.Getenv("MASTODON_LOGIN"), "Matodon auth.")
	mpasswd   = flag.String("mastodon-passwd", os.Getenv("MASTODON_PASSWORD"), "Matodon auth.")
	hashtag   = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	topicname = flag.String("topic", os.Getenv("TOPIC"), "Twitter hashtag")
)

type server struct {
	ps              pubsub.PublisherClient
	publishTimeChan chan int64
	timeInjectors   *prometheus.HistogramVec
	countInjectors  *prometheus.CounterVec
	streamError     chan bool
	mastodon        *mastodon.Client
	timeline        chan mastodon.Event
	ctx             context.Context
}

func main() {

	flag.Parse()
	var s server
	s.ctx = context.Background()

	// Client
	s.ps = s.connexionPublisher("pubsub.googleapis.com:443", os.Getenv("SECRET_PATH"), "https://www.googleapis.com/auth/pubsub")
	s.publishTimeChan = make(chan int64)
	s.streamError = make(chan bool)

	// mastondon
	s.mastodon = newMastodon(s.ctx)

	// Prometheus
	s.timeInjectors = PromHistogramVec()
	s.countInjectors = PromCounterVec()
	go func() {
		for {
			elapsed := <-s.publishTimeChan
			s.timeInjectors.WithLabelValues(os.Getenv("TOPIC_NAME")).Observe(float64(elapsed))
		}
	}()

	log.Println("Create filter")
	var err error
	s.timeline, err = s.mastodon.StreamingHashtag(s.ctx, *hashtag, false)
	if err != nil {
		log.Println("tesssssssssssssssst")
		log.Println(err)
	}

	go func() {
		for {
			e := <-s.timeline
			if _, ok := e.(*mastodon.ErrorEvent); !ok {
				s.publishmessage(e)
			} else {
				log.Println(e)
				s.streamError <- true
			}
		}
	}()

	// Receive messages until stopped or stream quits
	go s.reconnectStream(s.ctx)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Stopping Stream...")
		log.Println("Wait for 1 second to finish processing")
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	log.Println("launch server...")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
