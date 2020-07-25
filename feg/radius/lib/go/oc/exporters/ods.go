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

package exporters

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"fbc/lib/go/oc"
	"github.com/kelseyhightower/envconfig"
	"go.opencensus.io/stats/view"
)

// ConfigProvider is used to mock configs for unitests
type ConfigProvider interface {
	getConfig() Config
	setConfig(Config) Config
}

type prodConfigProvider struct {
	cfg Config
}

// Config needed in order to report to ODs should be also be defined in
// your docker-compose
type Config struct {
	Category        string        `envconfig:"ODS_CATEGORY_ID" required:"true"`
	DisablePost     bool          `envconfig:"DISABLE_POST" default:"false"`
	GraphURL        string        `envconfig:"GRAPH_URL" default:"https://graph.facebook.com/ods_metrics"`
	Prefix          string        `envconfig:"ODS_PREFIX" required:"true"`
	Token           string        `envconfig:"ODS_ACCESS_TOKEN" required:"true"`
	Entity          string        `envconfig:"ODS_ENTITY" required:"true"`
	ReportingPeriod time.Duration `envconfig:"ODS_REPORTING_PERIOD" default:"60s"`
	UniqueEntity    bool          `envconfig:"ODS_UNIQUE_ENTITY" default:"true"`
}

// Datapoint is used to Marshal JSON encoding for ODS data submission
// see https://phabricator.intern.facebook.com/diffusion/E/browse/tfb/trunk/www/flib/platform/graph/resources/ods/metrics/GraphOdsMetricsPost.php
// for types accepted
type Datapoint struct {
	Entity string   `json:"entity"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Time   int64    `json:"time"`
	Tags   []string `json:"tags"`
}

func (pscp *prodConfigProvider) setConfig(cfg Config) Config {
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("processing config: %v", err)
	}
	pscp.cfg = cfg
	return pscp.cfg
}

func (pscp *prodConfigProvider) getConfig() Config {
	return pscp.cfg
}

// PostToODS goes through all the timeseries in from a request,
// coverts the labels to keys/entities and posts to ODS via GraphAPI
func PostToODS(metricsData map[string]string, cfg Config) error {
	var datapoints []Datapoint
	var entity string
	ts := time.Now().Unix()
	if cfg.UniqueEntity {
		hostname, _ := os.Hostname()
		entity = fmt.Sprintf("%s.%s.%s", cfg.Prefix, cfg.Entity, hostname)
	} else {
		// Mostly used for unitests but could also disable the feature.
		entity = fmt.Sprintf("%s.%s", cfg.Prefix, cfg.Entity)
	}
	for k, v := range metricsData {
		datapoints = append(datapoints, Datapoint{
			Entity: entity,
			Key:    k,
			Value:  v,
			Time:   ts,
		})
	}

	if len(datapoints) == 0 {
		return fmt.Errorf(
			"empty datapoints: %v", "no valid datapoints found not posting to ODS")
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	urlValues, err := getURLValues(datapoints, cfg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, cfg.GraphURL, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return errors.WithMessage(err, "failed to create http post request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := oc.DefaultClient.Do(req)

	if err != nil {
		return errors.WithMessage(err, "failed to post form")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		errMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to post to ODS and failed to get reason why we failed")
		}
		return fmt.Errorf("failed to post to ODS: %v", string(errMsg))
	}
	log.Printf("\nsubmitting these datapoints %v\n", datapoints)
	return nil
}

func getURLValues(datapoints []Datapoint, cfg Config) (url.Values, error) {
	urlValues := url.Values{}
	datapointsJSON, err := json.Marshal(datapoints)
	if err != nil {
		return urlValues, errors.WithMessage(err, "error marshaling datapoints")
	}

	urlValues.Add("access_token", cfg.Token)
	urlValues.Add("category_id", cfg.Category)
	urlValues.Add("datapoints", string(datapointsJSON))
	return urlValues, nil
}

type odsMetricsExporter struct {
	metricsData    map[string]string
	configProvider ConfigProvider
}

// TODO: (@ayeletrd T44633984) Do not emit errors here, add callback.
func (ce *odsMetricsExporter) ExportView(vd *view.Data) {
	ce.metricsData = make(map[string]string)
	if len(vd.Rows) == 0 {
		return
	}
	if ce.configProvider == nil {
		// Production will always use prodConfigProvider, which will be mocked in tests
		ce.configProvider = new(prodConfigProvider)
		ce.configProvider.setConfig(Config{})
	}
	cfg := ce.configProvider.getConfig()
	var key string
	for _, row := range vd.Rows {
		for _, tag := range row.Tags {
			key += fmt.Sprintf("%v.", tag.Value)
		}
		val := row.Data
		key += vd.View.Name
		switch vd.View.Aggregation.Type.String() {
		case "Count":
			count, _ := val.(*view.CountData)
			ce.metricsData[key+".count"] = strconv.FormatInt(count.Value, 10)
		case "Sum":
			count, _ := val.(*view.SumData)
			ce.metricsData[key+".sum"] = strconv.FormatFloat(count.Value, 'f', 6, 64)
		case "LastValue":
			count, _ := val.(*view.LastValueData)
			ce.metricsData[key+".gauge"] = strconv.FormatFloat(count.Value, 'f', 6, 64)
		case "Distribution":
			dist, _ := val.(*view.DistributionData)
			ce.metricsData[key+".count"] = strconv.FormatInt(dist.Count, 10)
			ce.metricsData[key+".sum"] = fmt.Sprintf("%f", float64(dist.Count)*dist.Mean)
			ce.metricsData[key+".avg"] = fmt.Sprintf("%f", dist.Mean)
			ce.metricsData[key+".min"] = fmt.Sprintf("%f", dist.Min)
			ce.metricsData[key+".max"] = fmt.Sprintf("%f", dist.Max)
		}
		key = ""
	}
	if cfg.DisablePost == true {
		return
	}
	err := PostToODS(ce.metricsData, cfg)
	if err != nil {
		log.Printf("failed to send post request of %v (%v) to %v. err: %v",
			vd.View.Name, ce.metricsData, cfg.GraphURL, err.Error())
	}
}

// StartODSExporter function should be called from the app main function after
// registering the project views.
func StartODSExporter() context.Context {
	log.SetFlags(0)
	ce := odsMetricsExporter{}
	ce.configProvider = new(prodConfigProvider)
	cfg := ce.configProvider.setConfig(Config{})
	ctx := context.Background()
	odsExporter := new(odsMetricsExporter)
	view.RegisterExporter(odsExporter)
	view.SetReportingPeriod(cfg.ReportingPeriod)
	return ctx
}
