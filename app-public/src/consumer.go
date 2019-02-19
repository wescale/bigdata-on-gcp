package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"

	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"github.com/opentracing/opentracing-go"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

func connexionSubcriber(ctx context.Context, tracer opentracing.Tracer, subname string, address string, filename string, scope ...string) pubsub.SubscriberClient {
	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Println(err)
	}

	creds := credentials.NewClientTLSFromCert(pool, "")
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

	client := pubsub.NewSubscriberClient(conn)

	//create specific subscription
	log.Printf("create subscription: %s", subname)
	_, err = client.CreateSubscription(ctx, &pubsub.Subscription{
		Name:  subname,
		Topic: *topicName,
	})
	if err != nil {
		log.Println(err)
	}

	return client
}

func (s server) closeSubscription() {
	log.Printf("delete subscription: %s", s.sub)
	s.clt.DeleteSubscription(s.ctx, &pubsub.DeleteSubscriptionRequest{
		Subscription: s.sub,
	})
}

func (s server) consumemessage() {

	var pull pubsub.PullRequest
	pull.Subscription = s.sub
	pull.MaxMessages = 5

	if s.ctx == nil {
		log.Println("Context is nil")
	}
	if s.clt == nil {
		log.Println("s.sub is nil")
	}

	for {
		if resp, err := s.clt.Pull(s.ctx, &pull); err != nil {
			fmt.Println(err)
		} else {
			s.messagesreceive(resp, pull)
		}
	}
}

func (s server) messagesreceive(resp *pubsub.PullResponse, pull pubsub.PullRequest) {
	var ackMess pubsub.AcknowledgeRequest
	ackMess.Subscription = pull.Subscription
	for _, messRec := range resp.ReceivedMessages {
		ackMess.AckIds = append(ackMess.AckIds, messRec.GetAckId())
		s.msgreceive(messRec.GetMessage())
	}
	s.clt.Acknowledge(s.ctx, &ackMess)
}

func (s server) msgreceive(msg *pubsub.PubsubMessage) {
	log.Println(msg.Data)
	var ms libmetier.MessageSocial
	err := json.Unmarshal(msg.Data, &ms)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(ms)
		s.messages <- ms
	}
}
