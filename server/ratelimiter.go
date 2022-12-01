package server

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// AddIP adds the IP adress to the connection map in IPRateLimiter
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

// GetLimiter returns a pointer to the rate.Limiter in the IPRateLimiter connection map. If it doesn't exist, it adds it.
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]
	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}
	i.mu.Unlock()
	return limiter
}

// LimitMiddleware is a middleware handler function that handles connections using a token bucket rate-limiter algorithm,
// if there are too many connections from the same IP address, the client gets a StatusTooManyRequests error.
func (i *IPRateLimiter) LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := i.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CleanUpVisitorMap runs as a go routine and every 30 minutes reassigns i.ips to an empty map to avoid memory leak
func (i *IPRateLimiter) CleanUpVisitorMap() {
	for {
		time.Sleep(5 * time.Minute)

		if time.Since(i.lastCleaned) > 30*time.Minute {
			i.ips = make(map[string]*rate.Limiter)
			i.lastCleaned = time.Now()
		}
	}
}
