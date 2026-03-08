// Package middleware provides HTTP middleware for the todo.open API.
package middleware

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

const idempotencyKeyHeader = "X-Idempotency-Key"
const idempotencyTTL = 5 * time.Minute

type cachedResponse struct {
	status   int
	headers  http.Header
	body     []byte
	storedAt time.Time
}

// IdempotencyStore holds cached responses keyed by idempotency key.
type IdempotencyStore struct {
	mu    sync.RWMutex
	cache map[string]cachedResponse
	nowFn func() time.Time
}

// NewIdempotencyStore creates a new in-memory idempotency store.
func NewIdempotencyStore(nowFn func() time.Time) *IdempotencyStore {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &IdempotencyStore{
		cache: make(map[string]cachedResponse),
		nowFn: nowFn,
	}
}

// Middleware returns an HTTP middleware that deduplicates requests sharing the same X-Idempotency-Key.
// Cached responses are returned for 5 minutes after the first successful request.
func (s *IdempotencyStore) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(idempotencyKeyHeader)
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		now := s.nowFn()

		// Check cache.
		s.mu.RLock()
		cached, ok := s.cache[key]
		s.mu.RUnlock()

		if ok && now.Sub(cached.storedAt) < idempotencyTTL {
			// Return cached response.
			for k, vals := range cached.headers {
				for _, v := range vals {
					w.Header().Add(k, v)
				}
			}
			w.Header().Set("X-Idempotency-Replayed", "true")
			w.WriteHeader(cached.status)
			_, _ = w.Write(cached.body)
			return
		}

		// Capture the response.
		rw := &responseRecorder{header: make(http.Header), code: http.StatusOK}
		next.ServeHTTP(rw, r)

		// Cache it.
		s.mu.Lock()
		s.cache[key] = cachedResponse{
			status:   rw.code,
			headers:  rw.header.Clone(),
			body:     rw.body.Bytes(),
			storedAt: now,
		}
		// Evict stale entries opportunistically.
		for k, v := range s.cache {
			if now.Sub(v.storedAt) >= idempotencyTTL {
				delete(s.cache, k)
			}
		}
		s.mu.Unlock()

		// Write actual response.
		for k, vals := range rw.header {
			for _, v := range vals {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(rw.code)
		_, _ = w.Write(rw.body.Bytes())
	})
}

// responseRecorder captures a response for caching.
type responseRecorder struct {
	header http.Header
	code   int
	body   bytes.Buffer
}

func (r *responseRecorder) Header() http.Header         { return r.header }
func (r *responseRecorder) WriteHeader(code int)        { r.code = code }
func (r *responseRecorder) Write(b []byte) (int, error) { return r.body.Write(b) }
