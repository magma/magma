package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var testUnmarshallMetricsCases = []struct {
	Json   string
	Result []Series
	Err    error
}{
	{
		Json:   "",
		Result: []Series{},
		Err:    nil,
	},
	{
		Json: "[{\"target\": \"main\", \"datapoints\": [[1.0, 1468339853], [1.0, 1468339854], [null, 1468339855]]}]",
		Result: []Series{
			{
				Target: "main",
				Datapoints: []DataPoint{
					{"1468339853", "1.0"},
					{"1468339854", "1.0"},
					{"1468339855", "null"},
				},
				Tags: make(map[string]string),
			},
		},
		Err: nil,
	},
}

func TestUnmarshallMetrics(t *testing.T) {
	for _, tc := range testUnmarshallMetricsCases {
		res, err := unmarshallSeries([]byte(tc.Json))
		if !reflect.DeepEqual(res, tc.Result) {
			t.Errorf("Result is %+v, \n expected %+v", res, tc.Result)
		}

		if err != tc.Err {
			t.Errorf("E %+v", err)
		}
	}
}

func TestNewClientFromString(t *testing.T) {
	urlString := "http://domain.tld/path"
	client, _ := NewFromString(urlString)
	testRequest := RenderRequest{}
	shouldUrl := urlString + testRequest.toQueryString()
	if shouldUrl != client.queryAsString(testRequest) {
		t.Errorf("Resulting URL is %v, \n but should be %v", client.queryAsString(testRequest), shouldUrl)
	}
}

func TestGraphiteRequest_ToQueryString(t *testing.T) {
	testCases := []struct {
		Request RenderRequest
		Result  string
	}{
		{
			Request: RenderRequest{},
			Result:  "/render/?format=json",
		},
		{
			Request: RenderRequest{Targets: []string{"foo", "bar"}},
			Result:  "/render/?format=json&target=foo&target=bar",
		},
		{
			Request: RenderRequest{From: time.Unix(1468339853, 0), Until: time.Unix(1468339853, 0)},
			Result:  "/render/?format=json&from=1468339853&until=1468339853",
		},
		{
			Request: RenderRequest{MaxDataPoints: 10},
			Result:  "/render/?format=json&maxDataPoints=10",
		},
	}
	for _, tc := range testCases {
		res := tc.Request.toQueryString()
		if res != tc.Result {
			t.Errorf("Result should be \"%v\", but \"%v\" received", tc.Result, res)
		}
	}
}

func makeTest(t *testing.T, request RenderRequest, expectedQuery, result string, series []Series) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedQuery, _ := url.ParseQuery(expectedQuery)
		if !reflect.DeepEqual(r.URL.Query(), parsedQuery) {
			t.Errorf("Expected query is %+v but %+v got", parsedQuery, r.URL.Query())
		}
		if r.URL.Path != "/render/" {
			t.Errorf("Path should be `/render/` but %s found", r.URL.Path)
		}
		fmt.Fprintln(w, result)
	}))
	defer ts.Close()
	client, _ := NewFromString(ts.URL)
	res, _ := client.QueryRender(request)
	if !reflect.DeepEqual(res, series) {
		t.Errorf("Expected series %+v, but %+v got", series, res)
	}
}

var queryTestCases = []struct {
	Request       RenderRequest
	ExpectedQuery string
	Result        string
	Series        []Series
}{
	{
		Request:       RenderRequest{Targets: []string{"main1", "main2"}},
		ExpectedQuery: "format=json&target=main1&target=main2",
		Result:        "[{\"target\": \"main\", \"datapoints\": [[1.0, 1468339853], [1.0, 1468339854], [null, 1468339855]]}]",
		Series: []Series{
			{
				Target: "main",
				Datapoints: []DataPoint{
					{"1468339853", "1.0"},
					{"1468339854", "1.0"},
					{"1468339855", "null"},
				},
				Tags: make(map[string]string),
			},
		},
	},
	{
		Request:       RenderRequest{From: time.Unix(1468339853, 0), Until: time.Unix(1468339854, 0)},
		ExpectedQuery: "format=json&from=1468339853&until=1468339854",
		Result:        "[{\"target\": \"main\", \"datapoints\": [[1.0, 1468339853], [1.0, 1468339854], [null, 1468339855]]}]",
		Series: []Series{
			{
				Target: "main",
				Datapoints: []DataPoint{
					{"1468339853", "1.0"},
					{"1468339854", "1.0"},
					{"1468339855", "null"},
				},
				Tags: make(map[string]string),
			},
		},
	},
	{
		Request:       RenderRequest{MaxDataPoints: 1},
		ExpectedQuery: "format=json&maxDataPoints=1",
		Result:        "[]",
		Series:        []Series{},
	},
}

func TestGraphiteClient_Query(t *testing.T) {
	for _, tc := range queryTestCases {
		makeTest(t, tc.Request, tc.ExpectedQuery, tc.Result, tc.Series)
	}
}
