package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	config "github.com/leetcode-golang-classroom/golang-gocraft-with-background-jobs/internal/config"
)

// Redis Pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", config.C.REDIS_URI)
	},
}

type User struct {
	ID    int64
	Email string
	Name  string
}
type Context struct {
	currentUser *User
}

// Middleware to fetch the user Object from userID
func (c *Context) FindCurrentUser(job *work.Job, next work.NextMiddlewareFunc) error {
	// If there's a user_id param
	if _, ok := job.Args["userID"]; ok {
		userID := job.ArgInt64("userID")
		// Simuate query from db
		c.currentUser = &User{ID: userID, Email: "test" + strconv.Itoa(int(userID)) + "@gmail.com", Name: "Test User"}
		if err := job.ArgError(); err != nil {
			return err
		}
	}
	return next()
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Printf("Starting a new job: %v with ID %v", job.Name, job.ID)
	return next()
}

// Create job enqueuer

var enqueuer = work.NewEnqueuer(config.C.REDIS_NS, redisPool)

func main() {
	pool := work.NewWorkerPool(Context{}, 10, config.C.REDIS_NS, redisPool)
	// Middlewares
	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCurrentUser)
	// Name to job map
	pool.JobWithOptions("email",
		work.JobOptions{Priority: 10, MaxFails: 1},
		(*Context).SendEmail,
	)
	pool.JobWithOptions("report",
		work.JobOptions{Priority: 10, MaxFails: 1},
		(*Context).Report,
	)
	pool.Start()

	// Wait for a signal to quit
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (c *Context) SendEmail(job *work.Job) error {
	addr := c.currentUser.Email
	subject := job.ArgString("subject")

	if err := job.ArgError(); err != nil {
		return err
	}
	fmt.Printf("Sending mail to %v with subject %v\n", addr, subject)
	time.Sleep(time.Second * 2)
	return nil
}

func (c *Context) Report(job *work.Job) error {
	fmt.Println("Preparing report...")
	// for i := range 360 {
	// 	time.Sleep(time.Second * 10)
	// 	job.Checkin("i=" + fmt.Sprint(i))
	// }
	time.Sleep(time.Second * 10)
	// Send the report via email
	enqueuer.Enqueue("email", work.Q{"userID": c.currentUser.ID, "subject": "Report is Ready"})
	return nil
}
