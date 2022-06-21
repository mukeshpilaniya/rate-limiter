package ratelimiter

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"sync"
	"time"
)

type Limit int64

type FixedWindow struct {
	mu sync.Mutex
	windowSize time.Duration
	limit Limit
	db *redis.Client
}

func NewFixWindow(windowSize time.Duration,limit Limit, store *redis.Client) *FixedWindow{
	return &FixedWindow{
		windowSize: windowSize,
		limit: limit,
		db: store,
	}
}

func (fw *FixedWindow) getCount(tenantID string, t time.Time)  (int64, error){
	key :=fmt.Sprintf(tenantID+":"+strconv.FormatInt(t.Unix(),10))
	strCmd:= fw.db.Get(key)
	v, _ := strCmd.Int64()
	if v==0{
		return int64(0),nil
	}
	return v, nil
}

func (fw *FixedWindow) Allow(tenantID string, t time.Time) bool{
	var executePipe bool
	currWinT :=t.Truncate(fw.windowSize)
	key :=fmt.Sprintf(tenantID+":"+strconv.FormatInt(currWinT.Unix(),10))
	pipe := fw.db.Pipeline()
	fw.mu.Lock()
	defer func() {
		if executePipe{
			_, err:= pipe.Exec()
			if err!=nil{
				log.Println("error executing redis commands",err)
			}
		}
	}()
	cnt,_ := fw.getCount(tenantID,currWinT)
	if cnt == 0{
		pipe.Set(key,1,fw.windowSize*2)
		log.Printf("allowing request for the tenant %s at time %d with counter %d \n",tenantID,currWinT.Unix(), 1)
	}else{
		if cnt>=int64(fw.limit){
			log.Printf("rate limit is exceed for the tenant %s\n", tenantID)
			fw.mu.Unlock()
			return false
		}else{
			pipe.Incr(key)
			fmt.Printf("allowing request for the tenant %s at time %d with counter %d \n",tenantID,currWinT.Unix(), cnt+1)
		}
	}
	executePipe=true
	fw.mu.Unlock()
	return true
}