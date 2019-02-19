package laststat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
)

// Aggrega test
type Aggrega struct {
	InjectorMean   float64   `json:"mean_time_injector" datastore:"mt_inj"`
	InjectorNb     int64     `json:"count_injector" datastore:"nb_inj"`
	NormalizerMean float64   `json:"mean_time_normalizer" datastore:"mt_nor"`
	NormalizerNb   int64     `json:"count_normalizer" datastore:"nb_nor"`
	Num            int64     `json:"id" datastore:"num"`
	CreateTime     time.Time `json:"create" datastore:"create_timestamp"`
}

type statStatus struct {
	Status  string  `json:"status" default:"done"`
	Elapsed int64   `json:"time_ms"`
	Agg     Aggrega `json:"result_ms"`
}

var (
	ds        *datastore.Client
	ctx       context.Context
	projectid = "slavayssiere-sandbox"
)

func init() {
	var err error
	ctx = context.Background()
	ds, err = datastore.NewClient(ctx, projectid)
	if err != nil {
		log.Printf("Failed to create client: %v\n", err)
	} else {
		log.Println("connected to client!")
	}
}

// LastStat test function
func LastStat(w http.ResponseWriter, r *http.Request) {

	ret := getLastStatID()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func getLastStatID() Aggrega {

	var err error
	var ret []Aggrega

	q := datastore.NewQuery("aggregas").Order("-create_timestamp").Limit(1)
	_, err = ds.GetAll(ctx, q, &ret)
	if err != nil {
		log.Printf("datastore: could not list Aggrega: %v", err)
	}

	return ret[0]
}
