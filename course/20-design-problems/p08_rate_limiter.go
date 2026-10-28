package design

import (
	"sync"
	"time"
)

// Problem: Rate Limiter (token bucket)
// Difficulty: Medium
//
// Build a limiter that lets bursts of up to `capacity` requests through,
// then sustains `rate` requests per second. Allow() answers immediately
// (no blocking) whether one request may proceed; AllowN(n) asks for n at
// once, all or nothing. It must be safe to call from many goroutines.
//
// Example (capacity 3, rate 1/s, starting full):
//
//	t=0.0s  Allow ×4  -> true true true false   (burst used up)
//	t=1.0s  Allow     -> true                   (one token refilled)
//	t=10s   AllowN(5) -> false                  (bucket caps at 3)
//
// Reference: https://en.wikipedia.org/wiki/Token_bucket
//
// Pattern: lazy refill. Don't run a ticker goroutine; on each call compute
// how much time passed since the last call and add elapsed·rate tokens,
// capped at capacity. Then spend tokens if enough are available.
//
//	tokens = min(capacity, tokens + (now - last)·rate)
//	if tokens >= n { tokens -= n; allow }
//
// The clock is injected so tests are deterministic (no time.Sleep).
//
// Time: O(1) per call. Space: O(1).
//
// Safe for concurrent use: the read-refill-spend sequence must be atomic,
// so it runs under one sync.Mutex (two separate atomic ops would race).
//
// See the production-style version with per-key buckets and middleware:
// https://github.com/neerajsinghi/system-design-patterns-go (patterns/reliability/rate-limiter-token-bucket)

// TokenBucket is a concurrency-safe token-bucket rate limiter.
type TokenBucket struct {
	mu       sync.Mutex
	capacity float64
	rate     float64 // tokens added per second
	tokens   float64
	last     time.Time
	now      func() time.Time
}

// NewTokenBucket returns a full bucket. now may be nil to use time.Now.
func NewTokenBucket(capacity int, ratePerSecond float64, now func() time.Time) *TokenBucket {
	if now == nil {
		now = time.Now
	}
	return &TokenBucket{
		capacity: float64(capacity),
		rate:     ratePerSecond,
		tokens:   float64(capacity),
		last:     now(),
		now:      now,
	}
}

// refill adds tokens for the time elapsed since the last call. Caller holds mu.
func (b *TokenBucket) refill() {
	now := b.now()
	if elapsed := now.Sub(b.last).Seconds(); elapsed > 0 { // ignore clocks going backwards
		b.tokens = min(b.capacity, b.tokens+elapsed*b.rate)
		b.last = now
	}
}

// Allow reports whether one request may proceed now.
func (b *TokenBucket) Allow() bool { return b.AllowN(1) }

// AllowN reports whether n requests may proceed now, consuming n tokens if so.
func (b *TokenBucket) AllowN(n int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= float64(n) {
		b.tokens -= float64(n)
		return true
	}
	return false
}

// Tokens returns the current token count (after refilling), for observability.
func (b *TokenBucket) Tokens() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	return b.tokens
}
