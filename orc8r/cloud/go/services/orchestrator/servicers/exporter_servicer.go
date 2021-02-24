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
	"sync"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/protos"

	"github.com/golang/glog"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

const (
	pushInterval = time.Second * 30
)

var (
	prometheusNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	nonPromoChars       = regexp.MustCompile(`[^a-zA-Z\d_]`)
)

type PushExporterServicer struct {
	FamiliesByName map[string]*io_prometheus_client.MetricFamily
	ExportInterval time.Duration
	PushAddresses  []string
	sync.Mutex
}

// NewPushExporterServicer returns an exporter pushing metrics to Prometheus
// pushgateways at the passed addresses.
func NewPushExporterServicer(pushAddrs []string) protos.MetricsExporterServer {
	srv := &PushExporterServicer{
		FamiliesByName: make(map[string]*io_prometheus_client.MetricFamily),
		ExportInterval: pushInterval,
		PushAddresses:  ensureHTTP(pushAddrs),
	}
	go srv.exportEvery()
	return srv
}

func (s *PushExporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	s.Lock()
	defer s.Unlock()

	processedMetrics := processMetrics(req.GetMetrics())
	for _, family := range processedMetrics {
		familyName := family.GetName()
		if baseFamily, ok := s.FamiliesByName[familyName]; ok {
			addMetricsToFamily(baseFamily, family)
		} else {
			s.FamiliesByName[familyName] = family
		}
	}
	return &protos.SubmitMetricsResponse{}, nil
}

func processMetrics(metrics []*protos.ContextualizedMetric) []*io_prometheus_client.MetricFamily {
	processedMetrics := make([]*io_prometheus_client.MetricFamily, 0)
	for _, metricAndContext := range metrics {
		// Don't register family if it has 0 metrics. Would cause prometheus scrape
		// to fail.
		if len(metricAndContext.Family.Metric) == 0 {
			continue
		}
		originalFamily := metricAndContext.Family
		originalFamily.Name = sanitizePrometheusName(metricAndContext.Context.MetricName)
		// Convert all families to gauges to avoid name collisions of different
		// types.
		convertedFamilies := convertFamilyToGauges(originalFamily)
		for _, fam := range convertedFamilies {
			familyName := fam.GetName()
			fam.Metric = dropInvalidMetrics(fam.Metric, familyName)
			// if all metrics from this family were dropped, don't submit it
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
	}
	return processedMetrics
}

func (s *PushExporterServicer) exportEvery() {
	for range time.Tick(s.ExportInterval) {
		errs := s.pushFamilies()
		s.resetFamilies()
		if len(errs) > 0 {
			glog.Errorf("error pushing to pushgateway: %v", errs)
		}
	}
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

func addMetricsToFamily(baseFamily *io_prometheus_client.MetricFamily, newFamily *io_prometheus_client.MetricFamily) {
	baseFamily.Metric = append(baseFamily.Metric, newFamily.Metric...)
}

func familyToString(family *io_prometheus_client.MetricFamily) (string, error) {
	var buf bytes.Buffer
	_, err := expfmt.MetricFamilyToText(&buf, family)
	if err != nil {
		return "", fmt.Errorf("error writing family string: %v", err)
	}
	return buf.String(), nil
}

func (s *PushExporterServicer) pushFamilies() []error {
	var errs []error
	if len(s.FamiliesByName) == 0 {
		return []error{}
	}
	builder := strings.Builder{}

	s.Lock()
	for _, fam := range s.FamiliesByName {
		familyString, err := familyToString(fam)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		builder.WriteString(familyString)
		builder.WriteString("\n")
	}
	s.Unlock()

	body := builder.String()
	client := http.Client{}
	for _, address := range s.PushAddresses {
		resp, err := client.Post(address, "text/plain", bytes.NewBufferString(body))
		if err != nil {
			errs = append(errs, fmt.Errorf("error sending request to pushgateway %s: %v", address, err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			respBody, _ := ioutil.ReadAll(resp.Body)
			errs = append(errs, fmt.Errorf("non-200 response code from pushgateway %s: %s", address, respBody))
			continue
		}
	}
	return errs
}

func (s *PushExporterServicer) resetFamilies() {
	s.FamiliesByName = make(map[string]*io_prometheus_client.MetricFamily)
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
