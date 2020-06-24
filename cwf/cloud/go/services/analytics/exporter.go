package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"magma/orc8r/lib/go/metrics"
)

type HttpClient interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type Exporter interface {
	Export(Result, HttpClient) error
}

type wwwExporter struct {
	metricsPrefix   string
	appSecret       string
	appID           string
	metricExportURL string
	categoryName    string
}

func NewWWWExporter(metricsPrefix, appSecret, appID, metricExportURL, categoryName string) Exporter {
	return &wwwExporter{
		metricsPrefix:   metricsPrefix,
		appSecret:       appSecret,
		appID:           appID,
		metricExportURL: metricExportURL,
		categoryName:    categoryName,
	}
}

func (e *wwwExporter) Export(res Result, client HttpClient) error {
	exportURL := fmt.Sprintf("%s?access_token=%s|%s", e.metricExportURL, e.appID, e.appSecret)

	nID := res.labels[metrics.NetworkLabelName]
	if nID == "" {
		return fmt.Errorf("no networkID for exported metric")
	}
	sample := wwwDatapoint{
		Entity: e.FormatEntity(res, nID),
		Key:    e.FormatKey(res),
		Value:  fmt.Sprintf("%f", res.value),
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

func (e *wwwExporter) FormatKey(res Result) string {
	var keyBuffer bytes.Buffer
	keyBuffer.WriteString(res.metricName)
	for labelName, labelValue := range res.labels {
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
var forbiddenKeyLabelNames = []string{metrics.NetworkLabelName}

func labelIsForbidden(labelName string, forbiddenLabels []string) bool {
	for _, forbidden := range forbiddenLabels {
		if labelName == forbidden {
			return true
		}
	}
	return false
}

func (e *wwwExporter) FormatEntity(res Result, nID string) string {
	return fmt.Sprintf("%s.analytics.%s", e.metricsPrefix, nID)
}

type wwwDatapoint struct {
	Entity string `json:"entity"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Time   int    `json:"time,omitempty"`
}
