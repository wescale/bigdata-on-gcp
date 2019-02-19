package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	"github.com/opentracing/opentracing-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

func (s server) connexionPublisher(tracer opentracing.Tracer, address string, filename string, scope ...string) pubsub.PublisherClient {
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")
	log.Printf("Secret in %s\n", filename)
	perRPC, err := oauth.NewServiceAccountFromFile(filename, "https://www.googleapis.com/auth/pubsub")
	if err != nil {
		log.Println(err)
	}

	conn, err := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		log.Println(err)
	}

	return pubsub.NewPublisherClient(conn)
}

func (s server) publishmessage(tweet *twitter.Tweet, publishTime chan int64) {
	var message pubsub.PubsubMessage
	var request pubsub.PublishRequest

	start := time.Now()

	log.Println(tweet)

	b, err := json.Marshal(tweet)
	if err != nil {
		log.Println(err)
	}
	ctx := context.Background()
	message.Data = []byte(b)
	message.Attributes = make(map[string]string)
	message.Attributes["source"] = "twitter"
	message.Attributes["tag"] = *hashtag
	message.Attributes["injector_time"] = strconv.FormatInt(start.UnixNano(), 10)

	request.Topic = *topicname
	log.Println("send message to " + *topicname)
	request.Messages = append(request.Messages, &message)

	if _, err := s.ps.Publish(ctx, &request); err != nil {
		log.Println(err)
	}

	t := time.Now()
	elapsed := t.Sub(start)

	publishTime <- elapsed.Nanoseconds()
}
