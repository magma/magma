package log

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

// TestRateLimitLogger test the rate-limit logger
// h=whenever the 2 time.Sleep() call is made, the delay is long enough to make the next log call go through.
// so we randomly pick a value (0 or 1) & decide whether to sleep or no, thereby, whether we log with fmt.Printf() or logger.Debug()
// in any case, 20 lines are spitted out, 10 of each logger
// i use 2 loggers to show that the fact we create a new logger from existing one, still shares the same rate-limiter.
// so if multiple consecutive calls to time.Sleep() are skipped, multiple logger calls will be dropped, though they are executed on 2 different instances.
func TestRateLimitLogger(t *testing.T) {
	ratePerSec := 5
	delayAllowsNextLog := int(1000 / ratePerSec)
	require.True(t, ratePerSec <= 1000) // bcz we sleep in [msec]
	ctx := context.Background()
	rateLimiter := rate.NewLimiter(rate.Limit(ratePerSec), 1)
	rateLimitCore := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return NewRateLimitedLogger(core, ctx, rateLimiter)
	})
	logger, _ := zap.NewDevelopment()
	rateLimitLogger := logger.WithOptions(rateLimitCore)
	rateLimitLoggerWithCtx := rateLimitLogger.With(zap.String("logger_instance", "with field"))
	for i := 0; i < 2*ratePerSec; i++ {
		if rand.Intn(2) == 0 {
			time.Sleep(time.Duration(delayAllowsNextLog) * time.Millisecond)
		} else {
			fmt.Printf("next log will be dropped by rate-limited logger (count=%d)\n", i)
		}

		rateLimitLogger.Debug("write with limiter", zap.Int("count", i))
		if rand.Intn(2) == 0 {
			time.Sleep(time.Duration(delayAllowsNextLog) * time.Millisecond)
		} else {
			fmt.Printf("next log will be dropped by rate-limited logger (count=%d)\n", i)
		}
		rateLimitLoggerWithCtx.Debug("write with limiter", zap.Int("count", i))
	}
}
