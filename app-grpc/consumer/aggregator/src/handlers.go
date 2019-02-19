package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (s server) handlerUsersFunc(w http.ResponseWriter, r *http.Request) {
	us := s.getUsersCounterList(100)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

func (s server) handlerTopTenFunc(w http.ResponseWriter, r *http.Request) {
	us := s.top10()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(us); err != nil {
		panic(err)
	}
}

type statStatus struct {
	Status  string  `json:"status" default:"done"`
	Elapsed int64   `json:"time_ms"`
	Agg     Aggrega `json:"result_ms"`
}

func (s server) handlerStatsFunc(w http.ResponseWriter, r *http.Request) {
	var ret statStatus

	start := time.Now()
	ret.Agg = s.computeAggregas()
	s.writeAggrega("aggregas", ret.Agg)
	ret.Elapsed = int64(ret.Agg.CreateTime.Sub(start))
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

func (s server) handlerStatsIDFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]

	start := time.Now()

	var ret statStatus
	id, _ := strconv.ParseInt(ids, 10, 64)
	log.Printf("search stats n:%d", id)
	ret.Agg = s.getStatbyID(id)
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
