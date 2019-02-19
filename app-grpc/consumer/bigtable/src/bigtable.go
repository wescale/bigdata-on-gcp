package main

import (
	"context"
	"encoding/binary"
	"io/ioutil"
	"log"
	"time"
	"fmt"

	"cloud.google.com/go/bigtable"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"crypto/sha256"
)

// User-provided constants.
const (
	columnFamilyName = "ms"

	columnNameData   = "data"
	columnNameUser   = "user"
	columnNameSource = "source"
	columnTagSource  = "tag"
	columnNameDate   = "time"
	columnSentiment  = "sentiment"
	columnID         = "ID"
)

// sliceContains reports whether the provided string is present in the given slice of strings.
func sliceContains(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}

func bigtableClient(ctx context.Context) bigtable.Client {
	jsonKey, err := ioutil.ReadFile(*secretpath)
	config, err := google.JWTConfigFromJSON(jsonKey, bigtable.Scope) // or bigtable.AdminScope, etc.

	client, err := bigtable.NewClient(ctx, *projectid, *instanceid, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		log.Fatalf("Could not create data operations client: %v", err)
	}

	return *client
}

func (s server) writeMessage(ctx context.Context, mess libmetier.MessageSocial) {

	tbl := s.bt.Open(*tableid)

	// Mutation Way
	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, columnNameData, bigtable.Now(), []byte(mess.Data))
	mut.Set(columnFamilyName, columnNameUser, bigtable.Now(), []byte(mess.User))
	mut.Set(columnFamilyName, columnNameSource, bigtable.Now(), []byte(mess.Source))
	mut.Set(columnFamilyName, columnTagSource, bigtable.Now(), []byte(mess.Tag))
	mut.Set(columnFamilyName, columnID, bigtable.Now(), []byte(mess.ID))
	mut.Set(columnFamilyName, columnNameDate, bigtable.Now(),[]byte(mess.Date.Format(time.RFC3339Nano)))
	mut.Set(columnFamilyName, columnSentiment, bigtable.Now(),[]byte(fmt.Sprintf("%f", mess.Sentiment)))

	sha256 := sha256.Sum256([]byte(mess.User))

	key := fmt.Sprintf("%s-%s-%x-%d\n", mess.Source, mess.Tag, sha256, mess.Date.UnixNano())

	if err := tbl.Apply(ctx, key, mut); err != nil {
		log.Println(err)
	}
}

func (s server) readMessage(ctx context.Context) libmetier.MessageSocial {

	var mess libmetier.MessageSocial

	tbl := s.bt.Open(*tableid)

	rowKey := "test"

	row, err := tbl.ReadRow(ctx, rowKey)
	if err != nil {
		log.Fatalf("Could not read row with key %s: %v", rowKey, err)
	}
	log.Printf("Row key: %s\n", rowKey)
	mess.Data = string(row[columnFamilyName][0].Value)
	mess.User = string(row[columnFamilyName][1].Value)
	mess.Source = string(row[columnFamilyName][2].Value)
	mess.Date = time.Unix(0, int64(binary.LittleEndian.Uint64(row[columnFamilyName][3].Value)))
	log.Println("Data:", mess.Data)
	log.Println("Source:", mess.Source)
	log.Println("User:", mess.User)
	log.Println("Date:", mess.Date)

	return mess
}

func (s server) writeMessages(ctx context.Context) {

	for {
		mess := <-s.messages
		s.writeMessage(ctx, mess)
	}
}
