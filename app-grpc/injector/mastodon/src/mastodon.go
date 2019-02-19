package main

import (
	"context"
	"log"

	mastodon "github.com/mattn/go-mastodon"
)

func newMastodon(ctx context.Context) *mastodon.Client {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       *mserv,
		ClientID:     *mcid,
		ClientSecret: *mcsct,
	})

	err := client.Authenticate(ctx, *mlog, *mpasswd)
	if err != nil {
		log.Println(err)
	}

	return client
}

func (s server) reconnectStream(ctx context.Context) {
	for {
		log.Println("wait for error")
		test := <-s.streamError
		if test {
			log.Println("error receive")
			var err error
			s.timeline, err = s.mastodon.StreamingHashtag(ctx, *hashtag, false)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
