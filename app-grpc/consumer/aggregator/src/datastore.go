package main

import (
	"context"
	"io/ioutil"
	"log"

	"cloud.google.com/go/datastore"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func datastoreClient(ctx context.Context) *datastore.Client {
	jsonKey, err := ioutil.ReadFile(*secretpath)
	config, err := google.JWTConfigFromJSON(jsonKey, datastore.ScopeDatastore) // or bigtable.AdminScope, etc.
	client, err := datastore.NewClient(ctx, *projectid, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		log.Printf("Failed to create client: %v\n", err)
	} else {
		log.Println("connected to client!")
	}

	return client
}

func (s server) writeAggrega(table string, agg Aggrega) {
	// Saves the new entity.
	k := datastore.IncompleteKey(table, nil)
	if _, err := s.ds.Put(s.ctx, k, &agg); err != nil {
		log.Printf("Failed to save task: %v\n", err)
	}
}

func (s server) getStatbyID(id int64) Aggrega{
	var ret []Aggrega

	q := datastore.NewQuery("aggregas").Filter("num=",id)
	_, err := s.ds.GetAll(s.ctx, q, &ret)
	if err != nil {
		log.Printf("datastore: could not list Aggrega: %v", err)
	}

	return ret[0]
}

func (s server) writeBulkMessage(ads []libmetier.AggregatedData) {
	// Saves the new entity.
	kl := make([]*datastore.Key, len(ads))
	k := datastore.IncompleteKey("userstats", nil)
	for i:=0; i!=len(ads); i++ {
		kl[i]=k
	}
	if len(ads) > 500 {
		var min int
		var max int
		min=0
		max=400
		for {
			log.Printf("min: %d, &max: %d, &len(ads): %d", min, max, len(ads))
			adstmp := ads[min:max]
			kltmp := kl[min:max]
			if _, err := s.ds.PutMulti(s.ctx, kltmp, adstmp); err != nil {
				log.Printf("Failed to save task: %v\n", err)
			}
			min = max+1
			if max > (len(ads) - 400) {
				max = len(ads)
			} else {
				max = max + 400
			}
			if min > len(ads) {
				break
			}
		}
	} else {
		if _, err := s.ds.PutMulti(s.ctx, kl, ads); err != nil {
			log.Printf("Failed to save task: %v\n", err)
		}
	}
}

func (s server) readMessage() []libmetier.AggregatedData {
	var uss []libmetier.AggregatedData
	
	q := datastore.NewQuery("userstats").Order("count").Limit(10)
	_, err := s.ds.GetAll(s.ctx, q, &uss)
	if err != nil {
		log.Printf("datastoredb: could not list books: %v", err)
		return nil
	}

	return uss
}
