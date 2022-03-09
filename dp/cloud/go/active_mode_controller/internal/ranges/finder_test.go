package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestFindAvailable(t *testing.T) {
	const minValue = -200
	testData := []struct {
		name     string
		points   []ranges.Point
		length   int
		value    int
		expected []ranges.Range
	}{{
		name: "Should get ranges spanning across two points",
		points: []ranges.Point{
			{Pos: 36000, Value: minValue},
			{Pos: 36500, Value: 20},
		},
		length: 100,
		value:  10,
		expected: []ranges.Range{
			{Begin: 36000, End: 36400, Value: 20},
		},
	}, {
		name: "Should get ranges spanning across multiple points",
		points: []ranges.Point{
			{Pos: 36000, Value: minValue},
			{Pos: 36100, Value: 20},
			{Pos: 36200, Value: 15},
		},
		length: 200,
		value:  10,
		expected: []ranges.Range{
			{Begin: 36000, End: 36000, Value: 15},
		},
	}, {
		name: "Should handle skip ranges with too low value",
		points: []ranges.Point{
			{Pos: 36000, Value: minValue},
			{Pos: 36100, Value: 20},
			{Pos: 36200, Value: 5},
			{Pos: 36300, Value: 20},
		},
		length: 50,
		value:  10,
		expected: []ranges.Range{
			{Begin: 36000, End: 36050, Value: 20},
			{Begin: 36200, End: 36250, Value: 20},
		},
	}, {
		name: "Should handle complex case",
		points: []ranges.Point{
			{Pos: 35500, Value: minValue},
			{Pos: 35600, Value: 30},
			{Pos: 35800, Value: 20},
			{Pos: 35900, Value: 25},
			{Pos: 36000, Value: 20},
			{Pos: 36200, Value: 10},
			{Pos: 36500, Value: 15},
			{Pos: 36600, Value: 10},
		},
		length: 200,
		value:  15,
		expected: []ranges.Range{
			{Begin: 35500, End: 35800, Value: 20},
			{Begin: 36200, End: 36300, Value: 15},
		},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := ranges.FindAvailable(tt.points, tt.length, tt.value)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
