package api

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type Operations struct {
	Redis redis.Conn
}

func NewOperations() (o *Operations) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}

	o = &Operations{
		Redis: c,
	}
	return o
}

// GetRank returns the rank of a key from the Redis key/value store.
func (o *Operations) GetRank(s string) (rank int) {
	log.Print(s)
	rank, err := redis.Int(o.Redis.Do("ZREVRANK", "words", s))
	if err != nil {
		log.Print("Error: ", err)
	}
	return rank
}
