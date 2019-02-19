package getstat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// GetStat test function
func GetStat(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	ids := r.URL.Query().Get("id")

	var ret statStatus
	id, _ := strconv.ParseInt(ids, 10, 64)
	log.Printf("search stats n:%d", id)
	ret.Agg = getStatbyID(id)
	t := time.Now()
	ret.Elapsed = int64(t.Sub(start))
	ret.Status = "done"

	//convert to ms
	agg := ret.Agg
	agg.InjectorMean = agg.InjectorMean / 1000000
	agg.NormalizerMean = agg.NormalizerMean / 1000000
	ret.Agg = agg
	ret.Elapsed = ret.Elapsed / 1000000

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

func getStatbyID(id int64) Aggrega {

	var ret []Aggrega
	var err error

	q := datastore.NewQuery("aggregas").Filter("num=", id)
	_, err = ds.GetAll(ctx, q, &ret)
	if err != nil {
		log.Printf("datastore: could not list Aggrega: %v", err)
	}

	return ret[0]
}
