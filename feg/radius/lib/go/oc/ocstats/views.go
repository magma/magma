/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocstats

import (
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// HTTPServerResponseCountByStatusAndPath is an additional view for server response status code and path.
var HTTPServerResponseCountByStatusAndPath = &view.View{
	Name:        "opencensus.io/http/server/response_count_by_status_code_path",
	Description: "Server response count by status code and path",
	TagKeys:     []tag.Key{ochttp.StatusCode, ochttp.KeyServerRoute},
	Measure:     ochttp.ServerLatency,
	Aggregation: view.Count(),
}
