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

package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"magma/orc8r/cloud/go/services/analytics/protos"

	"net/http"
	"net/url"

	"magma/orc8r/lib/go/metrics"
)

type HttpClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type Exporter interface {
	Export(*protos.CalculationResult, HttpClient) error
}

type wwwExporter struct {
	metricsPrefix   string
	appSecret       string
	appID           string
	metricExportURL string
	categoryName    string
}

//NewWWWExporter exporter instance to export metrics
func NewWWWExporter(metricsPrefix, appID, appSecret, metricExportURL, categoryName string) Exporter {
	return &wwwExporter{
		metricsPrefix:   metricsPrefix,
		appID:           appID,
		appSecret:       appSecret,
		metricExportURL: metricExportURL,
		categoryName:    categoryName,
	}
}

func (e *wwwExporter) Export(res *protos.CalculationResult, client HttpClient) error {
	exportURL := fmt.Sprintf("%s?access_token=%s|%s", e.metricExportURL, e.appID, e.appSecret)

	nID := res.Labels[metrics.NetworkLabelName]
	sample := wwwDatapoint{
		Entity: e.FormatEntity(res, nID),
		Key:    e.FormatKey(res),
		Value:  fmt.Sprintf("%f", res.Value),
	}

	sampleJSON, err := json.Marshal([]wwwDatapoint{sample})
	if err != nil {
		return err
	}
	resp, err := client.PostForm(exportURL, url.Values{"datapoints": {string(sampleJSON)}, "category": {e.categoryName}})
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errMsg, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

func (e *wwwExporter) FormatKey(res *protos.CalculationResult) string {
	var keyBuffer bytes.Buffer
	keyBuffer.WriteString(res.MetricName)
	for labelName, labelValue := range res.Labels {
		if labelIsForbidden(labelName, forbiddenKeyLabelNames) {
			continue
		}
		keyBuffer.WriteString(".")
		keyBuffer.WriteString(labelName)
		keyBuffer.WriteString("-")
		keyBuffer.WriteString(labelValue)
	}
	return keyBuffer.String()
}

// Labels to not add to key
var forbiddenKeyLabelNames = []string{metrics.NetworkLabelName, metrics.ImsiLabelName}

func labelIsForbidden(labelName string, forbiddenLabels []string) bool {
	for _, forbidden := range forbiddenLabels {
		if labelName == forbidden {
			return true
		}
	}
	return false
}

func (e *wwwExporter) FormatEntity(res *protos.CalculationResult, nID string) string {
	return fmt.Sprintf("%s.analytics.%s", e.metricsPrefix, nID)
}

type wwwDatapoint struct {
	Entity string `json:"entity"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Time   int    `json:"time,omitempty"`
}
