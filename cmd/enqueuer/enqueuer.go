package main

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	config "github.com/leetcode-golang-classroom/golang-gocraft-with-background-jobs/internal/config"
)

// Redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", config.C.REDIS_URI)
	},
}

// Create job enqueuer
var enqueuer = work.NewEnqueuer(config.C.REDIS_NS, redisPool)

func main() {
	_, err := enqueuer.Enqueue("email",
		work.Q{"userID": 10, "subject": "Just testing"},
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = enqueuer.Enqueue("report",
		work.Q{"userID": 5},
	)
	if err != nil {
		log.Fatal(err)
	}
}
