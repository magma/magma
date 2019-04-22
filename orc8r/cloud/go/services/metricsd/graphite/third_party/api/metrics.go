package api

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// FindMetricRequest is struct describing request to /metrics/find api.
type FindMetricRequest struct {
	From      time.Time
	Until     time.Time
	Query     string
	Wildcards bool
}

func (r FindMetricRequest) toQueryString() string {
	values := url.Values{
		"format": []string{"treejson"},
	}
	if !r.From.IsZero() {
		values.Set("from", strconv.FormatInt(r.From.Unix(), 10))
	}
	if !r.Until.IsZero() {
		values.Set("until", strconv.FormatInt(r.Until.Unix(), 10))
	}
	if r.Query != "" {
		values.Set("query", r.Query)
	}
	if r.Wildcards {
		values.Set("wildcards", strconv.Itoa(1))
	}
	qs := values.Encode()
	return "/metrics/find?" + qs
}

type Metric struct {
	Id            string
	Text          string
	Expandable    int
	Leaf          int
	AllowChildren int
}

// FindMetrics perform request to /metrics/find API: http://graphite-api.readthedocs.io/en/latest/api.html#metrics-find
// It returns slice of Metric if all is OK or RequestError if things goes wrong.
func (c *Client) FindMetrics(r FindMetricRequest) ([]Metric, error) {
	empty := []Metric{}
	data, err := c.makeRequest(r)
	if err != nil {
		return empty, err
	}

	metrics, err := unmarshallMetrics(data)
	if err != nil {
		return empty, c.createError(r, "Can't unmarshall response")
	}
	return metrics, nil
}

func unmarshallMetrics(data []byte) ([]Metric, error) {
	var metrics []Metric
	err := json.Unmarshal(data, &metrics)
	return metrics, err
}

type ExpandMetricRequest struct {
	Query       string
	GroupByExpr bool
	LeavesOnly  bool
}

func (r ExpandMetricRequest) toQueryString() string {
	values := url.Values{}
	if r.Query != "" {
		values.Set("query", r.Query)
	}
	if r.GroupByExpr {
		values.Set("groupByExpr", strconv.Itoa(1))
	}
	if r.LeavesOnly {
		values.Set("leavesOnly", strconv.Itoa(1))
	}
	qs := values.Encode()
	return "/metrics/expand?" + qs
}

// Results is a list of metric ids.
type ExpandResult struct {
	Results []string
}

// FindMetrics perform request to /metrics/expand API: http://graphite-api.readthedocs.io/en/latest/api.html#metrics-expand
// It returns slice of Metric if all is OK or RequestError if things goes wrong.
func (c *Client) ExpandMetrics(r ExpandMetricRequest) (ExpandResult, error) {
	result := ExpandResult{}
	data, err := c.makeRequest(r)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, c.createError(r, "Cant unmarshall response")
	}
	return result, nil
}
