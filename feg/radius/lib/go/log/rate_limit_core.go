package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

// rateLimitCore A zap logger which has rate-limiter to limit the rate of logs
// all logs that are written through this logger are rate-limited without regard to their content
// Note the difference between this & the zapcore.Sampler() which takes a sample of a log stream.
// zapcore.Sampler is more appropriate when your log stream is high-rate by nature (e.g.: you have 100M request per minute so a single log per request is already a huge amount of logs).
// in the above case, sampling is a reasonable approach, so it will sample 1/M of logs. and since the rate is changing moderately over time, M can be adjusted every few month.
// This core is meant at logs which are expected to be at very low-rate (e.g.: Warn() & Error() calls) yet they are caused by an external event (e.g.: parameter received in Rest API), hence in case of a bug of DDOS attack, the rate can be very high
// here we would like to ensure that the external entity *cannot* flood the system with logs, thereby degrading its service.
type rateLimitCore struct {
	zapcore.Core
	// rate limiter is shared among all instances created from the initial logger when adding fields.
	limiter *rate.Limiter
	ctx     context.Context
}

// NewRateLimitedLogger returns a new rate-limited core
func NewRateLimitedLogger(core zapcore.Core, ctx context.Context, limiter *rate.Limiter) zapcore.Core {
	return &rateLimitCore{
		Core:    core,
		limiter: limiter,
		ctx:     ctx,
	}
}

// With returns a new instance of the rate limited core
// the new core shares the same rate-limiter so their aggregate logging is limited
func (rll *rateLimitCore) With(fields []zap.Field) zapcore.Core {
	return &rateLimitCore{
		Core:    rll.Core.With(fields),
		limiter: rll.limiter,
		ctx:     rll.ctx,
	}
}

// Check checks whether a log is allowed by rate limiter or not
func (rll *rateLimitCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !rll.Enabled(ent.Level) {
		return ce
	}
	if !rll.limiter.Allow() {
		return ce
	}
	return rll.Core.Check(ent, ce)
}
