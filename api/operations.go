package api

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/hermanschaaf/go-mafan"
)

type Operations struct {
	RedisPool *redis.Pool
}

func NewOperations(server string, password string) (o *Operations) {
	pool := &redis.Pool{
		MaxIdle:     3,
		MaxActive:   10,
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

	rank, err := redis.Int(conn.Do("ZSCORE", "ranks", s))
	if err != nil {
		log.Print("Error: ", err)
		rank = -1
	}
	log.Println(s, rank)
	return rank
}

// GetRank returns the ranks of words from the Redis key/value store.
func (o *Operations) GetRanks(words []string) (ranks []int) {
	conn := o.RedisPool.Get()
	defer conn.Close()

	unknown := 0

	for i := 0; i < len(words); i++ {
		if mafan.IsHanzi(words[i]) {
			rank, err := redis.Int(conn.Do("ZSCORE", "ranks", words[i]))
			if err != nil {
				log.Print("Error: ", err)
				rank = -1
			}
			if rank >= 0 {
				ranks = append(ranks, rank)
			} else {
				unknown += 1
			}
		}
	}

	// add unknown words to end of ranks list
	// as equal to biggest known word
	biggest := 0
	if len(ranks) > 0 {
		biggest = ranks[len(ranks)-1]
	}
	for i := 0; i < unknown; i++ {
		ranks = append(ranks, biggest)
	}
	return ranks
}
