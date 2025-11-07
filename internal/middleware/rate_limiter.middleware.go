package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
	"user-management-api/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Client struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	mutex   = &sync.Mutex{}
	clients = make(map[string]*Client)
)

func getClientIP(c *gin.Context) string {
	// implement logic to get client IP
	ip := c.ClientIP()
	if ip == "" {
		ip = c.Request.RemoteAddr // lay ip remote if client ip is empty
	}
	return ip
}

func getRateLimiter(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()
	client, exist := clients[ip]
	if !exist {
		// get env variable for rate limiter
		reqSec := utils.GetIntEnv("RATE_LIMITER_REQUEST_SEC", 5)
		brustSec := utils.GetIntEnv("RATE_LIMITER_REQUEST_BRUST", 10)

		limiter := rate.NewLimiter(rate.Limit(reqSec), brustSec) // 5 requests/s and 10 burst
		newClient := &Client{limiter, time.Now()}
		clients[ip] = newClient
		//log.Printf("A client with IP %s - {limiter: %+v , lastseen: %+s}", ip, newClient.Limiter, newClient.LastSeen)
		return limiter
	}

	//log.Printf("A client with IP %s - {limiter: %+v , lastseen: %+s}", ip, client.Limiter, client.LastSeen)

	client.LastSeen = time.Now()
	return client.Limiter
}
func CleanUpClients() {
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		for ip, client := range clients {
			if time.Since(client.LastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mutex.Unlock()
	}
}

// / ab -n 20 -c 1 http://localhost:8080
func RateLimiterMiddleware(ratelimitlogger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// implement rate limit logic here

		// get client IP
		ip := getClientIP(c)
		limiter := getRateLimiter(ip)

		if !limiter.Allow() {
			// check if rate limit exceeded to avoid log spam
			if shutLogRateLimiter(ip) {
				// write log if rate limit exceeded
				ratelimitlogger.Warn().
					Str("client_ip", c.ClientIP()).
					Str("protocol", c.Request.Proto).
					Str("user_agent", c.Request.UserAgent()).
					Str("referer", c.Request.Referer()).
					Str("host", c.Request.Host).
					Str("remote_address", c.Request.RemoteAddr).
					Str("request_uri", c.Request.RequestURI).
					Str("method", c.Request.Method).
					Interface("headers", c.Request.Header).
					Str("path", c.Request.URL.Path).
					Str("quey", c.Request.URL.RawQuery).
					Msg("Rate limit exceeded")

			}

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
				"code":  "TOO_MANY_REQUESTS",
			})
			return
		}
		//log.Printf("Rate limit for IP %s exceeded", ip)
		c.Next()
	}
}

var rateLimitLogCache = sync.Map{} //  cache for rate limit log with ip in a time\

var rateLimitLogTTL = 10 * time.Second //default value

func shutLogRateLimiter(ip string) bool {
	ttl_val, e := strconv.Atoi(utils.GetEnv("RATE_LIMITER_LOG_TTL", "10"))
	if e != nil {
		rateLimitLogTTL = 10 * time.Second
	} else {
		rateLimitLogTTL = time.Duration(ttl_val) * time.Second
	}
	now := time.Now()
	value, ok := rateLimitLogCache.Load(ip)
	if ok {
		if t, ok := value.(time.Time); ok && now.Sub(t) < rateLimitLogTTL {
			return true
		}
	}
	rateLimitLogCache.Store(ip, now)
	return true
}
