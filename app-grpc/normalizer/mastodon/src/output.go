package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	pubsub "google.golang.org/genproto/googleapis/pubsub/v1beta2"
)

func (s server) sendMessage() {

	for {
		log.Println("Wait for msgSTream...")
		msg, starttime := (<-s.msgStream)()

		var message pubsub.PubsubMessage
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
		message.Data = []byte(b)
		message.Attributes = make(map[string]string)
		message.Attributes["injector_time"] = strconv.FormatInt(starttime, 10)
		message.Attributes["normalizer_time"] = strconv.FormatInt(time.Now().UnixNano(), 10)

		s.publishmessage(&message)
	}
}
