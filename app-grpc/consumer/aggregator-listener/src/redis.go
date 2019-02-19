package main

import (
	"log"
	"time"

	"github.com/go-redis/redis"
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

func redisNew() *redis.Client {
	var client *redis.Client
	for {
		client = redis.NewClient(&redis.Options{
			Addr:     *redisaddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err := client.Ping().Result()
		if err != nil {
			log.Println("error in redis connection to " + *redisaddr)
			log.Println(err)
		} else {
			log.Println("MemoryStore connected !")
			break
		}
	}

	return client
}

func (s server) countUser(tag string, user string) {
	s.redis.Incr(user + "_" + tag)
	s.redis.SAdd("list_users_"+tag, user)
}

func (s server) addNormTime(tag string, normtime int64) {
	s.redis.LPush("normTimes_"+tag+"_"+string(s.getNbAggregation()), normtime)
}

func (s server) addInjectTime(tag string, injectime int64) {
	s.redis.LPush("injectTimes_"+tag+"_"+string(s.getNbAggregation()), injectime)
}

func (s server) addAggTime(tag string, aggtime int64) {
	s.redis.LPush("aggTimes_"+tag+"_"+string(s.getNbAggregation()), aggtime)
}

func (s server) getNbAggregation() int64 {
	nb, err := s.redis.Get("aggregas").Int64()
	if err != nil {
		log.Println(err)
	}
	return nb
}

func (s server) addTag(tag string) {
	s.redis.SAdd("tag_list", tag)
}

func (s server) writeMessagesToRedis() {

	for {
		mess, normtime, injectime := (<-s.messages)()
		s.addTag(mess.Tag)
		if len(mess.User) > 0 {
			s.countUser(mess.Tag, mess.User)
		}
		s.addNormTime(mess.Tag, time.Now().UnixNano()-normtime)
		s.addInjectTime(mess.Tag, time.Now().UnixNano()-injectime)
	}
}
