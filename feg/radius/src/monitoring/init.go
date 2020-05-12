package monitoring

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/monitoring/census"
	"fbc/cwf/radius/monitoring/ods"
	"fbc/cwf/radius/monitoring/scuba"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Init(config *config.MonitoringConfig, logger *zap.Logger) (*zap.Logger, error) {
	var result *zap.Logger = logger
	var err error

	// Init default configuration values
	if config == nil {
		return nil, errors.New("Could not find 'monitoring' section in configuration")
	}

	if config.Census != nil {
		census.Init(*config.Census, logger)
	}

	if config.Ods != nil {
		ods.Init(config.Ods, logger)
	}

	if config.Scuba != nil {
		scuba.Initialize(config.Scuba, logger)
		result, err = scuba.NewLogger("xwf_goradius")
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
