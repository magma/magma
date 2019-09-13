package coanas

import (
	"context"
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fmt"

	"go.uber.org/zap"

	"github.com/mitchellh/mapstructure"
)

// Config config has only one parameter which is the port to forward the request
type Config struct {
	Port string
}

// ModuleCtx ...
type ModuleCtx struct {
	port string
}

// Init module interface implementation
func Init(_ *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var coaConfig Config
	err := mapstructure.Decode(config, &coaConfig)
	if err != nil {
		return nil, err
	}

	if coaConfig.Port == "" {
		return nil, errors.New("coa module cannot be initialized with empty Port value")
	}

	return ModuleCtx{port: coaConfig.Port}, nil
}

// Handle module interface implementation
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)
	requestCode := r.Code
	// Checking that we have received a coa request
	if requestCode != radius.CodeDisconnectRequest && requestCode != radius.CodeCoARequest {
		return next(c, r)
	}

	// Extract the nas value from the coa package and forward the request
	coaNasAttribute, _ := rfc2865.NASIPAddress_Lookup(r.Packet)
	if coaNasAttribute == nil {
		return next(c, r)
	}

	// Sending the request to the ip specified in the nas attribute
	host := coaNasAttribute.String()
	res, err := radius.Exchange(context.Background(), r.Packet, fmt.Sprintf("%s:%s", host, mCtx.port))
	if err != nil {
		return nil, err
	}

	b, err := res.Encode()
	if err != nil {
		c.Logger.Info("failed to serialize CoA response")
	}

	return &modules.Response{
		Code:       res.Code,
		Attributes: res.Attributes,
		Raw:        b,
	}, nil
}
