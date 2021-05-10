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

package query_api

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type PrometheusAPI interface {
	Query(ctx context.Context, query string, ts time.Time) (model.Value, v1.Warnings, error)
	QueryRange(ctx context.Context, query string, r v1.Range) (model.Value, v1.Warnings, error)
}

// QueryPrometheusVector handles all the error cases of making an instant query
// with a PrometheusAPI
func QueryPrometheusVector(prometheusClient PrometheusAPI, query string) (model.Vector, error) {
	// TODO: catch the warning at _
	val, _, err := prometheusClient.Query(context.Background(), query, time.Now())
	if err != nil {
		return nil, err
	}
	vec, ok := val.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected ValueType: %v", val.Type())
	}
	if len(vec) == 0 {
		return nil, fmt.Errorf("no data returned from query")
	}
	return vec, nil
}

// QueryPrometheusMatrix handles all the error cases of making a range query
// with a PrometheusAPI
func QueryPrometheusMatrix(prometheusClient PrometheusAPI, query string, r v1.Range) (model.Matrix, error) {
	// TODO: catch the warning at _
	val, _, err := prometheusClient.QueryRange(context.Background(), query, r)
	if err != nil {
		return nil, err
	}
	matrix, ok := val.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected ValueType: %v", val.Type())
	}
	if len(matrix) == 0 {
		return nil, fmt.Errorf("no data returned from query")
	}
	return matrix, nil
}
