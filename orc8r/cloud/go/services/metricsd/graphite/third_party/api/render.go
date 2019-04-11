package api

import (
	"net/url"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
)

// RenderRequest is struct, describing request to graphite `/render/` api.
// No fields are required. If field has zero value it'll be just skipped in request.
// RenderRequest.Targets are slice of strings, were every entry is a path identifying one or several metrics,
// optionally with functions acting on those metrics.
//
// Warning. While wildcards could be used in Targets one should use them with caution, as
// using of the simple target like "main.cluster.*.cpu.*" could result in hundreds of series
// with megabytes of data inside.
type RenderRequest struct {
	From          time.Time
	Until         time.Time
	MaxDataPoints int
	Targets       []string
}

func (r RenderRequest) toQueryString() string {
	values := url.Values{
		"format": []string{"json"},
		"target": r.Targets,
	}
	if !r.From.IsZero() {
		values.Set("from", strconv.FormatInt(r.From.Unix(), 10))
	}
	if !r.Until.IsZero() {
		values.Set("until", strconv.FormatInt(r.Until.Unix(), 10))
	}
	if r.MaxDataPoints != 0 {
		values.Set("maxDataPoints", strconv.Itoa(r.MaxDataPoints))
	}
	qs := values.Encode()
	return "/render/?" + qs
}

// QueryRender performs query to graphite `/render/` api. Normally it should return `[]graphite.Series`,
// but if things go wrong it will return `graphite.RequestError` error.
func (c *Client) QueryRender(r RenderRequest) ([]Series, error) {
	empty := []Series{}
	data, err := c.makeRequest(r)
	if err != nil {
		return empty, err
	}

	metrics, err := unmarshallSeries(data)
	if err != nil {
		return empty, c.createError(r, "Can't unmarshall response")
	}
	return metrics, nil
}

// Series describes time series data for given target.
type Series struct {
	Target     string            `json:"target"`
	Datapoints []DataPoint       `json:"datapoints"`
	Tags       map[string]string `json:"tags"`
}

// DataPoint describes concrete point of time series.
type DataPoint []string

func unmarshallSeries(data []byte) ([]Series, error) {
	empty, result := []Series{}, []Series{}
	if len(data) == 0 {
		return empty, nil
	}
	var ie error = nil
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}

		datapoints, e := unmarshallDatapoints(value)
		if e != nil {
			ie = e
			return
		}

		target, e := jsonparser.GetString(value, "target")
		if e != nil {
			ie = e
			return
		}

		tags := unmarshallTags(value)

		result = append(result, Series{Target: target, Datapoints: datapoints, Tags: tags})
	})

	if err != nil {
		return empty, err
	}
	if ie != nil {
		return empty, ie
	}
	return result, nil
}

func unmarshallDatapoints(data []byte) ([]DataPoint, error) {
	empty, result := []DataPoint{}, []DataPoint{}
	rawData, _, _, err := jsonparser.Get(data, "datapoints")
	if err != nil {
		return empty, err
	}

	_, err = jsonparser.ArrayEach(rawData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		datapoint, e := unmarshallDatapoint(value)
		if e != nil {
			err = e
			return
		}
		result = append(result, datapoint)
	})
	if err != nil {
		return empty, err
	}
	return result, nil
}

func unmarshallTags(data []byte) map[string]string {
	tags := make(map[string]string)
	rawData, _, _, err := jsonparser.Get(data, "tags")
	if err != nil {
		return tags
	}

	err = jsonparser.ObjectEach(rawData, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
		if err != nil {
			return err
		}
		tags[string(key)] = string(value)
		return nil
	})
	return tags
}

func unmarshallDatapoint(data []byte) (DataPoint, error) {
	empty, result := DataPoint{}, make(DataPoint, 2)
	var err error = nil
	position := 0
	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}
		if position == 0 {
			if dataType == jsonparser.Null {
				result[1] = "null"
			} else {
				result[1] = string(value)
			}
		} else {
			result[0] = string(value)
		}
		position++
	})
	if err != nil {
		return empty, err
	}
	return result, nil
}
