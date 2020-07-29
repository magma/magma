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
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

// QueryRestrictor provides functionality to add restrictor labels to a
// Prometheus query
type QueryRestrictor struct {
	matchers []labels.Matcher
	Opts
}

// Opts contains optional configurations for the QueryRestrictor
type Opts struct {
	ReplaceExistingLabel bool
}

var DefaultOpts = Opts{ReplaceExistingLabel: true}

// NewQueryRestrictor returns a new QueryRestrictor to be built upon
func NewQueryRestrictor(opts Opts) *QueryRestrictor {
	return &QueryRestrictor{
		matchers: []labels.Matcher{},
		Opts:     opts,
	}
}

// AddMatcher takes a key and an arbitrary number of values. If only one value
// is provided, an Equal matcher will be added to the restrictor. Otherwise a
// Regex matcher with an OR of the values will be added. e.g. {label=~"value1|value2"}
// If values is empty the label will be matched to the empty string e.g. {label=""}
// effectively, this matches which do not contain this label
func (q *QueryRestrictor) AddMatcher(key string, values ...string) *QueryRestrictor {
	if len(values) < 1 {
		q.matchers = append(q.matchers, labels.Matcher{Type: labels.MatchEqual, Name: key})
		return q
	}

	if len(values) == 1 {
		q.matchers = append(q.matchers, labels.Matcher{Type: labels.MatchEqual, Name: key, Value: values[0]})
		return q
	}

	q.matchers = append(q.matchers, labels.Matcher{Type: labels.MatchRegexp, Name: key, Value: strings.Join(values, "|")})
	return q
}

// RestrictQuery appends a label selector to each metric in a given query so
// that only metrics with those labels are returned from the query.
func (q *QueryRestrictor) RestrictQuery(query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("empty query string")
	}

	promQuery, err := promql.ParseExpr(query)
	if err != nil {
		return "", fmt.Errorf("error parsing query: %v", err)
	}
	promql.Inspect(promQuery, q.addRestrictorLabels())
	return promQuery.String(), nil
}

// Matchers returns the list of label matchers for the restrictor
func (q *QueryRestrictor) Matchers() []labels.Matcher {
	return q.matchers
}

func (q *QueryRestrictor) addRestrictorLabels() func(n promql.Node, path []promql.Node) error {
	return func(n promql.Node, path []promql.Node) error {
		if n == nil {
			return nil
		}
		for _, matcher := range q.matchers {
			switch n := n.(type) {
			case *promql.VectorSelector:
				n.LabelMatchers = appendOrReplaceMatcher(n.LabelMatchers, matcher, q.ReplaceExistingLabel)
			case *promql.MatrixSelector:
				n.LabelMatchers = appendOrReplaceMatcher(n.LabelMatchers, matcher, q.ReplaceExistingLabel)
			}
		}
		return nil
	}
}

func appendOrReplaceMatcher(matchers []*labels.Matcher, newMatcher labels.Matcher, replaceExistingLabel bool) []*labels.Matcher {
	if replaceExistingLabel && getMatcherIndex(matchers, newMatcher.Name) >= 0 {
		return replaceLabelValue(matchers, newMatcher.Name, newMatcher.Value)
	} else {
		return append(matchers, &newMatcher)
	}
}

func getMatcherIndex(matchers []*labels.Matcher, name string) int {
	for idx, match := range matchers {
		if match.Name == name {
			return idx
		}
	}
	return -1
}

func replaceLabelValue(matchers []*labels.Matcher, name, value string) []*labels.Matcher {
	idx := getMatcherIndex(matchers, name)
	if idx >= -1 {
		matchers[idx].Value = value
	}
	return matchers
}
