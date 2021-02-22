package calculations

import (
	"context"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"time"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

const (
	sortTag         = "@timestamp"
	defaultLogField = "log"
)

// LogsMetricCalculation defines new metric calculation based on querying
// elastic logs.
type LogsMetricCalculation struct {
	BaseCalculation
	LogConfig     *LogConfig
	ElasticClient *elastic.Client
}

// Calculate runs elastic count queries based on the input log config and
// exports the result as a metric with input labels
func (x *LogsMetricCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculating Log Metrics")
	var results []*protos.CalculationResult

	if x.ElasticClient == nil {
		err := errors.New("Elastic client not found for LogMetricsCalculation")
		return results, err
	}

	q := elasticQuery(x.Hours, x.LogConfig)
	logCount, err := x.ElasticClient.Count().Index("").Query(q).Do(context.Background())
	if err != nil {
		err = errors.Wrapf(err, "Error %v querying elastic search for query %v", err, q)
		return nil, err
	}

	results = append(results, NewResult(float64(logCount), x.Name, x.Labels))
	glog.V(1).Infof("Log Metrics results %v", results)
	return results, nil
}

func elasticQuery(hours uint, logConfig *LogConfig) *elastic.BoolQuery {
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
