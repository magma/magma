/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"context"
	"fmt"

	"magma/lte/cloud/go/services/nprobe"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/eventd/obsidian/handlers"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"
	orc8r_config "magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	"github.com/thoas/go-funk"
)

// multi-stream endpoint query args
const (
	defaultQuerySize = 50

	sortTag             = "@timestamp"
	pathParamStreamName = "stream_name"
	queryParamEventType = "event_type"

	// We use the ES "keyword" type for exact match
	dotKeyword              = ".keyword"
	elasticFilterStreamName = pathParamStreamName + dotKeyword
	elasticFilterEventType  = queryParamEventType + dotKeyword
	elasticFilterEventTag   = "event_tag" + dotKeyword // We use event_tag as fluentd uses the "tag" field
	elasticFilterTimestamp  = "@timestamp"
)

// multi-stream query parameters
type multiStreamEventQueryParams struct {
	networkID string
	streams   []string
	events    []string
	tags      []string
	timestamp string
}

func getElasticClient() (*elastic.Client, error) {
	// Retrieve elastic config
	elasticConfig, err := orc8r_config.GetServiceConfig(orc8r.ModuleName, "elastic")
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate elastic config")
	}
	elasticHost := elasticConfig.MustGetString("elasticHost")
	elasticPort := elasticConfig.MustGetInt("elasticPort")

	client, err := elastic.NewSimpleClient(elastic.SetURL(fmt.Sprintf("http://%s:%d", elasticHost, elasticPort)))
	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate elastic client")
	}
	return client, nil
}

// getMultiStreamsQueryParameters takes a networkID, timestamp and list of tags
// and return a multi streams query parameters
func getMultiStreamsQueryParameters(networkID, timestamp string, tags []string) multiStreamEventQueryParams {
	return multiStreamEventQueryParams{
		networkID: networkID,
		streams:   nprobe.GetESStreams(),
		events:    nprobe.GetESEventTypes(),
		tags:      tags,
		timestamp: timestamp,
	}
}

// EventsCollector wraps an elastic client, query multiple streams per
// subscriber and retrieve all events since last timestamp
type EventsCollector struct {
	ElasticClient *elastic.Client
}

// NewEventsCollector validates params and returns a new EventsCollector.
func NewEventsCollector() (*EventsCollector, error) {
	client, err := getElasticClient()
	if err != nil {
		return nil, err
	}
	return &EventsCollector{
		ElasticClient: client,
	}, nil
}

// GetMultiStreamsEvents queries elastic search with a multi stream event
// query and returns a list of event
func (e *EventsCollector) GetMultiStreamsEvents(networkID, timestamp string, tags []string) ([]models.Event, error) {
	glog.V(3).Info("Collecting events for subscribers ", tags)

	queryParams := getMultiStreamsQueryParameters(networkID, timestamp, tags)
	elasticQuery := queryParams.toElasticBoolQuery()
	search := e.ElasticClient.Search().
		Index("").
		Size(defaultQuerySize).
		Sort(elasticFilterTimestamp, false).
		Query(elasticQuery)

	result, err := search.Do(context.Background())
	if err != nil {
		glog.Errorf("Error getting response from Elastic: %s", err)
		return nil, err
	}
	if result.Error != nil {
		glog.Errorf("Error getting response from Elastic: %v", result.Error)
		return nil, fmt.Errorf("Result error")
	}

	eventResults, err := handlers.GetEventResults(result.Hits.Hits)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return eventResults, nil
}

func stringsToInterfaces(st []string) []interface{} {
	ret := make([]interface{}, 0, len(st))
	for _, s := range st {
		ret = append(ret, s)
	}
	return ret
}

func (m multiStreamEventQueryParams) toElasticBoolQuery() *elastic.BoolQuery {
	ret := elastic.NewBoolQuery()
	if !funk.IsEmpty(m.streams) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterStreamName, stringsToInterfaces(m.streams)...))
	}
	if !funk.IsEmpty(m.events) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventType, stringsToInterfaces(m.events)...))
	}
	if !funk.IsEmpty(m.tags) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventTag, stringsToInterfaces(m.tags)...))
	}
	if m.timestamp != "" {
		timeRangeQuery := elastic.NewRangeQuery(sortTag).Format("strict_date_optional_time_nanos")
		timeRangeQuery.Gt(m.timestamp)
		ret.Must(timeRangeQuery)
	}
	return ret
}
