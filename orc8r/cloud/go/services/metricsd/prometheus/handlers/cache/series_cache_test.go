package cache

import (
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/handlers/mocks"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	sampleLabelSet = []model.LabelSet{{"__name__": "val1"}}
)

func TestValuesCache_Set(t *testing.T) {
	mockAPI := &mocks.SeriesAPI{}
	testCache := getTestCache(time.Minute, time.Minute, 150, mockAPI, 0)

	// Basic Cache Set works
	testCache.Set([]string{"testParam1"}, sampleLabelSet)
	assert.Len(t, testCache.responses, 1)
	time.Sleep(time.Millisecond * 500)
	testCache.Set([]string{"testParam1"}, sampleLabelSet)
	assert.Len(t, testCache.responses, 1)
	time.Sleep(time.Millisecond * 500)
	testCache.Set([]string{"testParam2"}, sampleLabelSet)
	assert.Len(t, testCache.responses, 2)
	time.Sleep(time.Millisecond * 500)

	// Cache deletes oldest series when full
	testCache.Set([]string{"testParam3"}, sampleLabelSet)
	assert.Len(t, testCache.responses, 2)
	assert.NotNil(t, testCache.responses["testParam3"])
	assert.NotNil(t, testCache.responses["testParam2"])
	assert.Nil(t, testCache.responses["testParam1"])
}

func TestValuesCache_Get(t *testing.T) {
	mockAPI := &mocks.SeriesAPI{}
	testCache := getTestCache(time.Minute, time.Minute, 10000, mockAPI, 0)

	// Basic Cache Get works
	_, ok := testCache.Get([]string{"testParam1"})
	assert.False(t, ok)
	testCache.Set([]string{"testParam1"}, sampleLabelSet)
	data, ok := testCache.Get([]string{"testParam1"})
	assert.True(t, ok)
	assert.Equal(t, sampleLabelSet, data)

	// No cache hit if result out of date
	mockAPI.On("Series", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sampleLabelSet, nil, nil)
	testCache.Set([]string{"testParam2"}, sampleLabelSet)
	testCache.responses["testParam2"].updateTime = time.Now().Add(-5 * time.Minute)
	_, ok = testCache.Get([]string{"testParam2"})
	assert.False(t, ok)
	// Wait for goroutine updateFunc to be called
	time.Sleep(time.Second)
	mockAPI.AssertNumberOfCalls(t, "Series", 1)
	// Now that it's been updated, cache get should return
	_, ok = testCache.Get([]string{"testParam2"})
	assert.True(t, ok)
}

func getTestCache(oldestAcceptable, ttl time.Duration, limit int, mockAPI SeriesAPI, updateFreq time.Duration) *SeriesCache {
	return NewSeriesCache(Params{
		Specs: Specs{
			OldestAcceptable: oldestAcceptable,
			TTL:              ttl,
			LimitBytes:       limit,
		},
		UpdateFreq: updateFreq,
	}, GetCacheUpdateProvider(mockAPI))
}
