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

package servicers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"magma/fbinternal/cloud/go/metrics/ods"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
)

const (
	ODSMetricsExportInterval = time.Second * 15
	// ODSMetricsQueueLength derivation:
	// a sample is 10 bytes
	// right now 50 metrics from each gateway, 35 metrics from each cloud
	// service per minute assume we support 100 metrics from each gateway,
	// 70 metrics from each cloud service. with 1000 gws, we will have 100000
	// metrics per minute from gws. with 30 cloud services,
	// we have 2100 from cloud.
	// this needs 10 * 102100 = 1021000 B
	ODSMetricsQueueLength = 102000

	deviceIDLabelName = "deviceID"
	serviceLabelName  = "service"
	tagsLabelName     = "tags"
)

var (
	// odsForbiddenKeyLabels lists labels that should not be appended to the ODS key
	odsForbiddenKeyLabels = []string{
		metrics.CloudHostLabelName,
		metrics.NetworkLabelName,
		metrics.GatewayLabelName,
		deviceIDLabelName,
		tagsLabelName,
	}
)

type HTTPClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type ExporterServicer struct {
	odsURL         string
	categoryID     string
	queue          []exporters.Sample
	queueMu        sync.Mutex
	maxQueueLength int
	exportInterval time.Duration
	metricsPrefix  string
}

// NewExporterServicer returns a MetricExporter with ODS datasink.
func NewExporterServicer(
	baseUrl string,
	appId string,
	appSecret string,
	categoryID string,
	metricsPrefix string,
	maxQueueLength int,
	exportInterval time.Duration,
) protos.MetricsExporterServer {
	srv := &ExporterServicer{
		odsURL:         fmt.Sprintf("%s?access_token=%s|%s", baseUrl, appId, appSecret),
		categoryID:     categoryID,
		maxQueueLength: maxQueueLength,
		exportInterval: exportInterval,
		metricsPrefix:  metricsPrefix,
	}
	go srv.exportEvery()
	return srv
}

func (s *ExporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	ret := &protos.SubmitMetricsResponse{}
	if len(req.GetMetrics()) == 0 {
		return ret, nil
	}

	var convertedSamples []exporters.Sample
	for _, metricAndContext := range req.GetMetrics() {
		for _, metric := range metricAndContext.Family.GetMetric() {
			newSamples := exporters.GetSamplesForMetrics(metricAndContext, metric)
			convertedSamples = append(convertedSamples, newSamples...)
		}
	}

	s.queueMu.Lock()
	defer s.queueMu.Unlock()

	// Don't append more samples than available queue space.
	endAppendIdx := int(math.Min(float64(len(convertedSamples)), float64(s.maxQueueLength-len(s.queue))))
	if endAppendIdx > 0 {
		s.queue = append(s.queue, convertedSamples[:endAppendIdx]...)
	}
	if endAppendIdx < len(convertedSamples) {
		droppedSampleCount := len(convertedSamples) - endAppendIdx
		return ret, fmt.Errorf("ODS queue full, dropping %d samples", droppedSampleCount)
	}
	return ret, nil
}

// Export syncs metrics in the exporter's queue to ODS. If Export fails, the
// exporter's queue will still be cleared (i.e. the samples will be dropped).
func (s *ExporterServicer) Export(client HTTPClient) error {
	s.queueMu.Lock()
	samples := s.queue
	s.queue = []exporters.Sample{}
	s.queueMu.Unlock()

	if len(samples) != 0 {
		err := s.write(client, samples)
		if err != nil {
			return fmt.Errorf("failed to sync to ODS: %s", err)
		}
	}
	return nil
}

func (s *ExporterServicer) exportEvery() {
	client := &http.Client{Timeout: 30 * time.Second}
	for range time.Tick(s.exportInterval) {
		err := s.Export(client)
		if err != nil {
			glog.Errorf("Error in syncing to ods: %v", err)
		}
	}
}

// Write to ODS from queued samples or error
func (s *ExporterServicer) write(client HTTPClient, samples []exporters.Sample) error {
	var datapoints []ods.Datapoint
	for _, sample := range samples {
		key := FormatKey(sample)
		entity := FormatEntity(sample, s.metricsPrefix)
		datapoints = append(datapoints, ods.Datapoint{
			Entity: entity,
			Key:    key,
			Value:  sample.Value(),
			Tags:   GetTags(sample),
			Time:   int(sample.TimestampMs()),
		})
	}

	datapointsJson, err := json.Marshal(datapoints)
	if err != nil {
		return err
	}
	resp, err := client.PostForm(s.odsURL, url.Values{"datapoints": {string(datapointsJson)}, "category": {s.categoryID}})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return errors.New(string(errMsg))
	}
	return err
}

// FormatKey generates an entity name for the Sample for use with ODS.
// The entity name is a dot separated concatenation of sorted label
// key value pairs.
func FormatKey(sample exporters.Sample) string {
	var keyBuffer bytes.Buffer
	var prefixBuffer bytes.Buffer // stores the service when found
	for _, labelPair := range sample.Labels() {
		// Don't add tags or deviceID label to key
		if labelIsForbidden(labelPair.GetName(), odsForbiddenKeyLabels) {
			continue
		}
		if labelPair.GetName() == serviceLabelName {
			prefixBuffer.WriteString(labelPair.GetValue())
			prefixBuffer.WriteString(".")
		} else {
			keyBuffer.WriteString(".")
			keyBuffer.WriteString(labelPair.GetName())
			keyBuffer.WriteString("-")
			keyBuffer.WriteString(labelPair.GetValue())
		}
	}
	// return combined strings
	prefixBuffer.WriteString(sample.Name())
	prefixBuffer.Write(keyBuffer.Bytes())
	return prefixBuffer.String()
}

// FormatEntity handles the special case of device metrics and appends deviceID
// to the entity if it is present.
func FormatEntity(sample exporters.Sample, metricsPrefix string) string {
	baseEntity := fmt.Sprintf("%s.%s", metricsPrefix, sample.Entity())
	deviceID := ""
	for _, labelPair := range sample.Labels() {
		if labelPair.GetName() == deviceIDLabelName {
			deviceID = labelPair.GetValue()
			break
		}
	}
	if deviceID != "" {
		return fmt.Sprintf("%s.%s", baseEntity, deviceID)
	}
	return baseEntity
}

// GetTags parses labels for tags and splits them from a comma-separated string.
func GetTags(sample exporters.Sample) []string {
	for _, labelPair := range sample.Labels() {
		if labelPair.GetName() == tagsLabelName {
			return strings.Split(labelPair.GetValue(), ",")
		}
	}
	return []string{}
}

func labelIsForbidden(labelName string, forbiddenLabels []string) bool {
	for _, forbidden := range forbiddenLabels {
		if labelName == forbidden {
			return true
		}
	}
	return false
}
