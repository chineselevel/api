package api

import "github.com/garyburd/redigo/redis"

type Operations struct {
	Redis redis.Conn
}

func NewOperations() (o *Operations) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}
	defer c.Close()

	o = &Operations{
		Redis: c,
	}
	return o
}
