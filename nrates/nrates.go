package nrates

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/influx6/npkg/nerror"
	"github.com/influx6/npkg/ntrace"
	openTracing "github.com/opentracing/opentracing-go"
)

// Rate is the rate of allowed requests. We support
// r/min and r/second.
type Rate int

const (
	// PerSecond allows us to accept x requests per second
	PerSecond Rate = iota
	// PerMinute allows us to accept x requests per minute
	PerMinute
)

// HHMMSS formats a timestamp as HH:MM:SS
// Reference: https://yourbasic.org/golang/format-parse-string-time-date-example/
// const HHMMSS = "15:04:05"

// HHMM formats a timestamp as HH:MM
// Reference: https://yourbasic.org/golang/format-parse-string-time-date-example/
// const HHMM = "15:04"

type RedisIncr struct {
	Client *redis.Client
}

func NewRedisIncr(config *redis.Options) (*RedisIncr, error) {
	client := redis.NewClient(config)

	return &RedisIncr{
		Client: client,
	}, nil
}

func (b *RedisIncr) Reset(ctx context.Context, r Request) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewMethodSpanFromContext(ctx); span != nil {
		defer span.Finish()
	}

	status := b.Client.Ping(ctx)
	if err := status.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	_, err := b.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		return pipe.Set(ctx, r.Owner(), 0, 0).Err()
	})
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

func (b *RedisIncr) Count(ctx context.Context, r Request) (int64, error) {
	var span openTracing.Span
	if ctx, span = ntrace.NewMethodSpanFromContext(ctx); span != nil {
		defer span.Finish()
	}

	status := b.Client.Ping(ctx)
	if err := status.Err(); err != nil {
		return -1, nerror.WrapOnly(err)
	}

	var res = b.Client.Get(ctx, r.Owner())
	if err := res.Err(); err != nil {
		return -1, nerror.WrapOnly(err)
	}

	var count, readErr = strconv.Atoi(res.Val())
	if readErr != nil {
		return -1, nerror.WrapOnly(readErr)
	}

	return int64(count), nil
}

func (b *RedisIncr) Inc(ctx context.Context, r Request, dur time.Duration) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewMethodSpanFromContext(ctx); span != nil {
		defer span.Finish()
	}

	status := b.Client.Ping(ctx)
	if err := status.Err(); err != nil {
		return nerror.WrapOnly(err)
	}

	_, err := b.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, r.Owner())
		pipe.Expire(ctx, r.Owner(), dur)
		return nil
	})
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}

type Request interface {
	Owner() string
	Data() interface{}
}

type IncrementStore interface {
	Reset(ctx context.Context, r Request) error
	Count(ctx context.Context, r Request) (int64, error)
	Inc(ctx context.Context, r Request, dur time.Duration) (int64, error)
}

// NewRateLimiter returns a new Limiter.
func NewFactory(db IncrementStore, rate Rate) *LimiterFactory {
	return &LimiterFactory{Store: db, Rate: rate}
}

type LimiterFactory struct {
	Store IncrementStore
	Rate  Rate
}

// NewLimiter creates a new Limiter.
func (f LimiterFactory) New(max int64) *RateLimiter {
	return &RateLimiter{
		store: f.Store,
		rate:  f.Rate,
		max:   max,
	}
}

// NewLimiter creates a new Limiter.
func (f LimiterFactory) NewLimiter(rate Rate, max int64) *RateLimiter {
	return &RateLimiter{
		store: f.Store,
		rate:  rate,
		max:   max,
	}
}

type RateLimiter struct {
	store IncrementStore
	rate  Rate
	max   int64
}

// RateLimit applies basic rate limiting to an HTTP request as described
// in Redis' onboarding documentation.
// Reference: https://redislabs.com/redis-best-practices/basic-rate-limiting/
func (l *RateLimiter) RateLimit(ctx context.Context, r Request) error {
	var span openTracing.Span
	if ctx, span = ntrace.NewMethodSpanFromContext(ctx); span != nil {
		defer span.Finish()
	}

	var expiry time.Duration

	if l.rate == PerSecond {
		expiry = time.Second
	} else {
		expiry = time.Minute
	}

	var count, err = l.store.Count(ctx, r)
	if err != nil {
		return nerror.WrapOnly(err)
	}

	if count >= l.max {
		return nerror.New("requests are throttled, try again later")
	}

	_, err = l.store.Inc(ctx, r, expiry)
	if err != nil {
		return nerror.WrapOnly(err)
	}
	return nil
}
