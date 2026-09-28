package security

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"
	"sync"
	"time"

	contractshttp "github.com/goravel/framework/contracts/http"
)

const (
	loginWindow = 30 * time.Minute
	bucketTTL   = 2 * time.Hour
)

type attemptBucket struct {
	Failures    int
	LastFailure time.Time
	LockedUntil time.Time
}

type AuthLimiter struct {
	mu      sync.Mutex
	buckets map[string]*attemptBucket
	now     func() time.Time
}

func NewAuthLimiter() *AuthLimiter {
	return &AuthLimiter{
		buckets: make(map[string]*attemptBucket),
		now:     time.Now,
	}
}

var defaultAuthLimiter = NewAuthLimiter()

func LoginAttempts() *AuthLimiter {
	return defaultAuthLimiter
}

func ClientAddress(ctx contractshttp.Context) string {
	if ctx != nil && ctx.Request() != nil {
		if request := ctx.Request().Origin(); request != nil {
			remote := strings.TrimSpace(request.RemoteAddr)
			if host, _, err := net.SplitHostPort(remote); err == nil && host != "" {
				return normalizeIP(host)
			}
			if remote != "" {
				return normalizeIP(remote)
			}
		}

		if ip := strings.TrimSpace(ctx.Request().Ip()); ip != "" {
			return normalizeIP(ip)
		}
	}

	return "unknown"
}

func IdentityFingerprint(identity string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(identity))))
	return hex.EncodeToString(sum[:])[:16]
}

func (l *AuthLimiter) Check(ip, identity string) (time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.cleanup(now)

	var retry time.Duration
	for _, key := range l.keys(ip, identity) {
		bucket := l.bucket(key, now)
		if bucket.LockedUntil.After(now) {
			if wait := bucket.LockedUntil.Sub(now); wait > retry {
				retry = wait
			}
		}
	}
	return retry, retry > 0
}

func (l *AuthLimiter) RecordFailure(ip, identity string) time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.cleanup(now)

	var retry time.Duration
	for _, keyed := range l.typedKeys(ip, identity) {
		bucket := l.bucket(keyed.key, now)
		bucket.Failures++
		bucket.LastFailure = now

		if duration := lockDuration(keyed.kind, bucket.Failures); duration > 0 {
			until := now.Add(duration)
			if until.After(bucket.LockedUntil) {
				bucket.LockedUntil = until
			}
		}
		if bucket.LockedUntil.After(now) {
			if wait := bucket.LockedUntil.Sub(now); wait > retry {
				retry = wait
			}
		}
	}

	return retry
}

func (l *AuthLimiter) RecordSuccess(ip, identity string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.buckets, pairKey(ip, identity))
	delete(l.buckets, identityKey(identity))

	ipKey := ipKey(ip)
	if bucket, ok := l.buckets[ipKey]; ok {
		if bucket.Failures <= 2 {
			delete(l.buckets, ipKey)
		} else {
			bucket.Failures -= 2
			bucket.LockedUntil = time.Time{}
		}
	}
}

func (l *AuthLimiter) ClearIdentity(identity string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fingerprint := IdentityFingerprint(identity)
	delete(l.buckets, "identity:"+fingerprint)
	suffix := ":" + fingerprint
	for key := range l.buckets {
		if strings.HasPrefix(key, "pair:") && strings.HasSuffix(key, suffix) {
			delete(l.buckets, key)
		}
	}
}

func (l *AuthLimiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]*attemptBucket)
}

type typedKey struct {
	kind string
	key  string
}

func (l *AuthLimiter) typedKeys(ip, identity string) []typedKey {
	return []typedKey{
		{kind: "pair", key: pairKey(ip, identity)},
		{kind: "ip", key: ipKey(ip)},
		{kind: "identity", key: identityKey(identity)},
	}
}

func (l *AuthLimiter) keys(ip, identity string) []string {
	typed := l.typedKeys(ip, identity)
	return []string{typed[0].key, typed[1].key, typed[2].key}
}

func (l *AuthLimiter) bucket(key string, now time.Time) *attemptBucket {
	bucket, ok := l.buckets[key]
	if !ok {
		bucket = &attemptBucket{}
		l.buckets[key] = bucket
		return bucket
	}

	if !bucket.LastFailure.IsZero() && now.Sub(bucket.LastFailure) > loginWindow && !bucket.LockedUntil.After(now) {
		bucket.Failures = 0
		bucket.LockedUntil = time.Time{}
	}
	return bucket
}

func (l *AuthLimiter) cleanup(now time.Time) {
	for key, bucket := range l.buckets {
		if !bucket.LastFailure.IsZero() && now.Sub(bucket.LastFailure) > bucketTTL && !bucket.LockedUntil.After(now) {
			delete(l.buckets, key)
		}
	}
}

func pairKey(ip, identity string) string {
	return "pair:" + strings.TrimSpace(ip) + ":" + IdentityFingerprint(identity)
}

func ipKey(ip string) string {
	return "ip:" + strings.TrimSpace(ip)
}

func identityKey(identity string) string {
	return "identity:" + IdentityFingerprint(identity)
}

func normalizeIP(value string) string {
	value = strings.Trim(value, "[] ")
	if parsed := net.ParseIP(value); parsed != nil {
		return parsed.String()
	}
	return value
}

func lockDuration(kind string, failures int) time.Duration {
	switch kind {
	case "pair":
		switch {
		case failures >= 15:
			return 30 * time.Minute
		case failures >= 10:
			return 5 * time.Minute
		case failures >= 7:
			return time.Minute
		case failures >= 5:
			return 15 * time.Second
		}
	case "ip":
		switch {
		case failures >= 60:
			return time.Hour
		case failures >= 35:
			return 10 * time.Minute
		case failures >= 20:
			return time.Minute
		}
	case "identity":
		switch {
		case failures >= 40:
			return 10 * time.Minute
		case failures >= 20:
			return 30 * time.Second
		}
	}
	return 0
}
