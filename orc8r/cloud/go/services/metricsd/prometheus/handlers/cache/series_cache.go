package cache

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	estimatedBytesPerSeries = 100
)

// SeriesAPI interfaces the prometheus /series API to be used by the values
// handler for updating cache values
type SeriesAPI interface {
	// Series finds series by label matchers.
	Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, v1.Warnings, error)
}

type SeriesCache struct {
	responses     map[string]*cacheData
	specs         Specs
	backfillSpecs BackfillSpecs
	// updateFunc knows how to update the cache values so it can periodically
	// keep itself up to date
	updateFunc func(params []string, start, end time.Time) (values, error)
	sync.Mutex
}

type Params struct {
	Specs      Specs
	Backfill   BackfillSpecs
	UpdateFreq time.Duration
}

type Specs struct {
	// oldestAcceptable oldest age of a cache value to be returned
	OldestAcceptable time.Duration
	// ttl oldest age of individual values to keep if not seen recently
	TTL time.Duration
	// limitBytes maximum size of cache in memory. Delete old results if full.
	LimitBytes int
}

type BackfillSpecs struct {
	Lookback time.Duration
	Width    time.Duration
	Steps    int
}

func NewSeriesCache(params Params, updateFunc func(params []string, start, end time.Time) (values, error)) *SeriesCache {
	c := &SeriesCache{
		responses:     make(map[string]*cacheData),
		specs:         params.Specs,
		backfillSpecs: params.Backfill,
		updateFunc:    updateFunc,
	}
	if params.UpdateFreq != 0 {
		go c.updatePeriodically(params.UpdateFreq)
	}
	return c
}

// GetCacheUpdateProvider provides a function which will query the given series API
// for a specified set of params
func GetCacheUpdateProvider(api SeriesAPI) func(params []string, start, end time.Time) (values, error) {
	return func(params []string, start, end time.Time) (values, error) {
		res, _, err := api.Series(context.Background(), params, start, end)
		if err != nil {
			return values{}, err
		}
		return makeSeriesValuesNow(res), nil
	}
}

// Get returns a cached value if it exists and is new enough. If it exists but
// updated too long ago, it triggers an update
func (c *SeriesCache) Get(params []string) ([]model.LabelSet, bool) {
	key := paramsToKey(params)
	resp, exists := c.responses[key]
	if !exists {
		return []model.LabelSet{}, false
	}
	if time.Since(resp.updateTime) > c.specs.OldestAcceptable {
		go c.updateResponse(resp.params, resp.updateTime, time.Now())
		return []model.LabelSet{}, false
	}
	resp.requestTime = time.Now()
	return resp.getSeries(), true
}

// Set initializes a response with values and a function to update this response
func (c *SeriesCache) Set(params []string, series []model.LabelSet) {
	c.Lock()
	defer c.Unlock()
	// Delete oldest responses until this data fits
	for c.getEstimatedSize() > c.specs.LimitBytes {
		c.deleteOldestResponse()
	}
	vals := makeSeriesValuesNow(series)
	c.responses[paramsToKey(params)] = &cacheData{
		data:        vals,
		params:      params,
		requestTime: time.Now(),
		updateTime:  time.Now(),
	}
	go c.backfillValues(params)
}

// updatePeriodically iterates through each cache value and updates the value
// after "dur" has elapsed. This is used to keep cache values up to date to
// decrease cache misses
func (c *SeriesCache) updatePeriodically(dur time.Duration) {
	for range time.Tick(dur) {
		for _, resp := range c.responses {
			c.updateResponse(resp.params, time.Now().Add(-dur), time.Now())
		}
	}
}

// updateResponse takes merges old values with new values, updating their "lastSeen"
// time if they are already in this response
func (c *SeriesCache) updateResponse(params []string, start, end time.Time) {
	key := paramsToKey(params)
	newValues, err := c.updateFunc(params, start, end)
	if err != nil {
		return
	}
	c.Lock()
	defer c.Unlock()
	c.responses[key].updateTime = time.Now()
	c.responses[key].data = mergeData(c.responses[key].data, newValues, c.specs.TTL)
}

// backfillValues updates a cache value by iteratively looking back in time to
// get series. This is to avoid a single massive query which could overload
// the prometheus server and timeout.
func (c *SeriesCache) backfillValues(params []string) {
	if c.backfillSpecs.Lookback == 0 {
		return
	}
	now := time.Now()
	stepSize := c.backfillSpecs.Lookback / time.Duration(c.backfillSpecs.Steps)

	for step := 0; step < c.backfillSpecs.Steps; step++ {
		start := now.Add(-stepSize * time.Duration(step))
		end := start.Add(c.backfillSpecs.Width)
		c.updateResponse(params, start, end)
	}
}

// deleteOldestResponse deletes the least recently used response in this cache
// to make room for new responses
func (c *SeriesCache) deleteOldestResponse() {
	var oldest time.Duration = 0
	currentTime := time.Now()
	oldestKey := ""
	for key, resp := range c.responses {
		age := currentTime.Sub(resp.requestTime)
		if age > oldest {
			oldest = age
			oldestKey = key
		}
	}
	if oldestKey != "" {
		delete(c.responses, oldestKey)
	}
}

// getEstimatedSize calculates approximately the size in memory of a cache value.
// This is approximate since calculating a real size would be complex and costly
func (c *SeriesCache) getEstimatedSize() int {
	numSeries := 0
	for _, resp := range c.responses {
		numSeries += len(resp.data)
	}
	return estimatedBytesPerSeries * numSeries
}

// mergeData takes new values and adds them to the cache if they don't already
// exist in the cache, and removes values that are too old
// Note: uses labelset.Fingerprint() to get a key for each series. I've tested that
// function on a set of >100k series and it had 0 collisions so I think it will
// be fine here especially since accuracy is not super important.
func mergeData(oldv, newv values, ttl time.Duration) values {
	seriesSet := map[model.Fingerprint]struct{}{}
	mergedSeries := make([]value, 0, len(oldv))
	now := time.Now()
	for _, val := range oldv {
		age := now.Sub(val.lastSeen)
		if age > ttl {
			continue
		}
		seriesSet[val.series.Fingerprint()] = struct{}{}
		val.lastSeen = now
		mergedSeries = append(mergedSeries, val)
	}
	for _, val := range newv {
		if _, ok := seriesSet[val.series.Fingerprint()]; !ok {
			mergedSeries = append(mergedSeries, val)
		}
	}
	return mergedSeries
}

func paramsToKey(params []string) string {
	sort.Strings(params)
	return strings.Join(params, ",")
}

func makeSeriesValuesNow(series []model.LabelSet) values {
	now := time.Now()
	vals := make(values, len(series))
	for idx, ser := range series {
		vals[idx] = value{
			series:   ser,
			lastSeen: now,
		}
	}
	return vals
}

// cacheData is a struct that holds a list of labelsets (series) along
// with the last time each string has been seen, so it can remove strings that
// have not been seen in a long time.
type cacheData struct {
	data        values
	params      []string
	requestTime time.Time
	updateTime  time.Time
}

func (r *cacheData) getSeries() []model.LabelSet {
	series := make([]model.LabelSet, len(r.data))
	for idx, val := range r.data {
		series[idx] = val.series
	}
	return series
}

type values []value

type value struct {
	series   model.LabelSet
	lastSeen time.Time
}
