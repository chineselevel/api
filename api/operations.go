package main

import "github.com/garyburd/redigo"

type Operations struct {
	redis redigo.Conn
}

func NewOperations() (o *Operations) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
	}
	defer c.Close()

	o = Operations{
		redis: c,
	}
	return o
}
