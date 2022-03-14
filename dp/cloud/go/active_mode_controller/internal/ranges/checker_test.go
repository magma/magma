package ranges_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/ranges"
)

func TestCheckIfContain(t *testing.T) {
	rs := []ranges.Range{
		{Begin: 35700, End: 35800},
		{Begin: 36000, End: 36000},
		{Begin: 36200, End: 36300},
	}
	points := []int{36500, 36250, 35500, 36000, 35850}
	expected := []bool{false, true, false, true, false}
	actual := ranges.CheckIfContain(rs, points)
	assert.Equal(t, expected, actual)
}
