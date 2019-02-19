package main

import (
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func newTwitter(consumerKey *string, consumerSecret *string, accessToken *string, accessSecret *string) *twitter.Client {
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	return twitter.NewClient(httpClient)
}

func (tc twitterClient) filterTwitter(hashtag string) *twitter.Stream {
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{hashtag},
		StallWarnings: twitter.Bool(true),
	}
	log.Printf("Starting Stream... for %s", hashtag)
	var err error
	strm, err := tc.clt.Streams.Filter(filterParams)
	if err != nil {
		log.Println(err)
	}

	return strm
}

func (s server) reconnectStream(tc *twitterClient) {
	for {
		log.Println("wait for error")
		test := <-s.streamError
		if test {
			log.Println("error receive")
			tc.strm.Stop()
			log.Println("Re-Create filter")
			tc.strm = (*tc).filterTwitter(*hashtag)
			go (*tc).demux.HandleChan(tc.strm.Messages)
		}
	}
}
