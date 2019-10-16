package coadynamic

import (
	"context"
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/coadynamic/radiustracker"
	"fbc/lib/go/radius"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/mitchellh/mapstructure"
)

var radiusTracker radiustracker.RadiusTracker

// Config config has only one parameter which is the ip to forward the request
type Config struct {
	Port           int16
	TimeoutSeconds uint
}

// ModuleCtx Context for the Module
type ModuleCtx struct {
	port    int16
	timeout uint
}

const DEFAULT_COA_TIMEOUT_SECONDS uint = 5

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	radiusTracker = radiustracker.NewRadiusTracker()

	var coaConfig Config
	err := mapstructure.Decode(config, &coaConfig)
	if err != nil {
		return nil, err
	}

	if coaConfig.Port == 0 {
		return nil, errors.New("coa module cannot be initialized with empty port value")
	}

	if coaConfig.TimeoutSeconds == 0 {
		coaConfig.TimeoutSeconds = DEFAULT_COA_TIMEOUT_SECONDS
	}
	logger.Debug("setting timeout for CoA", zap.Uint("value", coaConfig.TimeoutSeconds))

	return ModuleCtx{
		port:    coaConfig.Port,
		timeout: coaConfig.TimeoutSeconds,
	}, nil
}

// Handle module interface implementation
// For radius requests we try to match the called and calling fields to the latest tracked called,calling,ip
// For non coa radius requests we store a mapping of called,calling and ip
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	mod := m.(ModuleCtx)
	// We received a coa request
	if r.Code == radius.CodeCoARequest || r.Code == radius.CodeDisconnectRequest {
		target, err := radiusTracker.Get(r)
		if err != nil {
			return nil, err
		}

		destination := fmt.Sprintf("%s:%d", target, mod.port)
		ctx, dispose := context.WithTimeout(context.Background(), time.Second*time.Duration(mod.timeout))
		defer dispose()
		res, err := radius.Exchange(ctx, r.Packet, destination)
		if err != nil {
			c.Logger.Debug(
				"failed sending CoA",
				zap.Int("code", int(r.Code)),
				zap.String("dest", destination),
				zap.Error(err),
			)
			return nil, err
		}

		// Send response back
		b, err := res.Encode()
		if err != nil {
			c.Logger.Info("failed to serialize CoA response")
		}
		response := &modules.Response{
			Code:       res.Code,
			Attributes: res.Attributes,
			Raw:        b,
		}
		c.Logger.Debug(
			"successfully sent CoA",
			zap.Int("code", int(r.Code)),
			zap.String("dest", destination),
			zap.Any("response", response),
		)
		return response, nil
	}

	// Regular non coa request
	err := radiusTracker.Set(r)
	if nil != err {
		c.Logger.Error("unable to track the packet", zap.Error(err))
	}
	return next(c, r)
}

// GetRadiusTracker gets the radius tracker associated with this instance
func GetRadiusTracker() radiustracker.RadiusTracker {
	return radiusTracker
}
