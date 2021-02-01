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

package servicers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/metrics/ods"
	"magma/fbinternal/cloud/go/services/fbinternal"
	"magma/fbinternal/cloud/go/services/fbinternal/servicers"
	"magma/fbinternal/cloud/go/services/fbinternal/servicers/mocks"
	"magma/fbinternal/cloud/go/services/fbinternal/test_init"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/test_common"
	"magma/orc8r/lib/go/metrics"

	prometheus_models "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestODSSubmit(t *testing.T) {
	srv := servicers.NewExporterServicer(
		"",
		"",
		"",
		"magma",
		"magma",
		2,
		time.Second*10,
	)
	test_init.StartTestServiceInternal(t, srv)
	exporter := exporters.NewRemoteExporter(fbinternal.ServiceName)

	singleMetricTestFamily := test_common.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, []*prometheus_models.LabelPair{})
	metricContext := exporters.MetricContext{MetricName: "test", AdditionalContext: &exporters.GatewayMetricContext{NetworkID: "testId1", GatewayID: "testId2"}}
	err := exporter.Submit([]exporters.MetricAndContext{{Family: singleMetricTestFamily, Context: metricContext}})
	assert.NoError(t, err)

	// Submitting to a full queue should drop metrics
	multiMetricTestFamily := test_common.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 100, []*prometheus_models.LabelPair{})
	err = exporter.Submit([]exporters.MetricAndContext{
		{Family: multiMetricTestFamily, Context: metricContext},
		{Family: singleMetricTestFamily, Context: metricContext},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ODS queue full, dropping 100 samples")

	err = exporter.Submit([]exporters.MetricAndContext{{Family: singleMetricTestFamily, Context: metricContext}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ODS queue full, dropping 1 samples")
}

func TestExport(t *testing.T) {
	srv := servicers.NewExporterServicer(
		"",
		"",
		"",
		"100",
		"magma",
		2,
		time.Second*10,
	)
	test_init.StartTestServiceInternal(t, srv)
	exporter := exporters.NewRemoteExporter(fbinternal.ServiceName)
	exporterSrv := srv.(*servicers.ExporterServicer)

	entity := "testId1.testId2"
	nameStr := "test"
	tagLabelPair := prometheus_models.LabelPair{
		Name:  test_common.MakeStrPtr("tags"),
		Value: test_common.MakeStrPtr("Tag1,Tag2"),
	}
	sample := exporters.NewSample(nameStr, "0", int64(0), []*prometheus_models.LabelPair{&tagLabelPair}, entity)

	client := &mocks.HTTPClient{}
	resp := &http.Response{StatusCode: 200}
	var datapoints []ods.Datapoint
	datapoints = append(datapoints, ods.Datapoint{
		Entity: fmt.Sprintf("magma.%s.%s", "testId1", "testId2"),
		Key:    servicers.FormatKey(sample),
		Tags:   servicers.GetTags(sample),
		Value:  sample.Value()})
	datapointsJson, err := json.Marshal(datapoints)
	assert.NoError(t, err)

	expectedDatapoints := `[{"entity":"magma.testId1.testId2","key":"test","value":"0","time":0,"tags":["Tag1","Tag2"]}]`
	assert.Equal(t, expectedDatapoints, string(datapointsJson))

	client.On("PostForm", mock.AnythingOfType("string"), url.Values{"datapoints": {string(datapointsJson)}, "category": {"100"}}).Return(resp, nil)

	// Export called on empty queue
	err = exporterSrv.Export(client)
	assert.NoError(t, err)
	client.AssertNotCalled(t, "PostForm", mock.AnythingOfType("string"), mock.AnythingOfType("url.Values"))

	testFamily := test_common.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 1, []*prometheus_models.LabelPair{&tagLabelPair})
	context := exporters.MetricContext{MetricName: "test", AdditionalContext: &exporters.GatewayMetricContext{NetworkID: "testId1", GatewayID: "testId2"}}
	err = exporter.Submit([]exporters.MetricAndContext{{Family: testFamily, Context: context}})
	assert.NoError(t, err)

	err = exporterSrv.Export(client)
	assert.NoError(t, err)
	client.AssertExpectations(t)

	// Fill queue (drop some samples), assert we didn't exceed queue length cap
	multiTestFamily := test_common.MakeTestMetricFamily(prometheus_models.MetricType_GAUGE, 100, []*prometheus_models.LabelPair{&tagLabelPair})
	err = exporter.Submit([]exporters.MetricAndContext{{Family: multiTestFamily, Context: context}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ODS queue full, dropping 98 samples")

	// We expect 2 samples to be exported
	datapoints = append(datapoints, datapoints...)
	datapointsJson, err = json.Marshal(datapoints)
	assert.NoError(t, err)
	client.On("PostForm", mock.Anything, url.Values{"datapoints": {string(datapointsJson)}, "category": {"100"}}).Return(resp, nil)

	err = exporterSrv.Export(client)
	assert.NoError(t, err)
	client.AssertExpectations(t)
}

func TestFormatKey(t *testing.T) {
	entity := "testId1.testId2"
	// Test where key should be prepended with service name
	testSampleWithService := exporters.NewSample(
		"test_sample_with_service",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("service"),
				Value: test_common.MakeStrPtr("mme"),
			},
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)

	key := servicers.FormatKey(testSampleWithService)
	assert.Equal(t, key, "mme.test_sample_with_service.result-success.cause-foo")

	// Test where no service name provided
	testSampleNoService := exporters.NewSample(
		"test_sample_no_service",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)

	key = servicers.FormatKey(testSampleNoService)
	assert.Equal(t, key, "test_sample_no_service.result-success.cause-foo")

	// Test where tags are provided and no service name is provided.
	testSampleWithTags := exporters.NewSample(
		"test_sample_with_tags",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("results"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
			{
				Name:  test_common.MakeStrPtr("tags"),
				Value: test_common.MakeStrPtr("Magma"),
			},
		},
		entity,
	)

	key = servicers.FormatKey(testSampleWithTags)
	assert.Equal(t, key, "test_sample_with_tags.results-success.cause-foo")

	// Test where both tags and service name is provided.
	testSampleWithTagsAndService := exporters.NewSample(
		"test_sample_with_tags_and_service",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("service"),
				Value: test_common.MakeStrPtr("mme"),
			},
			{
				Name:  test_common.MakeStrPtr("results"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
			{
				Name:  test_common.MakeStrPtr("tags"),
				Value: test_common.MakeStrPtr("Magma"),
			},
		},
		entity,
	)

	key = servicers.FormatKey(testSampleWithTagsAndService)
	assert.Equal(t, key, "mme.test_sample_with_tags_and_service.results-success.cause-foo")

	// Test where deviceID is provided
	testSampleDeviceID := exporters.NewSample(
		"test_sample_deviceID",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("deviceID"),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)

	key = servicers.FormatKey(testSampleDeviceID)
	assert.Equal(t, key, "test_sample_deviceID.result-success")

	// Test where cloudHost, gatewayID, and networkID are provided
	testForbiddenLabels := exporters.NewSample(
		"test_sample",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr(metrics.CloudHostLabelName),
				Value: test_common.MakeStrPtr("foo"),
			},
			{
				Name:  test_common.MakeStrPtr(metrics.GatewayLabelName),
				Value: test_common.MakeStrPtr("foo"),
			},
			{
				Name:  test_common.MakeStrPtr(metrics.NetworkLabelName),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)

	key = servicers.FormatKey(testForbiddenLabels)
	assert.Equal(t, key, "test_sample.result-success")
}

func TestFormatTags(t *testing.T) {
	entity := "testId1.testId2"
	// Test where key should be prepended with service name
	testSampleWithService := exporters.NewSample(
		"test_sample_with_service",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("service"),
				Value: test_common.MakeStrPtr("mme"),
			},
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
			{
				Name:  test_common.MakeStrPtr("tags"),
				Value: test_common.MakeStrPtr("Magma,Bootcamp"),
			},
		},
		entity,
	)

	tag := servicers.GetTags(testSampleWithService)
	assert.Equal(t, tag, []string{"Magma", "Bootcamp"})
}

func TestFormatEntity(t *testing.T) {
	// Test entity with no deviceID
	entity := "testId1.testId2"
	prefix := "magma"

	// Test sample with no device ID
	testSampleNoDeviceID := exporters.NewSample(
		"test_sample_no_deviceID",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("cause"),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)
	assert.Equal(t, servicers.FormatEntity(testSampleNoDeviceID, prefix), fmt.Sprintf("magma.%s", entity))

	// Test sample with deviceID
	testSampleDeviceID := exporters.NewSample(
		"test_sample_deviceID",
		"val",
		int64(0),
		[]*prometheus_models.LabelPair{
			{
				Name:  test_common.MakeStrPtr("result"),
				Value: test_common.MakeStrPtr("success"),
			},
			{
				Name:  test_common.MakeStrPtr("deviceID"),
				Value: test_common.MakeStrPtr("foo"),
			},
		},
		entity,
	)
	assert.Equal(t, servicers.FormatEntity(testSampleDeviceID, prefix), fmt.Sprintf("magma.%s.%s", entity, "foo"))
}
