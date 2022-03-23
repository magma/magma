package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestSelectPoint(t *testing.T) {
	rs := []ranges.Range{
		{Begin: 35499, End: 35601, Value: 4},
		{Begin: 35800, End: 36000, Value: 6},
		{Begin: 36200, End: 36250, Value: 8},
	}
	testData := []struct {
		name     string
		index    int
		multiply int
		expected ranges.Point
		found    bool
	}{{
		name:     "Should handle begin of non divisible boundaries",
		index:    0,
		multiply: 100,
		expected: ranges.Point{Value: 4, Pos: 35500},
		found:    true,
	}, {
		name:     "Should handle end of non divisible boundaries",
		index:    1,
		multiply: 100,
		expected: ranges.Point{Value: 4, Pos: 35600},
		found:    true,
	}, {
		name:     "Should handle begin of divisible boundaries",
		index:    2,
		multiply: 100,
		expected: ranges.Point{Value: 6, Pos: 35800},
		found:    true,
	}, {
		name:     "Should handle end of divisible boundaries",
		index:    4,
		multiply: 100,
		expected: ranges.Point{Value: 6, Pos: 36000},
		found:    true,
	}, {
		name:     "Should handle only begin visible",
		index:    5,
		multiply: 100,
		expected: ranges.Point{Value: 8, Pos: 36200},
		found:    true,
	}, {
		name:     "Should handle overlap",
		index:    64,
		multiply: 100,
		expected: ranges.Point{Value: 6, Pos: 36000},
		found:    true,
	}, {
		name:     "Should handle not found case",
		index:    0,
		multiply: 10000,
		found:    false,
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual, found := ranges.Select(rs, tt.index, tt.multiply)
			assert.Equal(t, tt.expected, actual)
			assert.Equal(t, tt.found, found)
		})
	}
}
