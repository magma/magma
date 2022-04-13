package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestSplit(t *testing.T) {
	testData := []struct {
		name           string
		ranges         []ranges.Range
		points         []int
		expectedRanges []ranges.Range
		expectedPoints []ranges.Point
	}{{
		name: "Should split range in the middle",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		points: []int{35650},
		expectedRanges: []ranges.Range{
			{Begin: 35600, End: 35649, Value: 30},
			{Begin: 35651, End: 35700, Value: 30},
		},
		expectedPoints: []ranges.Point{
			{Pos: 35650, Value: 30},
		},
	}, {
		name: "Should split range in the ends",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		points: []int{35600, 35700},
		expectedRanges: []ranges.Range{
			{Begin: 35601, End: 35699, Value: 30},
		},
		expectedPoints: []ranges.Point{
			{Pos: 35600, Value: 30},
			{Pos: 35700, Value: 30},
		},
	}, {
		name: "Should split range multiple times",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		points: []int{35601, 35620, 35650, 35699, 35700},
		expectedRanges: []ranges.Range{
			{Begin: 35600, End: 35600, Value: 30},
			{Begin: 35602, End: 35619, Value: 30},
			{Begin: 35621, End: 35649, Value: 30},
			{Begin: 35651, End: 35698, Value: 30},
		},
		expectedPoints: []ranges.Point{
			{Pos: 35601, Value: 30},
			{Pos: 35620, Value: 30},
			{Pos: 35650, Value: 30},
			{Pos: 35699, Value: 30},
			{Pos: 35700, Value: 30},
		},
	}, {
		name: "Should take out all elements of range",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35600, Value: 30},
		},
		points:         []int{35600},
		expectedRanges: nil,
		expectedPoints: []ranges.Point{
			{Pos: 35600, Value: 30},
		},
	}, {
		name: "Should not split if point not in range",
		ranges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		points: []int{35800},
		expectedRanges: []ranges.Range{
			{Begin: 35600, End: 35700, Value: 30},
		},
		expectedPoints: nil,
	}, {
		name: "Should handle complex case",
		ranges: []ranges.Range{
			{Begin: 35700, End: 35800, Value: 30},
			{Begin: 36000, End: 36000, Value: 20},
			{Begin: 36200, End: 36300, Value: 10},
		},
		points: []int{35500, 35850, 36000, 36250, 36300, 36500},
		expectedRanges: []ranges.Range{
			{Begin: 35700, End: 35800, Value: 30},
			{Begin: 36200, End: 36249, Value: 10},
			{Begin: 36251, End: 36299, Value: 10},
		},
		expectedPoints: []ranges.Point{
			{Pos: 36000, Value: 20},
			{Pos: 36250, Value: 10},
			{Pos: 36300, Value: 10},
		},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actualRanges, actualPoints := ranges.Split(tt.ranges, tt.points)
			assert.Equal(t, tt.expectedRanges, actualRanges)
			assert.Equal(t, tt.expectedPoints, actualPoints)
		})
	}
}
