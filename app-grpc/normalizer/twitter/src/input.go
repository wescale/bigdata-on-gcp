package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

func (s server) consumemessage() {

	var pull pubsub.PullRequest
	pull.Subscription = *subname
	pull.MaxMessages = 5
	ctx := context.Background()

	for {
		if ctx == nil {
			log.Println("Context is nil")
		}
		if s.sub == nil {
			log.Println("s.sub is nil")
		}
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
	if starttime, err := strconv.ParseInt(msg.Attributes["injector_time"], 10, 64); err != nil {
		log.Println(err)
	} else {
		var elapsedTime float64
		elapsedTime = float64(time.Now().Round(time.Millisecond).UnixNano() - starttime)
		s.timeProm.WithLabelValues("time").Observe(elapsedTime)

		var tweet twitter.Tweet
		tag := msg.Attributes["tag"]
		log.Println("Tag: " + tag)
		err := json.Unmarshal(msg.Data, &tweet)
		if err != nil {
			log.Println(err)
		} else {
			s.tweetStream <- (func() (twitter.Tweet, int64, string) { return tweet, starttime, tag })
		}
	}
}
