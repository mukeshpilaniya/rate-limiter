package ratelimiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

type SlidingWindow struct {
	mu sync.Mutex
	windowSize time.Duration
	limit Limit
	db *redis.Client
}

func NewSlidingWindow(windowSize time.Duration,limit Limit, store *redis.Client) *SlidingWindow{
	return &SlidingWindow{
		windowSize: windowSize,
		limit: limit,
		db: store,
	}
}

func (sw *SlidingWindow) getCount(tenantID string, t time.Time)  (int64, error){
	key :=fmt.Sprintf(tenantID+":"+strconv.FormatInt(t.Unix(),10))
	strCmd:= sw.db.Get(key)
	v, _ := strCmd.Int64()
	if v==0{
		return int64(0),nil
	}
	return v, nil
}

func (sw *SlidingWindow) Allow(tenantID string, t time.Time) bool{
	var executePipe bool
	currWinT :=t.Truncate(sw.windowSize)
	prevWinT :=t.Add(-sw.windowSize)
	prevWinT =prevWinT.Truncate(sw.windowSize)
	key :=fmt.Sprintf(tenantID+":"+strconv.FormatInt(currWinT.Unix(),10))

	sw.mu.Lock()
	prevCnt, _ :=sw.getCount(tenantID,prevWinT)
	currCnt, _ :=sw.getCount(tenantID,currWinT)
 	weightedCnt :=prevCnt*((int64(sw.windowSize)-(t.Unix()-currWinT.Unix()))/int64(sw.windowSize))+currCnt
	if weightedCnt+1> int64(sw.limit){
		log.Println("rate limit is exceed for the tenant", tenantID)
		sw.mu.Unlock()
		return false
	}
	pipe := sw.db.Pipeline()
	defer func() {
		if executePipe{
			_, err:= pipe.Exec()
			if err!=nil{
				log.Println("error executing redis commands",err)
			}
		}
	}()
	if currCnt == int64(0){
		pipe.Set(key,1,sw.windowSize*2)
		log.Printf("allowing request for the tenant %s at time %d with counter %d \n",tenantID,currWinT.Unix(), 1)
	}else{
			pipe.Incr(key)
			log.Printf("allowing request for the tenant %s at time %d with counter %d \n",tenantID,currWinT.Unix(), currCnt+1)
	}
	executePipe=true
	sw.mu.Unlock()
	return true
}