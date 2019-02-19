package libmetier

import (
	"encoding/json"
	"log"
	"time"
)

// AggregatedData a common social msg
type AggregatedData struct {
	User  string    `json:"user" datastore:"user"`
	Count int64     `json:"count" datastore:"count"`
	Date  time.Time `json:"timestamp" datastore:"date"`
}

// ListAggregatedData list aggregatedata
type ListAggregatedData *AggregatedData

func (ms AggregatedData) toAggregatedData(adtpl []byte) {
	err := json.Unmarshal(adtpl, &ms)
	if err != nil {
		log.Println(err)
	}
}

func (ms AggregatedData) toByteArray() []byte {
	b, err := json.Marshal(ms)
	if err != nil {
		log.Println(err)
	}
	return []byte(b)
}
