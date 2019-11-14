package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"magma/orc8r/cloud/go/metrics"
)

type Calculation interface {
	Calculate(PrometheusAPI) ([]Result, error)
}

type Result struct {
	value      float64
	metricName string
	labels     prometheus.Labels
}

// XAPCalculation holds the parameters needed to run a XAP query and the registered
// prometheus gauge that the resulting value should be stored in
type XAPCalculation struct {
	Days            int
	ThresholdBytes  int
	QueryStepSize   time.Duration
	RegisteredGauge *prometheus.GaugeVec
	Labels          prometheus.Labels
	Name            string
}

// Calculate returns the number of unique users who have had a session in the
// past X days and have used over `thresholdBytes` data in that time
func (x *XAPCalculation) Calculate(prometheusClient PrometheusAPI) ([]Result, error) {
	// List the users who have had an active session over the last X days
	uniqueUsersQuery := fmt.Sprintf(`count(max_over_time(active_sessions[%dd]) >= 1) by (imsi)`, x.Days)
	// List the users who have used at least x.ThresholdBytes of data in the last X days
	usersOverThresholdQuery := fmt.Sprintf(`count(sum(increase(octets_in[%dd])) by (imsi) > %d)`, x.Days, x.ThresholdBytes)
	// Count the users who match both conditions
	intersectionQuery := fmt.Sprintf(`count(%s and %s) by (%s)`, uniqueUsersQuery, usersOverThresholdQuery, metrics.NetworkLabelName)

	val, err := prometheusClient.Query(context.Background(), intersectionQuery, time.Now())
	if err != nil {
		return nil, err
	}
	vec, ok := val.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("XAP query returned unexpected ValueType: %v", val.Type())
	}
	if len(vec) == 0 {
		return nil, fmt.Errorf("no data returned from query")
	}

	var results []Result
	for _, v := range vec {
		// Get labels from query result
		queryLabels := make(map[string]string, 0)
		for label, value := range v.Metric {
			queryLabels[string(label)] = string(value)
		}
		combinedLabels := combineLabels(x.Labels, queryLabels)
		results = append(results, Result{
			metricName: x.Name,
			labels:     combinedLabels,
			value:      float64(v.Value),
		})
		x.RegisteredGauge.With(combinedLabels).Set(float64(vec[0].Value))
	}
	return results, nil
}

func combineLabels(l1, l2 map[string]string) map[string]string {
	retLabels := make(map[string]string)
	for l, v := range l1 {
		retLabels[l] = v
	}
	for l, v := range l2 {
		retLabels[l] = v
	}
	return retLabels
}
