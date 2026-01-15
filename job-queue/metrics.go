package main

import "sync/atomic"

type Counter struct{
	val int64
}

func (c *Counter)Inc()  {
	atomic.AddInt64(&c.val,1)
}

func (c *Counter)Get() int64 {
	return atomic.LoadInt64(&c.val)

}

func (c *Counter)Dec()  {
	atomic.AddInt64(&c.val, -1)
}

var(
	metricsJobSuccess = &Counter{}
	metricsJobFailure = &Counter{}
	metricsJobInProgress = &Counter{}
	metricsJobsDead = &Counter{}
)