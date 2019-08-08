package coadynamic

import (
	"context"
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/coa_dynamic/radiustracker"
	"fbc/lib/go/radius"
	"fmt"

	"go.uber.org/zap"

	"github.com/mitchellh/mapstructure"
)

var radiusTracker radiustracker.RadiusTracker

// Config config has only one parameter which is the ip to forward the request
type Config struct {
	Port string
}

var port string

// Init module interface implementation
func Init(_ *zap.Logger, config modules.ModuleConfig) error {
	radiusTracker = radiustracker.NewRadiusTracker()

	var coaConfig Config
	err := mapstructure.Decode(config, &coaConfig)
	if err != nil {
		return err
	}

	if coaConfig.Port == "" {
		return errors.New("coa module cannot be initialized with empty port value")
	}

	port = coaConfig.Port
	return nil
}

// Handle module interface implementation
// For radius requests we try to match the called and calling fields to the latest tracked called,calling,ip
// For non coa radius requests we store a mapping of called,calling and ip
func Handle(c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {

	// We received a coa request
	if r.Code == radius.CodeCoARequest || r.Code == radius.CodeDisconnectRequest {
		target, err := radiusTracker.Get(r)
		if err != nil {
			return nil, err
		}

		destination := fmt.Sprintf("%s:%s", target, port)
		res, err := radius.Exchange(context.Background(), r.Packet, destination)
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
		response := &modules.Response{
			Code:       res.Code,
			Attributes: res.Attributes,
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
