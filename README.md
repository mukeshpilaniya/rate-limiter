
### Usages of Sliding window rate limiter
```azure
package main

import (
	"github.com/go-redis/redis"
	ratelimiter "github.com/mukeshpilaniya/rate-limiter"
	"time"
)

func main() {
	// Create a Redis Client for storing a keys
	redisClient :=redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
	})

	// Check whether the Redis connection is in working state or not
	strCmd :=redisClient.Ping()
	if strCmd.Err() !=nil{
		panic(strCmd.Err())
	}

	// Create a rate limiter of sliding window type
	rtl :=ratelimiter.NewRateLimiter(time.Minute,25,redisClient,ratelimiter.SLIDINGWINDOW)

	// Create a rate limiter of Fixed window type
	//rtl :=ratelimiter.NewRateLimiter(time.Minute,25,redisClient,FIXEDWINDOW)

	// test rate limiter
	for {
		rtl.Allow("pilaniya1", time.Now())
		time.Sleep(time.Second*1)
	}
}
```

### Test RateLimiter by sending concurrent request
```azure
package main

import (
	"github.com/go-redis/redis"
	ratelimiter "github.com/mukeshpilaniya/rate-limiter"
	"time"
)

func main() {
	// Create a Redis Client for storing a keys
	redisClient :=redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
	})

	// Check whether the Redis connection is in working state or not
	strCmd :=redisClient.Ping()
	if strCmd.Err() !=nil{
		panic(strCmd.Err())
	}

	// Create a rate limiter of sliding window type
	rtl :=ratelimiter.NewRateLimiter(time.Minute,25,redisClient,ratelimiter.SLIDINGWINDOW)

	// Create a rate limiter of Fixed window type
	//rtl :=NewRateLimiter(time.Minute,25,redisClient,FIXEDWINDOW)

	// test rate limiter
	go func(){
		for {
			go  rtl.Allow("pilaniya1", time.Now())
			time.Sleep(time.Second*1)
		}
	}()
	go func(){
		for {
			go  rtl.Allow("pilaniya1", time.Now())
			time.Sleep(time.Second*1)
		}
	}()

	go func(){
		for {
			go  rtl.Allow("pilaniya1", time.Now())
			time.Sleep(time.Second*1)
		}
	}()
	go func(){
		for {
			go  rtl.Allow("pilaniya1", time.Now())
			time.Sleep(time.Second*1)
		}
	}()
	go func(){
		for {
			go  rtl.Allow("pilaniya1", time.Now())
			time.Sleep(time.Second*1)
		}
	}()
	for{

	}
}
```