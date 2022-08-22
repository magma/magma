/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package restrictor

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/assert"
)

type restrictorTestCase struct {
	name          string
	input         string
	expected      string
	expectedError string
	restrictor    *QueryRestrictor
}

func (tc *restrictorTestCase) RunTest(t *testing.T) {
	output, err := tc.restrictor.RestrictQuery(tc.input)
	if tc.expectedError == "" {
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, output)
		return
	}
	assert.EqualError(t, err, tc.expectedError)
}

var singleMatcher, _ = labels.NewMatcher(labels.MatchEqual, "networkID", "test")

func TestQueryRestrictor_RestrictQuery(t *testing.T) {
	singleLabelRestrictor := NewQueryRestrictor(DefaultOpts).AddMatcher("networkID", "test")
	testCases := []*restrictorTestCase{
		{
			name:       "basic query",
			input:      "up",
			expected:   `up{networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with function",
			input:      "sum(up)",
			expected:   `sum(up{networkID="test"})`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with labels",
			input:      `up{label="value"}`,
			expected:   `up{label="value",networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with multiple metrics",
			input:      "metric1 or metric2",
			expected:   `metric1{networkID="test"} or metric2{networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with multiple metrics and labels",
			input:      `metric1 or metric2{label="value"}`,
			expected:   `metric1{networkID="test"} or metric2{label="value",networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with matrix selector",
			input:      "up[5m]",
			expected:   `up{networkID="test"}[5m]`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with matrix and functions",
			input:      "sum_over_time(metric1[5m])",
			expected:   `sum_over_time(metric1{networkID="test"}[5m])`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with existing networkID",
			input:      `metric1{networkID="test"}`,
			expected:   `metric1{networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "query with existing wrong networkID",
			input:      `metric1{networkID="malicious"}`,
			expected:   `metric1{networkID="test"}`,
			restrictor: singleLabelRestrictor,
		},
		{
			name:       "restricts with multiple labels",
			input:      `metric1`,
			expected:   `metric1{newLabel1="value1",newLabel2="value2"}`,
			restrictor: NewQueryRestrictor(DefaultOpts).AddMatcher("newLabel1", "value1").AddMatcher("newLabel2", "value2"),
		},
		{
			name:       "creates an OR with multiple values",
			input:      `metric1`,
			expected:   `metric1{newLabel1=~"value1|value2"}`,
			restrictor: NewQueryRestrictor(DefaultOpts).AddMatcher("newLabel1", "value1", "value2"),
		},
		{
			name:       "creates an OR along with another label",
			input:      `metric1{newLabel1="value1"}`,
			expected:   `metric1{newLabel1="value1",newLabel2=~"value2|value3"}`,
			restrictor: NewQueryRestrictor(DefaultOpts).AddMatcher("newLabel2", "value2", "value3"),
		},
		{
			name:       "doesn't overwrite existing label if configured",
			input:      `metric1{newLabel1="value1"}`,
			expected:   `metric1{newLabel1="value1",newLabel1=~"value2|value3"}`,
			restrictor: NewQueryRestrictor(Opts{ReplaceExistingLabel: false}).AddMatcher("newLabel1", "value2", "value3"),
		},
		{
			name:       "Empty matcher value works",
			input:      `metric1`,
			expected:   `metric1{newLabel1=""}`,
			restrictor: NewQueryRestrictor(DefaultOpts).AddMatcher("newLabel1"),
		},
		{
			name:          "Empty query",
			input:         ``,
			expectedError: `empty query string`,
			restrictor:    singleLabelRestrictor,
		},
		{
			name:          "invalid query",
			input:         `!test`,
			expectedError: `error parsing query: parse error at char 1: unexpected character after '!': 't'`,
			restrictor:    singleLabelRestrictor,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}
