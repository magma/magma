package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestSelectPoint(t *testing.T) {
	rs := []ranges.Range{
		{Begin: 35500, End: 35600, Value: 4},
		{Begin: 35700, End: 35900, Value: 6},
		{Begin: 36000, End: 36000, Value: 8},
	}
	testData := []struct {
		name     string
		index    int
		expected ranges.Point
	}{{
		name:     "Should pick first point in range",
		index:    101,
		expected: ranges.Point{Value: 6, Pos: 35700},
	}, {
		name:     "Should pick last point in range",
		index:    100,
		expected: ranges.Point{Value: 4, Pos: 35600},
	}, {
		name:     "Should pick middle point in range",
		index:    201,
		expected: ranges.Point{Value: 6, Pos: 35800},
	}, {
		name:     "Should pick point from one point range",
		index:    302,
		expected: ranges.Point{Value: 8, Pos: 36000},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			r := &stubIndexProvider{index: tt.index}
			actual := ranges.SelectPoint(rs, r)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

type stubIndexProvider struct {
	index int
}

func (s *stubIndexProvider) Intn(n int) int {
	return s.index % n
}
