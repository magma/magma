/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package eventd_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	"github.com/thoas/go-funk"
)

const (
	pathParamStreamName  = "stream_name"
	pathParamNetworkID   = "network_id"
	queryParamEventType  = "event_type"
	queryParamHardwareID = "hardware_id"
	queryParamTag        = "tag"

	defaultQuerySize = 50

	// We use the ES "keyword" type for exact match
	dotKeyword              = ".keyword"
	elasticFilterStreamName = pathParamStreamName + dotKeyword
	elasticFilterNetworkID  = pathParamNetworkID + dotKeyword
	elasticFilterEventType  = queryParamEventType + dotKeyword
	elasticFilterHardwareID = "hw_id" + dotKeyword
	elasticFilterEventTag   = "event_tag" + dotKeyword // We use event_tag as fluentd uses the "tag" field
	elasticFilterTimestamp  = "@timestamp"
)

// GetElasticClient parses es config and instanciates a new es client
func GetElasticClient() (*elastic.Client, error) {
	elasticConfig, err := config.GetServiceConfig(orc8r.ModuleName, "elastic")
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

// EventQueryParams represents a single stream es query
type EventQueryParams struct {
	StreamName string
	EventType  string
	HardwareID string
	NetworkID  string
	Tag        string
}

func (b *EventQueryParams) toElasticBoolQuery() *elastic.BoolQuery {
	query := elastic.NewBoolQuery()
	query.Filter(elastic.NewTermQuery(elasticFilterStreamName, b.StreamName))
	query.Filter(elastic.NewTermQuery(elasticFilterNetworkID, b.NetworkID))
	if len(b.EventType) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterEventType, b.EventType))
	}
	if len(b.HardwareID) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterHardwareID, b.HardwareID))
	}
	if len(b.Tag) > 0 {
		query.Filter(elastic.NewTermQuery(elasticFilterEventTag, b.Tag))
	}
	return query
}

// MultiStreamEventQueryParams exposes more query options.
// Primarily the ability to query across multiple streams/tags and specifying
// a time range. It also accepts an optional query size limit and offset for
// paginated queries.
type MultiStreamEventQueryParams struct {
	NetworkID   string
	Streams     []string
	Events      []string
	Tags        []string
	HardwareIDs []string
	From        int
	Size        int
	Start       *time.Time
	End         *time.Time
}

func (m *MultiStreamEventQueryParams) toElasticBoolQuery() *elastic.BoolQuery {
	ret := elastic.NewBoolQuery().Filter(elastic.NewTermQuery(elasticFilterNetworkID, m.NetworkID))
	if !funk.IsEmpty(m.Streams) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterStreamName, stringsToInterfaces(m.Streams)...))
	}
	if !funk.IsEmpty(m.Events) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventType, stringsToInterfaces(m.Events)...))
	}
	if !funk.IsEmpty(m.Tags) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterEventTag, stringsToInterfaces(m.Tags)...))
	}
	if !funk.IsEmpty(m.HardwareIDs) {
		ret.Filter(elastic.NewTermsQuery(elasticFilterHardwareID, stringsToInterfaces(m.HardwareIDs)...))
	}
	if m.Start != nil || m.End != nil {
		ret.Must(elastic.NewRangeQuery(elasticFilterTimestamp).From(m.Start).To(m.End))
	}
	return ret
}

// GetEvents query es and return a list of events for a specific stream
func GetEvents(ctx context.Context, params EventQueryParams, client *elastic.Client) ([]models.Event, error) {
	elasticQuery := params.toElasticBoolQuery()
	search := client.Search().
		Index("").
		Size(defaultQuerySize).
		Sort(elasticFilterTimestamp, false).
		Query(elasticQuery)
	return doSearch(ctx, search)
}

//  GetMultiStreamEvents exposes more query options than EventsHandler,
func GetMultiStreamEvents(ctx context.Context, params MultiStreamEventQueryParams, client *elastic.Client) ([]models.Event, error) {
	query := params.toElasticBoolQuery()
	search := client.Search().
		Index("eventd*").
		From(params.From).
		Size(params.Size).
		Sort(elasticFilterTimestamp, true).
		Query(query)
	return doSearch(ctx, search)
}

// GetEventCount queries es and returns the number of events based on the provided query options.
func GetEventCount(ctx context.Context, params MultiStreamEventQueryParams, client *elastic.Client) (int64, error) {
	query := params.toElasticBoolQuery()
	result, err := client.Count().
		Index("eventd*").
		Query(query).
		Do(ctx)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func doSearch(ctx context.Context, search *elastic.SearchService) ([]models.Event, error) {
	result, err := search.Do(ctx)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("Elastic Error Type: %s, Reason: %s", result.Error.Type, result.Error.Reason)
	}

	eventResults, err := GetEventResults(result.Hits.Hits)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return eventResults, nil
}

type eventElasticHit struct {
	StreamName string `json:"stream_name"`
	EventType  string `json:"event_type"`
	// FluentBit logs sent from AGW are tagged with hw_id
	HardwareID string `json:"hw_id"`
	// We use event_tag as fluentd uses the "tag" field
	Tag       string `json:"event_tag"`
	Timestamp string `json:"@timestamp"`
	Value     string `json:"value"`
}

// Retrieve Event properties from the _source of
// ES Hits, including event metadata
func GetEventResults(hits []*elastic.SearchHit) ([]models.Event, error) {
	results := []models.Event{}
	for _, hit := range hits {
		var eventHit eventElasticHit
		// Get Value from the _source
		if err := json.Unmarshal(hit.Source, &eventHit); err != nil {
			return nil, fmt.Errorf("Unable to Unmarshal JSON from elastic.Hit. "+
				"elastic.Hit.Source: %s, Error: %s", hit.Source, err)
		}
		// Skip hits without an event value
		if eventHit.Value == "" {
			return nil, fmt.Errorf("eventResult %s does not contain a value", eventHit)
		}
		var eventValue map[string]interface{}
		if err := json.Unmarshal([]byte(eventHit.Value), &eventValue); err != nil {
			return nil, fmt.Errorf("Unable to Unmarshal JSON from eventResult.Value. "+
				"eventHit.Value: %s, Error: %s", hit.Source, err)
		}
		results = append(results, models.Event{
			StreamName: eventHit.StreamName,
			EventType:  eventHit.EventType,
			HardwareID: eventHit.HardwareID,
			Tag:        eventHit.Tag,
			Timestamp:  eventHit.Timestamp,
			Value:      eventValue,
		})
	}
	return results, nil
}

func stringsToInterfaces(st []string) []interface{} {
	ret := make([]interface{}, 0, len(st))
	for _, s := range st {
		ret = append(ret, s)
	}
	return ret
}
