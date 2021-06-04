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

package calculations

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

type CertLifespanCalculation struct {
	CertsDirectory string
	Certs          []string
	calculations.BaseCalculation
}

func (x *CertLifespanCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(2).Infof("Calculating %s", metrics.CertExpiresInHoursMetric)
	var results []*protos.CalculationResult

	metricConfig, ok := x.AnalyticsConfig.Metrics[metrics.CertExpiresInHoursMetric]
	if !ok {
		glog.Errorf("%s metric not found in metric config", metrics.CertExpiresInHoursMetric)
		return results, nil
	}

	for _, certName := range x.Certs {
		result, err := calculateCertLifespanHours(x.CertsDirectory, certName, metricConfig.Labels)
		if err != nil {
			glog.Errorf("Could not get lifespan for cert %s: %+v", certName, err)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func calculateCertLifespanHours(certsDirectory string, certName string, metricConfigLabels map[string]string) (*protos.CalculationResult, error) {
	dat, err := getCert(certsDirectory + certName)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(dat)
	if err != nil {
		return nil, err
	}
	// Hours remaining
	hoursLeft := math.Floor(time.Until(cert.NotAfter).Hours())
	labels := prometheus.Labels{
		metrics.CertNameLabel: certName,
	}
	labels = calculations.CombineLabels(labels, metricConfigLabels)
	result := calculations.NewResult(hoursLeft, metrics.CertExpiresInHoursMetric, labels)
	glog.V(2).Infof("Calculated metric %s for %s: %f", metrics.CertExpiresInHoursMetric, certName, hoursLeft)
	return result, nil
}

func getCert(certPath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(certPath, ".pem") || strings.HasSuffix(certPath, ".crt") {
		block, _ := pem.Decode(dat)
		if block == nil || block.Type != "PUBLIC KEY" {
			log.Fatal("failed to decode PEM block containing public key")
		}
		return block.Bytes, nil
	}
	return dat, nil
}
