package main

import (
	"golang.org/x/time/rate"
)

// TODO: not thread safe
// TODO: limiters memory leak
type RequestsLimiter struct {
	limiters map[string]*rate.Limiter
}

func (rl *RequestsLimiter) Allow(key string) bool {
	limiter := reqsLimiter.limiters[key]
	if limiter == nil {
		limiter = rate.NewLimiter(1, 2)
		rl.limiters[key] = limiter
	}
	return limiter.Allow()
}
