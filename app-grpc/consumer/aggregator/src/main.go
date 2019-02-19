package main

import (
	"github.com/go-redis/redis"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"context"

	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"

	"github.com/opentracing/opentracing-go"
)

// LoggerMiddleware add logger and metrics
func LoggerMiddleware(inner http.HandlerFunc, name string) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		inner.ServeHTTP(w, r)

		if strings.Compare(name,"health") != 0 {
			time := time.Since(start)
			log.Printf(
				"%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time,
			)
		}
	})
}

var (
	addr           = flag.String("listen-address", ":"+os.Getenv("PROM_PORT"), "The address to listen on for HTTP requests.")
	hashtag        = flag.String("hashtag", os.Getenv("HASHTAG"), "Twitter hashtag")
	projectid      = flag.String("project-id", os.Getenv("PROJECT_ID"), "Twitter hashtag")
	subname        = flag.String("sub-name", os.Getenv("SUB_NAME"), "Twitter hashtag")
	secretpath     = flag.String("secret-path", os.Getenv("SECRET_PATH"), "Twitter hashtag")
	redisaddr      = flag.String("redis-address", os.Getenv("REDIS_HOST")+":6379", "The address to listen on for HTTP requests.")
	aggregasub     = flag.String("aggrega-sub", os.Getenv("SUB_AGGREGA"), "subscription for cloud scheduler")
	pathprefix     = flag.String("path-prefix", os.Getenv("PATH_PREFIX"), "Path prefix")
	zipkinuri      = flag.String("zipkin-endpoint", os.Getenv("ZIPKIN_ENDPOINT"), "Zipkin endpoint")
)

type server struct {
	sub pubsub.SubscriberClient
	ds *datastore.Client
	timeProm *prometheus.HistogramVec
	redis *redis.Client
	ctx context.Context
}

func main() {

	flag.Parse()
	var s server

	// Define globals
	s.ctx = context.Background()


	///////////////////////////////// Zipkin Connection ////////////////////////////////
	collector, err := zipkin.NewHTTPCollector(*zipkinuri)
	if err != nil {
		log.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		os.Exit(-1)
	}

	// Create our recorder.
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:8080", "aggregator")

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
	s.ds = datastoreClient(s.ctx)
	s.timeProm = promHistogramVec()
	s.redis = redisNew()

	log.Println("launch consume thread")
	go s.consumeAggregatorMsg()

	router := mux.NewRouter().StrictSlash(true)

	var handlerStatus http.Handler
	handlerStatus = LoggerMiddleware(libmetier.HandlerStatusFunc, "root")
	router.
		Methods("GET").
		Path(*pathprefix + "/").
		Name("root").
		Handler(handlerStatus)

	var handlerHealth http.Handler
	handlerHealth = LoggerMiddleware(libmetier.HandlerHealthFunc, "health")
	router.
		Methods("GET").
		Path("/healthz").
		Name("health").
		Handler(handlerHealth)
	
	var handlerUsers http.Handler
	handlerUsers = LoggerMiddleware(s.handlerUsersFunc, "users_get")
	router.
		Methods("GET").
		Path(*pathprefix + "/users").
		Name("users_get").
		Handler(handlerUsers)


	var handlerTopTen http.Handler
	handlerTopTen = LoggerMiddleware(s.handlerTopTenFunc, "top_ten")
	router.
		Methods("GET").
		Path(*pathprefix + "/top10").
		Name("top_ten").
		Handler(handlerTopTen)


	var handlerStats http.Handler
	handlerStats = LoggerMiddleware(s.handlerStatsFunc, "stats")
	router.
		Methods("POST").
		Path(*pathprefix + "/stats").
		Name("stats").
		Handler(handlerStats)

	
	var handlerStatsID http.Handler
	handlerStatsID = LoggerMiddleware(s.handlerStatsIDFunc, "stats_id")
	router.
		Methods("GET").
		Path(*pathprefix + "/stats/{id}").
		Name("stats_id").
		Handler(handlerStatsID)

	router.Methods("GET").Path("/metrics").Name("Metrics").Handler(promhttp.Handler())

	// CORS
	headersOk := handlers.AllowedHeaders([]string{"authorization", "content-type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Printf("Start server...")
	http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
