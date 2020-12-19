package calculations

import (
	"context"
	"fmt"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/service/config"
	"time"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	sortTag         = "@timestamp"
	defaultLogField = "message"
)

func GetElasticClient() *elastic.Client {
	elasticConfig, err := config.GetServiceConfig(orc8r.ModuleName, "elastic")
	if err != nil {
		glog.Errorf("Error %v reading elastic service configuration", err)
		return nil
	}
	elasticHost := elasticConfig.MustGetString("elasticHost")
	elasticPort := elasticConfig.MustGetInt("elasticPort")

	client, err := elastic.NewSimpleClient(elastic.SetURL(fmt.Sprintf("http://%s:%d", elasticHost, elasticPort)))
	if err != nil {
		glog.Errorf("Error %v getting client handle to elastic service", err)
		return nil
	}
	return client
}

func elasticQuery(hours int, logConfig *LogConfig) *elastic.BoolQuery {
	query := elastic.NewBoolQuery()

	startTime := time.Now().Add(time.Duration(-hours) * time.Hour).UnixNano()
	startTime = startTime / int64(time.Millisecond)
	timeRangeQuery := elastic.NewRangeQuery(sortTag)
	timeRangeQuery.Gte(startTime)
	query.Must(timeRangeQuery)

	simpleQuery := elastic.NewSimpleQueryStringQuery(logConfig.Query).AnalyzeWildcard(true)
	if logConfig.Fields != nil {
		for _, field := range logConfig.Fields {
			simpleQuery.Field(field)
		}
	} else {
		simpleQuery.Field(defaultLogField)
	}
	query.Filter(simpleQuery)

	for key, value := range logConfig.Tags {
		query.Filter(elastic.NewTermQuery(key, value))
	}
	return query
}

//LogsMetricCalculation defines new metric calculation based on querying the orc8r logs
type LogsMetricCalculation struct {
	BaseCalculation
	LogConfig     *LogConfig
	ElasticClient *elastic.Client
}

func (x *LogsMetricCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculating Log Metrics")
	results := []*protos.CalculationResult{}
	if x.ElasticClient == nil {
		err := fmt.Errorf("Elastic client not found for LogMetricsCalculation")
		return results, err
	}

	q := elasticQuery(x.Hours, x.LogConfig)
	logCount, err := x.ElasticClient.Count().Index("").Query(q).Do(context.Background())
	if err != nil {
		glog.Errorf("Error %v querying elastic search for query %v ", err, q)
		return nil, err
	}

	// add metric
	if x.Labels == nil {
		x.Labels = make(map[string]string)
	}
	results = append(results, NewResult(float64(logCount), x.Name, x.Labels))
	glog.V(1).Infof("Log Metrics results %v", results)
	return results, nil
}

// GetLogMetricsCalculations gets all the log calculations based on the provided log configs
func GetLogMetricsCalculations(elasticClient *elastic.Client, analyticsConfig *AnalyticsConfig) []Calculation {
	allCalculations := make([]Calculation, 0)

	for metricName, metricConfig := range analyticsConfig.Metrics {
		if metricConfig.LogConfig != nil {
			continue
		}

		labels := []string{}
		for k := range metricConfig.Labels {
			labels = append(labels, k)
		}
		glog.V(1).Infof("Adding Log Calculation for %s", metricName)
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricName}, labels)
		allCalculations = append(allCalculations, &LogsMetricCalculation{
			BaseCalculation: BaseCalculation{
				&CalculationParams{
					Name:                metricName,
					Hours:               3,
					AnalyticsConfig:     analyticsConfig,
					Labels:              metricConfig.Labels,
					RegisteredGauge:     gauge,
					ExpectedGaugeLabels: labels,
				}},
			LogConfig:     metricConfig.LogConfig,
			ElasticClient: elasticClient,
		})
	}
	return allCalculations
}
