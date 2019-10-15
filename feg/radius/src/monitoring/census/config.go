package census

import (
	ocprom "contrib.go.opencensus.io/exporter/prometheus"
)

type (
	// Config offers a declarative way to construct a census.
	Config struct {
		DisableStats        bool      `env:"NO_STATS" long:"no-stats" description:"Disables statistics gathering and exporting" json:"disable_stats"`
		StatViews           StatViews `env:"VIEWS" long:"view" default:"proc" description:"Set of metric types to expose" json:"stat_views"`
		SamplingProbability float64   `env:"SAMPLING_PROBABILITY" long:"sampling-probability" default:"1.0" description:"Trace sampling probability" json:"sampling_probability"`
	}

	// StatViews attaches flags methods to []string.
	StatViews []string
)

// Option a wrapper for an OC config option
type Option func(*ocprom.Options) error
