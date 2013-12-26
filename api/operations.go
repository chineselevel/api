package api

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type Operations struct {
	RedisPool *redis.Pool
}

func NewOperations(server string, password string) (o *Operations) {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	o = &Operations{
		RedisPool: pool,
	}
	return o
}

// GetRank returns the rank of a key from the Redis key/value store.
func (o *Operations) GetRank(s string) (rank int) {
	conn := o.RedisPool.Get()
	defer conn.Close()

	log.Println("GETRANK")
	rank, err := redis.Int(conn.Do("ZREVRANK", "words", s))
	if err != nil {
		log.Print("Error: ", err)
		rank = -1
	}
	log.Println(s, rank)
	return rank
}
