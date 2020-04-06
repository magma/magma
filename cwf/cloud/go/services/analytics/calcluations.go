package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"magma/orc8r/lib/go/metrics"
)

const (
	APNLabel       = "apn"
	DaysLabel      = "days"
	DirectionLabel = "direction"
)

type Calculation interface {
	Calculate(PrometheusAPI) ([]Result, error)
}

type CalculationParams struct {
	Days            int
	RegisteredGauge *prometheus.GaugeVec
	Labels          prometheus.Labels
	Name            string
}

type Result struct {
	value      float64
	metricName string
	labels     prometheus.Labels
}

// XAPCalculation holds the parameters needed to run a XAP query and the registered
// prometheus gauge that the resulting value should be stored in
type XAPCalculation struct {
	CalculationParams
	ThresholdBytes int
}

// Calculate returns the number of unique users who have had a session in the
// past X days and have used over `thresholdBytes` data in that time
func (x *XAPCalculation) Calculate(prometheusClient PrometheusAPI) ([]Result, error) {
	// List the users who have had an active session over the last X days
	uniqueUsersQuery := fmt.Sprintf(`count(max_over_time(active_sessions[%dd]) >= 1) by (imsi,networkID)`, x.Days)
	// List the users who have used at least x.ThresholdBytes of data in the last X days
	usersOverThresholdQuery := fmt.Sprintf(`count(sum(increase(octets_in[%dd])) by (imsi,networkID) > %d) by (imsi,networkID)`, x.Days, x.ThresholdBytes)
	// Count the users who match both conditions
	intersectionQuery := fmt.Sprintf(`count(%s and %s) by (%s)`, uniqueUsersQuery, usersOverThresholdQuery, metrics.NetworkLabelName)

	vec, err := queryPrometheusVector(prometheusClient, intersectionQuery)
	if err != nil {
		return nil, fmt.Errorf("User Consumption query error: %s", err)
	}

	results := makeVectorResults(vec, x.Labels, x.Name)
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}

type APThroughputCalculation struct {
	CalculationParams
	QueryStepSize time.Duration
	Direction     ConsumptionDirection
}

func (x *APThroughputCalculation) Calculate(prometheusClient PrometheusAPI) ([]Result, error) {
	// Get datapoints for throughput when the value is not 0 segmented by apn
	avgRateQuery := fmt.Sprintf(`avg(rate(octets_%s[3m]) > 0) by (%s, %s)`, x.Direction, APNLabel, metrics.NetworkLabelName)

	timeRange := v1.Range{End: time.Now(), Start: time.Now().Add(-time.Duration(x.Days * int(time.Hour) * 24)), Step: x.QueryStepSize}
	avgRateMatrix, err := queryPrometheusMatrix(prometheusClient, avgRateQuery, timeRange)
	if err != nil {
		return nil, fmt.Errorf("AP Throughput query error: %s", err)
	}

	results := make([]Result, 0)
	for _, apnAverages := range avgRateMatrix {
		apn := string(apnAverages.Metric[APNLabel])
		nID := string(apnAverages.Metric[metrics.NetworkLabelName])
		avgThroughputOverTime := averageDatapoints(apnAverages.Values)
		if apn == "" || nID == "" {
			glog.Errorf("Missing tags from AP Throughput Calculation: APN: %s, NetworkID: %s", apn, nID)
			continue
		}
		results = append(results, Result{
			value:      avgThroughputOverTime,
			metricName: x.Name,
			labels:     combineLabels(x.Labels, map[string]string{APNLabel: apn, metrics.NetworkLabelName: nID, DirectionLabel: string(x.Direction)}),
		})
	}
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}

type ConsumptionDirection string

const (
	ConsumptionIn  ConsumptionDirection = "in"
	ConsumptionOut ConsumptionDirection = "out"
)

type UserThroughputCalculation struct {
	CalculationParams
	QueryStepSize time.Duration
	Direction     ConsumptionDirection
}

func (x *UserThroughputCalculation) Calculate(prometheusClient PrometheusAPI) ([]Result, error) {
	// Get datapoints for throughput when the value is not 0 segmented
	avgRateQuery := fmt.Sprintf(`avg(rate(octets_%s[3m]) > 0) by (%s)`, x.Direction, metrics.NetworkLabelName)

	timeRange := v1.Range{End: time.Now(), Start: time.Now().Add(-time.Duration(x.Days * int(time.Hour) * 24)), Step: x.QueryStepSize}
	avgRateMatrix, err := queryPrometheusMatrix(prometheusClient, avgRateQuery, timeRange)
	if err != nil {
		return nil, fmt.Errorf("User Throughput query error: %s", err)
	}

	results := make([]Result, 0)
	for _, apnAverages := range avgRateMatrix {
		nID := string(apnAverages.Metric[metrics.NetworkLabelName])
		avgThroughputOverTime := averageDatapoints(apnAverages.Values)
		if nID == "" {
			glog.Error("Missing NetworkID from Throughput Calculation")
			continue
		}
		results = append(results, Result{
			value:      avgThroughputOverTime,
			metricName: x.Name,
			labels:     combineLabels(x.Labels, map[string]string{metrics.NetworkLabelName: nID, DirectionLabel: string(x.Direction)}),
		})
	}
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}

type UserConsumptionCalculation struct {
	CalculationParams
	Direction ConsumptionDirection
}

func (x *UserConsumptionCalculation) Calculate(prometheusClient PrometheusAPI) ([]Result, error) {
	consumptionQuery := fmt.Sprintf(`sum(increase(octets_%s[%dd])) by (%s)`, x.Direction, x.Days, metrics.NetworkLabelName)

	vec, err := queryPrometheusVector(prometheusClient, consumptionQuery)
	if err != nil {
		return nil, fmt.Errorf("User Consumption query error: %s", err)
	}

	baseLabels := combineLabels(x.Labels, map[string]string{DirectionLabel: string(x.Direction)})
	results := makeVectorResults(vec, baseLabels, x.Name)
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}

func queryPrometheusVector(prometheusClient PrometheusAPI, query string) (model.Vector, error) {
	// TODO: catch the warning at _
	val, _, err := prometheusClient.Query(context.Background(), query, time.Now())
	if err != nil {
		return nil, err
	}
	vec, ok := val.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected ValueType: %v", val.Type())
	}
	if len(vec) == 0 {
		return nil, fmt.Errorf("no data returned from query")
	}
	return vec, nil
}

func queryPrometheusMatrix(prometheusClient PrometheusAPI, query string, r v1.Range) (model.Matrix, error) {
	// TODO: catch the warning at _
	val, _, err := prometheusClient.QueryRange(context.Background(), query, r)
	if err != nil {
		return nil, err
	}
	matrix, ok := val.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected ValueType: %v", val.Type())
	}
	if len(matrix) == 0 {
		return nil, fmt.Errorf("no data returned from query")
	}
	return matrix, nil
}

func averageDatapoints(samples []model.SamplePair) float64 {
	sum := float64(0)
	for _, val := range samples {
		sum += float64(val.Value)
	}
	return sum / float64(len(samples))
}

func makeVectorResults(vec model.Vector, baseLabels prometheus.Labels, metricName string) []Result {
	var results []Result
	for _, v := range vec {
		// Get labels from query result
		queryLabels := make(map[string]string, 0)
		for label, value := range v.Metric {
			queryLabels[string(label)] = string(value)
		}
		combinedLabels := combineLabels(baseLabels, queryLabels)
		results = append(results, Result{
			metricName: metricName,
			labels:     combinedLabels,
			value:      float64(v.Value),
		})
	}
	return results
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
