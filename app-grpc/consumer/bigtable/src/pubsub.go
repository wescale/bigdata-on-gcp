package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"log"

	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"github.com/opentracing/opentracing-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

func (s server) connexionSubcriber(tracer opentracing.Tracer, address string, filename string, scope ...string) pubsub.SubscriberClient {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := oauth.NewServiceAccountFromFile(filename, scope...)
	conn, _ := grpc.Dial(
		"pubsub.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)

	return pubsub.NewSubscriberClient(conn)
}

func (s server) consumemessage() {

	var pull pubsub.PullRequest
	pull.Subscription = *subname
	pull.MaxMessages = 5
	ctx := context.Background()

	if ctx == nil {
		log.Println("Context is nil")
	}
	if s.sub == nil {
		log.Println("s.sub is nil")
	}

	for {
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			log.Println(err)
		} else {
			s.messagesreceive(ctx, resp, pull)
		}
	}
}

func (s server) messagesreceive(ctx context.Context, resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.msgreceive(messRec.GetMessage())
	}
	s.sub.Acknowledge(ctx, &ackMess)
}

func (s server) msgreceive(msg *pubsub.PubsubMessage) {
	var ms libmetier.MessageSocial
	err := json.Unmarshal(msg.Data, &ms)
	if err != nil {
		log.Println(err)
	} else {
		s.messages <- ms
	}
}

// MessageAggregator a message send by cloud sheduler
type MessageAggregator struct {
	Script string `json:"run"`
}

func (s server) consumeAggregatorMsg() {

	var pull pubsub.PullRequest
	pull.Subscription = *aggregasub
	pull.MaxMessages = 5
	ctx := context.Background()

	if ctx == nil {
		log.Println("Context is nil")
	}
	if s.sub == nil {
		log.Println("s.sub is nil")
	}

	for {
		resp, err := s.sub.Pull(ctx, &pull)
		if err != nil {
			log.Println(err)
		} else {
			s.messageAggregatorreceive(ctx, resp, pull)
		}
	}
}

func (s server) messageAggregatorreceive(ctx context.Context, resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.mAreceive(messRec.GetMessage())
	}
	s.sub.Acknowledge(ctx, &ackMess)
}

func (s server) mAreceive(msg *pubsub.PubsubMessage) {
	var ma MessageAggregator
	err := json.Unmarshal(msg.Data, &ma)
	if err != nil {
		log.Println(err)
	} else {
		if ma.Script == "dataset" {
			log.Println("dataset generator")
		} else {
			log.Println("aggrega generator")
		}
	}
}
