package main

import (
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/slavayssiere/sandbox-gcp/app-grpc/libmetier"
)

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

// UserCounter UserCounter
type UserCounter struct {
	Users []libmetier.AggregatedData `json:"users"`
	Tag   string                     `json:"tag"`
}

func (s server) getUsersCounter(tag string, limit int) []libmetier.AggregatedData {
	var ret []libmetier.AggregatedData
	users, err := s.redis.SMembers("list_users_" + tag).Result()
	if err != nil {
		log.Println(err)
	}
	for id := range users {
		var user libmetier.AggregatedData
		user.User = users[id]
		user.Count, err = s.redis.Get(users[id] + "_" + tag).Int64()
		if err != nil {
			log.Println(err)
		}
		user.Date = time.Now()
		ret = append(ret, user)
		if limit > 0 {
			if id > limit {
				break
			}
		}
	}
	return ret
}

func (s server) getUsersCounterList(limit int) []UserCounter {
	var listUsers []UserCounter
	tags, err := s.redis.SMembers("tag_list").Result()
	if err != nil {
		log.Println(err)
	}
	for _, tag := range tags {
		var uc UserCounter
		uc.Tag = tag
		uc.Users = s.getUsersCounter(tag, limit)
		listUsers = append(listUsers, uc)
	}
	return listUsers
}

func (s server) getMeanTimes(tag string, key string, aggrega int64) (float64, int64) {
	nb, erra := s.redis.LLen(key + tag + "_" + string(aggrega)).Result()
	if erra != nil {
		log.Println(erra)
	}
	val, errb := s.redis.LRange(key+tag+"_"+string(aggrega), 0, nb).Result()
	if errb != nil {
		log.Println(errb)
	}
	s.redis.Del(key + tag + "_" + string(aggrega))
	var sum int64
	var i int64
	sum = 0
	for i = 0; i != nb; i++ {
		tmp, _ := strconv.ParseInt(val[i], 10, 64)
		sum = sum + tmp
	}
	var ret float64
	ret = (float64(sum) / float64(nb))
	return ret, nb
}

func (s server) computeAggregas() Aggrega {
	var agg Aggrega

	agg.Num = s.getNbAggregation()
	s.addAggregation()

	tags, err := s.redis.SMembers("tag_list").Result()
	if err != nil {
		log.Println(err)
	}
	for _, tag := range tags {
		agg.InjectorMean, agg.InjectorNb = s.getMeanTimes(tag, "injectTimes_", agg.Num)
		agg.NormalizerMean, agg.NormalizerNb = s.getMeanTimes(tag, "normTimes_", agg.Num)
	}

	agg.CreateTime = time.Now()

	return agg
}

func (s server) addAggregation() {
	s.redis.Incr("aggregas")
}

func (s server) getNbAggregation() int64 {
	nb, err := s.redis.Get("aggregas").Int64()
	if err != nil {
		log.Println(err)
	}
	return nb
}
