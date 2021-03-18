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
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/protos"

	"github.com/golang/glog"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var (
	prometheusNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	nonPromoChars       = regexp.MustCompile(`[^a-zA-Z\d_]`)
)

type PushExporterServicer struct {
	PushAddresses []string
}

// NewPushExporterServicer returns an exporter pushing metrics to Prometheus
// pushgateways at the passed addresses.
func NewPushExporterServicer(pushAddrs []string) protos.MetricsExporterServer {
	srv := &PushExporterServicer{
		PushAddresses: ensureHTTP(pushAddrs),
	}
	return srv
}

func (s *PushExporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	return &protos.SubmitMetricsResponse{}, s.pushFamilies(processMetrics(req.GetMetrics()))
}

func processMetrics(metrics []*protos.ContextualizedMetric) []*io_prometheus_client.MetricFamily {
	processedMetrics := make([]*io_prometheus_client.MetricFamily, 0)
	for _, metricAndContext := range metrics {
		// Don't register family if it has 0 metrics. Would cause prometheus scrape
		// to fail.
		if len(metricAndContext.Family.Metric) == 0 {
			continue
		}
		fam := metricAndContext.Family
		fam.Name = sanitizePrometheusName(metricAndContext.Context.MetricName)
		fam.Metric = dropInvalidMetrics(fam.Metric, fam.GetName())
		if len(fam.Metric) == 0 {
			continue
		}
		for _, metric := range fam.Metric {
			if metric.TimestampMs == nil || *metric.TimestampMs == 0 {
				timeStamp := time.Now().Unix() * 1000
				metric.TimestampMs = &timeStamp
			}
		}
		processedMetrics = append(processedMetrics, fam)
	}
	return processedMetrics
}

// dropInvalidMetrics because invalid label names would cause the entire scrape
// to fail. Drop them here and log to allow good metrics through
func dropInvalidMetrics(metrics []*io_prometheus_client.Metric, familyName string) []*io_prometheus_client.Metric {
	validMetrics := make([]*io_prometheus_client.Metric, 0, len(metrics))
	for _, metric := range metrics {
		if err := validateLabels(metric); err != nil {
			glog.Errorf("Dropping metric %s because of invalid label: %v", familyName, err)
		} else {
			validMetrics = append(validMetrics, metric)
		}
	}
	return validMetrics
}

func validateLabels(metric *io_prometheus_client.Metric) error {
	for _, label := range metric.Label {
		if !prometheusNameRegex.MatchString(label.GetName()) {
			return fmt.Errorf("label %s invalid", label.GetName())
		}
	}
	return nil
}

func familyToString(family *io_prometheus_client.MetricFamily) (string, error) {
	var buf bytes.Buffer
	_, err := expfmt.MetricFamilyToText(&buf, family)
	if err != nil {
		return "", fmt.Errorf("error writing family string: %v", err)
	}
	return buf.String(), nil
}

func (s *PushExporterServicer) pushFamilies(fams []*io_prometheus_client.MetricFamily) error {
	if len(fams) == 0 {
		return nil
	}
	builder := strings.Builder{}

	for _, fam := range fams {
		familyString, err := familyToString(fam)
		if err != nil {
			glog.Errorf("Family dropped during push: %s. Error: %v", fam.GetName(), err)
			continue
		}
		builder.WriteString(familyString)
		builder.WriteString("\n")
	}

	body := builder.String()
	client := http.Client{}
	var err error
	for _, address := range s.PushAddresses {
		resp, pushErr := client.Post(address, "text/plain", bytes.NewBufferString(body))
		if pushErr != nil {
			err = fmt.Errorf("%w; Error sending request to push receiver: %s: %s", err, address, pushErr)
			continue
		}
		if resp.StatusCode/100 != 2 {
			respBody, _ := ioutil.ReadAll(resp.Body)
			err = fmt.Errorf("%w; non-200 response code from push receiver: %s: Status: %d, %s", err, address, resp.StatusCode, respBody)
			continue
		}
	}
	return err
}

func ensureHTTP(addrs []string) []string {
	for i, addr := range addrs {
		if !strings.HasPrefix(addr, "http") {
			addrs[i] = fmt.Sprintf("http://%s", addr)
		}
	}
	return addrs
}

func makeStringPointer(str string) *string {
	return &str
}

func sanitizePrometheusName(name string) *string {
	sanitizedName := nonPromoChars.ReplaceAllString(name, "_")
	// If still doesn't match, must be because digit is first character.
	if !prometheusNameRegex.MatchString(sanitizedName) {
		sanitizedName = "_" + sanitizedName
	}
	return &sanitizedName
}
