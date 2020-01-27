/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package restrictor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type restrictorTestCase struct {
	name       string
	input      string
	expected   string
	restrictor *QueryRestrictor
}

func (tc *restrictorTestCase) RunTest(t *testing.T) {
	output, err := tc.restrictor.RestrictQuery(tc.input)
	assert.NoError(t, err)
	assert.Equal(t, tc.expected, output)
}

func TestQueryRestrictor_RestrictQuery(t *testing.T) {
	singleLabelRestictor := NewQueryRestrictor(map[string]string{"networkID": "test"})
	testCases := []*restrictorTestCase{
		{
			name:       "basic query",
			input:      "up",
			expected:   `up{networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with function",
			input:      "sum(up)",
			expected:   `sum(up{networkID="test"})`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with labels",
			input:      `up{label="value"}`,
			expected:   `up{label="value",networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with multiple metrics",
			input:      "metric1 or metric2",
			expected:   `metric1{networkID="test"} or metric2{networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with multiple metrics and labels",
			input:      `metric1 or metric2{label="value"}`,
			expected:   `metric1{networkID="test"} or metric2{label="value",networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with matrix selector",
			input:      "up[5m]",
			expected:   `up{networkID="test"}[5m]`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with matrix and functions",
			input:      "sum_over_time(metric1[5m])",
			expected:   `sum_over_time(metric1{networkID="test"}[5m])`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with existing networkID",
			input:      `metric1{networkID="test"}`,
			expected:   `metric1{networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "query with existing wrong networkID",
			input:      `metric1{networkID="malicious"}`,
			expected:   `metric1{networkID="test"}`,
			restrictor: singleLabelRestictor,
		},
		{
			name:       "restricts with multiple labels",
			input:      `metric1`,
			expected:   `metric1{newLabel1="value1",newLabel2="value2"}`,
			restrictor: NewQueryRestrictor(map[string]string{"newLabel1": "value1", "newLabel2": "value2"}),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}
