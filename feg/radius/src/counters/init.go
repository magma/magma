package counters

import (
	"fbc/lib/go/oc"
	"net/http"

	"go.uber.org/zap"
)

// Init Initialize views and Prometheus exporter
func Init(config oc.Config, logger *zap.Logger) {
	// Create metrics server
	census, err := config.Build(oc.WithLogger(logger))
	if err != nil {
		logger.Error("Failed building census", zap.Error(err))
		return
	}
	http.Handle("/metrics", census.StatsHandler)
	go func() {
		defer census.Close()
		http.ListenAndServe(":9100", nil)
	}()
}
