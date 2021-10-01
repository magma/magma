// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/golang/mock/gomock"
	"github.com/magma/magma/src/go/protos/magma/metricsd"
	"github.com/magma/magma/src/go/protos/magma/metricsd/mock_metricsd"
	"github.com/pkg/errors"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

type metricsContainerMatcher struct{ t *metricsd.MetricsContainer }

func MetricsContainerMatcher(t *metricsd.MetricsContainer) gomock.Matcher {
	return &metricsContainerMatcher{t}
}

func (o *metricsContainerMatcher) Matches(x interface{}) bool {
	return fmt.Sprintf("%+v", x) == fmt.Sprintf("%+v", o.t)
}

func (o *metricsContainerMatcher) String() string {
	return fmt.Sprintf("%+v", o.t)
}

func TestPrometheusMetricSender(t *testing.T) {
	// Create the mock client first
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mClient := mock_metricsd.NewMockMetricsControllerClient(ctrl)

	// New sender
	sender := NewPrometheusMetricSender(mClient, "gw1")

	// Set values for a Gauge
	mClock := clock.NewMock()
	gauge := NewGauge(mClock, "metric1", "test metric", []string{})
	gauge.Set(1, map[string]string{"label1": "asdf"})
	mClock.Add(time.Second)
	gauge.Set(2, map[string]string{})

	// Trigger collection of metric samples
	mClient.EXPECT().Collect(
		gomock.Any(),
		MetricsContainerMatcher(
			&metricsd.MetricsContainer{
				GatewayId: "gw1",
				Family: []*io_prometheus_client.MetricFamily{
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric1"),
						Help: proto.String("test metric"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label1"),
										Value: proto.String("asdf"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(1)},
								TimestampMs: proto.Int64(0),
							},
							&io_prometheus_client.Metric{
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(2)},
								TimestampMs: proto.Int64(1000),
							},
						},
					},
				},
			},
		),
		gomock.Any(),
	).Return(nil, errors.New("gRPC test error"))

	// First attempt a send, and simulate a failure
	gauges := []Gauge{gauge}
	err := sender.Send(gauges)
	assert.NotNil(t, err)

	// Add another gauge and sample
	gauge2 := NewGauge(mClock, "metric2", "test metric 2", []string{})
	mClock.Add(time.Second)
	gauge.Set(4, map[string]string{"label1": "asdf"})
	gauge2.Set(3, map[string]string{"label2": "abcd"})

	// Next time that collection and send is triggered, the previous metrics
	// and new metric samples should be sent
	mClient.EXPECT().Collect(
		gomock.Any(),
		MetricsContainerMatcher(
			&metricsd.MetricsContainer{
				GatewayId: "gw1",
				Family: []*io_prometheus_client.MetricFamily{
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric1"),
						Help: proto.String("test metric"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label1"),
										Value: proto.String("asdf"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(1)},
								TimestampMs: proto.Int64(0),
							},
							&io_prometheus_client.Metric{
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(2)},
								TimestampMs: proto.Int64(1000),
							},
						},
					},
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric2"),
						Help: proto.String("test metric 2"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label2"),
										Value: proto.String("abcd"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(3)},
								TimestampMs: proto.Int64(2000),
							},
						},
					},
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric1"),
						Help: proto.String("test metric"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label1"),
										Value: proto.String("asdf"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(4)},
								TimestampMs: proto.Int64(2000),
							},
						},
					},
				},
			},
		),
		gomock.Any(),
	).Return(nil, nil)

	gauges = []Gauge{gauge, gauge2}
	err = sender.Send(gauges)
	assert.Nil(t, err)

	// On a third collect + send, since the previous send succeeded, we should
	// only send the new samples reported since the last send
	mClock.Add(time.Second)
	gauge.Set(11, map[string]string{"label1": "asdf"})
	gauge2.Set(12, map[string]string{"label2": "abcd"})

	mClient.EXPECT().Collect(
		gomock.Any(),
		MetricsContainerMatcher(
			&metricsd.MetricsContainer{
				GatewayId: "gw1",
				Family: []*io_prometheus_client.MetricFamily{
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric2"),
						Help: proto.String("test metric 2"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label2"),
										Value: proto.String("abcd"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(12)},
								TimestampMs: proto.Int64(3000),
							},
						},
					},
					&io_prometheus_client.MetricFamily{
						Name: proto.String("metric1"),
						Help: proto.String("test metric"),
						Type: io_prometheus_client.MetricType_GAUGE.Enum(),
						Metric: []*io_prometheus_client.Metric{
							&io_prometheus_client.Metric{
								Label: []*io_prometheus_client.LabelPair{
									&io_prometheus_client.LabelPair{
										Name:  proto.String("label1"),
										Value: proto.String("asdf"),
									},
								},
								Gauge:       &io_prometheus_client.Gauge{Value: proto.Float64(11)},
								TimestampMs: proto.Int64(3000),
							},
						},
					},
				},
			},
		),
		gomock.Any(),
	).Return(nil, nil)

	err = sender.Send(gauges)
	assert.Nil(t, err)
}
