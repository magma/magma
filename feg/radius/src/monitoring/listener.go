package monitoring

import (
	"context"
	"fbc/lib/go/radius"
	"fmt"
	"sync"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type (
	// RequestCounter ...
	RequestCounter interface {
		Failure(errorCode string)
		GotResponse(radiusCode radius.Code)
	}

	requestCounter struct {
		listenerTag tag.Mutator
		request     Operation
		response    *stats.Int64Measure
	}

	// ListenerCounters ...
	ListenerCounters interface {
		StartRequest(radiusCode radius.Code) RequestCounter
	}

	listenerCounters struct {
		listenerName string
		opCounters   sync.Map
	}
)

func (c *requestCounter) Failure(errorCode string) {
	c.request.Failure(errorCode)
}

func (c *requestCounter) GotResponse(radiusCode radius.Code) {
	c.request.Success()
	stats.RecordWithTags(
		context.Background(),
		[]tag.Mutator{
			tag.Upsert(ResponseCodeTag, radiusCode.String()),
		},
		c.response.M(1),
	)
}

func (c *listenerCounters) StartRequest(radiusCode radius.Code) RequestCounter {
	listenerTag := tag.Upsert(ListenerTag, c.listenerName)

	// Get counter
	requestOpCounter, _ := c.opCounters.LoadOrStore(radiusCode, NewOperation(
		"handle_request",
		listenerTag,
		tag.Upsert(RequestCodeTag, radiusCode.String()),
	))

	//
	result := &requestCounter{
		listenerTag: listenerTag,
		request:     (requestOpCounter.(Operation)).Start(),
		response: stats.Int64(
			"response",
			fmt.Sprintf("Operation '%s' started", c.listenerName),
			stats.UnitDimensionless,
		),
	}

	// expose the response counter
	view.Register(&view.View{
		Name:        "response",
		Measure:     result.response,
		Description: fmt.Sprintf("The number of time '%s' was started", c.listenerName),
		Aggregation: view.Count(),
		TagKeys:     AllTagKeys(),
	})

	return result
}

// CreateListenerCounters ...
func CreateListenerCounters(name string) ListenerCounters {
	result := &listenerCounters{
		listenerName: name,
		opCounters:   sync.Map{},
	}
	return result
}
