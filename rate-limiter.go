package ratelimiter

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

const  (
	SLIDINGWINDOW int=iota+1
	FIXEDWINDOW
)

type RateLimiter interface {
	Allow(tenantID string, t time.Time) bool
}

func NewRateLimiter(windowSize time.Duration, limit Limit, store *redis.Client, limitingAlgo int) RateLimiter{
	switch limitingAlgo {
	case SLIDINGWINDOW:
		log.Println("activating sliding window rate limiter")
		return NewSlidingWindow(windowSize,limit,store)
	case FIXEDWINDOW:
		log.Println("activating fixed window rate limiter")
		return NewFixWindow(windowSize,limit,store)
	default:
		log.Println("activating default sliding window rate limiter")
		return NewSlidingWindow(windowSize,limit,store)
	}
}