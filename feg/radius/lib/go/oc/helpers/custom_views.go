package helpers

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	// MLatencyMs measures the latency in milliseconds
	MLatencyMs = stats.Int64("latency", "The latency in milliseconds per request", stats.UnitMilliseconds)

	// MErrors measures the amount of general errors.
	MErrors = stats.Int64("errors", "The number of failed events encountered", stats.UnitDimensionless)

	// MSuccess measures the amount of general successful events.
	MSuccess = stats.Int64("success", "The number of successful events encountered", stats.UnitDimensionless)

	// MEvents measures the number of times an event occurred.
	MEvents = stats.Int64("events", "The number of times an event occurred", stats.UnitDimensionless)
)

// TagKeys for the stats quickstart.
var (
	KeyComponent, _ = tag.NewKey("component")
	KeyError, _     = tag.NewKey("error_code")
	KeyOperation, _ = tag.NewKey("operation")
)

// Views for the stats quickstart.
var (
	LatencyView = &view.View{
		Name:        "latency",
		Measure:     MLatencyMs,
		Description: "operation latency in ms",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{KeyComponent, KeyOperation}}

	ErrorCountView = &view.View{
		Name:        "errors",
		Measure:     MErrors,
		Description: "The number of failed events encountered",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent, KeyOperation, KeyError}}

	SuccessCountView = &view.View{
		Name:        "success",
		Measure:     MSuccess,
		Description: "The number of successful events encountered",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent, KeyOperation}}
	CountView = &view.View{
		Name:        "counters",
		Measure:     MEvents,
		Description: "The number of times an event occurred",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{KeyComponent, KeyOperation},
	}
)
