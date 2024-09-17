package plmn_filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	plmns_5       = []string{"00102", "00103"}
	plmns_6       = []string{"008043", "009899"}
	plmns_5_and_6 = []string{"00102", "001027"}
)

func TestPlmn5(t *testing.T) {
	plmnIdVals := GetPlmnVals(plmns_5)

	assert.True(t, CheckImsiOnPlmnIdListIfAny("001020000000055", plmnIdVals),
		"IMSI 001020000000055 should be valid PLMN, PLMN ID Map", plmnIdVals)
	assert.True(t, plmnIdVals.Check("001030000000055"),
		"IMSI 001030000000055 should be valid PLMN, PLMN ID Map", plmnIdVals)
	assert.False(t, CheckImsiOnPlmnIdListIfAny("001010000000055", plmnIdVals),
		"IMSI 001010000000055 should NOT be valid PLMN, PLMN ID Map", plmnIdVals)
}

func TestPlmn6(t *testing.T) {
	plmnIdVals := GetPlmnVals(plmns_6, "test plmns_6")

	assert.True(t, CheckImsiOnPlmnIdListIfAny("008043000000055", plmnIdVals),
		"IMSI 008043000000055 should be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
	assert.True(t, plmnIdVals.Check("009899000000055"),
		"IMSI 009899000000055 should be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
	assert.False(t, CheckImsiOnPlmnIdListIfAny("008040000000055", plmnIdVals),
		"IMSI 008040000000055 should NOT be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
}

func TestPlmn6and6(t *testing.T) {
	plmnIdVals := GetPlmnVals(plmns_5_and_6)

	assert.True(t, CheckImsiOnPlmnIdListIfAny("001023000000055", plmnIdVals),
		"IMSI 001023000000055 should be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
	assert.True(t, CheckImsiOnPlmnIdListIfAny("001027000000055", plmnIdVals),
		"IMSI 001027000000055 should be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
	assert.True(t, plmnIdVals.Check("001020000000055"),
		"IMSI 001020000000055 should NOT be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
	assert.False(t, CheckImsiOnPlmnIdListIfAny("001090000000055", plmnIdVals),
		"IMSI 001090000000055 should NOT be valid PLMN, PLMN ID Map: %+v", plmnIdVals)
}
