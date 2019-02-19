package main

import (
	"context"
	"log"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func newNaturalLanguage(ctx context.Context) *language.Client {

	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func (s server) analyzeText(txt string) float32 {
	// Detects the sentiment of the text.
	sentiment, err := s.language.AnalyzeSentiment(s.ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: txt,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Printf("Failed to analyze text: %v", err)
	}

	var ret float32
	ret = 0.0
	if sentiment != nil {
		if sentiment.DocumentSentiment != nil {
			ret = sentiment.DocumentSentiment.Score
		}
	}

	return ret
}

func (s server) convert() {
	for {
		tweet, starttime, tag := (<-s.tweetStream)()
		var u libmetier.MessageSocial
		u.Data = tweet.Text
		u.User = tweet.User.Name
		u.Source = "twitter"
		u.Tag = tag
		u.Date = time.Unix(0, starttime)
		u.ID = tweet.IDStr
		t := time.Now()
		if s.el == true {
			u.Sentiment = s.analyzeText(tweet.Text)
		} else {
			u.Sentiment = 0.0
		}
		log.Println(time.Now().Sub(t).Seconds())
		s.msgStream <- (func() (libmetier.MessageSocial, int64) { return u, starttime })
	}
}
