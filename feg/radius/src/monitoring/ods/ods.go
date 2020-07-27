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

package ods

import (
	"crypto/tls"
	"encoding/json"
	"fbc/cwf/radius/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

var (
	// DefaultTransport is the default tracing transport and is used by DefaultClient.
	DefaultTransport http.RoundTripper = &ochttp.Transport{}

	// DefaultClient is the default tracing http client.
	DefaultClient = &http.Client{
		Transport: DefaultTransport,
	}
)

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

// PostToODS goes through all the timeseries in from a request,
// coverts the labels to keys/entities and posts to ODS via GraphAPI
func PostToODS(metricsData map[string]string, cfg config.Ods) error {
	var datapoints []Datapoint
	var entity string
	ts := time.Now().Unix()
	hostname, _ := os.Hostname()
	entity = fmt.Sprintf("%s.%s", cfg.Entity, hostname)
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
	req, err := http.NewRequest(
		http.MethodPost,
		cfg.GraphURL,
		strings.NewReader(urlValues.Encode()),
	)
	if err != nil {
		return errors.WithMessage(err, "failed to create http post request")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := DefaultClient.Do(req)
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
	if cfg.DebugPrints {
		for _, dp := range datapoints {
			log.Printf("%v\n", dp)
		}
	}
	return nil
}

func getURLValues(datapoints []Datapoint, cfg config.Ods) (url.Values, error) {
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
	config config.Ods
}

// TODO: (@ayeletrd T44633984) Do not emit errors here, add callback.
func (ce *odsMetricsExporter) ExportView(vd *view.Data) {
	metricsData := make(map[string]string)
	if len(vd.Rows) == 0 {
		return
	}
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
			metricsData[key+".count"] = strconv.FormatInt(count.Value, 10)
		case "Sum":
			count, _ := val.(*view.SumData)
			metricsData[key+".sum"] = strconv.FormatFloat(count.Value, 'f', 6, 64)
		case "LastValue":
			count, _ := val.(*view.LastValueData)
			metricsData[key+".gauge"] = strconv.FormatFloat(count.Value, 'f', 6, 64)
		case "Distribution":
			dist, _ := val.(*view.DistributionData)
			metricsData[key+".count"] = strconv.FormatInt(dist.Count, 10)
			metricsData[key+".sum"] = fmt.Sprintf("%f", float64(dist.Count)*dist.Mean)
			metricsData[key+".avg"] = fmt.Sprintf("%f", dist.Mean)
			metricsData[key+".min"] = fmt.Sprintf("%f", dist.Min)
			metricsData[key+".max"] = fmt.Sprintf("%f", dist.Max)
		}
		key = ""
	}
	if ce.config.DisablePost == true {
		return
	}

	go func() {
		err := PostToODS(metricsData, ce.config)
		if err != nil {
			log.Printf("failed to send post request of %v (%v) to %v. err: %v",
				vd.View.Name, metricsData, ce.config.GraphURL, err.Error())
		}
	}()
}

// Init Should be called once if ODS counters are to be emmitted
func Init(odsConfig *config.Ods, logger *zap.Logger) {
	// If no ODS configuration is there - skip initialization
	if odsConfig == nil {
		logger.Info("no ODS configuration, skipping initialization")
		return
	}
	config, _ := json.Marshal(odsConfig)
	logger.Info("initializing ODS counters", zap.String("config", string(config)))

	log.SetFlags(0)
	view.RegisterExporter(&odsMetricsExporter{
		config: *odsConfig,
	})
	view.SetReportingPeriod(odsConfig.ReportingPeriod.Duration)
}
