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
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/pkg/errors"
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

	networks, err := configurator.ListNetworkIDs()
	if err != nil {
		glog.Errorf("Unable to retrieve any networks: %+v", err)
		return results, nil
	}

	for _, certName := range x.Certs {
		cert_results, err := calculateCertLifespanHours(x.CertsDirectory, certName, metricConfig.Labels, networks)
		if err != nil {
			glog.Errorf("Could not get lifespan for cert %s: %+v", certName, err)
			continue
		}
		results = append(results, cert_results...)
	}
	return results, nil
}

func calculateCertLifespanHours(certsDirectory string, certName string, metricConfigLabels map[string]string, networks []string) ([]*protos.CalculationResult, error) {
	var results []*protos.CalculationResult
	dat, err := getCert(certsDirectory + certName)
	if err != nil {
		return nil, err
	}
	cert, err := x509.ParseCertificate(dat)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse x509 certificate data for %s", certName))
	}
	// Hours remaining
	hoursLeft := math.Floor(time.Until(cert.NotAfter).Hours())
	glog.V(2).Infof("Calculated metric %s for %s: %f", metrics.CertExpiresInHoursMetric, certName, hoursLeft)

	// Here we are broadcasting infra level certificate alert on all networks
	// This is not a typical pattern, however we are currently doing this to
	// enable the certificate expiry alert to be displayed on per tenant NMS portal
	for _, networkID := range networks {
		labels := prometheus.Labels{
			metrics.NetworkLabelName: networkID,
			metrics.CertNameLabel:    certName,
		}
		labels = calculations.CombineLabels(labels, metricConfigLabels)
		result := calculations.NewResult(hoursLeft, metrics.CertExpiresInHoursMetric, labels)
		results = append(results, result)
	}
	return results, nil
}

func getCert(certPath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read cert file %s", certPath))
	}
	if strings.HasSuffix(certPath, ".pem") || strings.HasSuffix(certPath, ".crt") {
		block, _ := pem.Decode(dat)
		if block == nil {
			return nil, fmt.Errorf("failed to decode a PEM block containing public key for certificate %s", certPath)
		} else if block.Type != "PUBLIC KEY" && block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("certificate %s has a PEM block type that is not PUBLIC KEY or CERTIFICATE, and instead is %s", certPath, block.Type)
		}
		return block.Bytes, nil
	}
	return dat, nil
}
