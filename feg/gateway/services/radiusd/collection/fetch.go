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

package collection

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"magma/feg/cloud/go/protos/mconfig"
	gwmcfg "magma/gateway/mconfig"

	"github.com/golang/glog"
)

const (
	defaultMetricsServerHost     = "radius"
	defaultMetricsServerHostPort = 9100
	defaultMetricsPath           = "metrics"
	defaultUpdateIntervalSecs    = 30
)

type MetricsRequester struct {
	metricsUrl string
}

func NewMetricsRequester() (*MetricsRequester, error) {
	metricsUrl := getMetricsUrl()
	return &MetricsRequester{
		metricsUrl: metricsUrl,
	}, nil
}

// FetchMetrics makes a request to the radius metrics server.
// The GET response body is returned, and this does not process the prometheus
// metrics info in any way.
func (r *MetricsRequester) FetchMetrics() (string, error) {
	resp, err := http.Get(r.metricsUrl)
	if err != nil {
		return "", fmt.Errorf("Failed GET request: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read GET response body: %s", err)
	}

	return string(body), nil
}

// RefreshConfig tries to refresh configs
func (r *MetricsRequester) RefreshConfig() {
	r.metricsUrl = getMetricsUrl()
}

func getMetricsUrl() string {
	radiusdCfg := GetRadiusdConfig()
	host := radiusdCfg.GetRadiusMetricsHost()
	port := radiusdCfg.GetRadiusMetricsPort()
	path := radiusdCfg.GetRadiusMetricsPath()
	return fmt.Sprintf("http://%s:%d/%s", host, port, path)
}

// GetRadiusdConfig attempts to retrieve a RadiusdConfig  from mconfig
// If this retrieval fails, or retrieves an invalid config, the config is
// set to use default values
func GetRadiusdConfig() *mconfig.RadiusdConfig {
	radiusdCfg := &mconfig.RadiusdConfig{}
	err := gwmcfg.GetServiceConfigs("radiusd", radiusdCfg)
	if err != nil {
		glog.Infof("Unable to retrieve Radiusd Config from mconfig: %s; Using default values...", err)
		return &mconfig.RadiusdConfig{
			RadiusMetricsPort:  defaultMetricsServerHostPort,
			RadiusMetricsPath:  defaultMetricsPath,
			UpdateIntervalSecs: defaultUpdateIntervalSecs,
			RadiusMetricsHost:  defaultMetricsServerHost,
		}
	}
	glog.Info("Using mconfig values for radiusd parameters")
	return radiusdCfg
}
