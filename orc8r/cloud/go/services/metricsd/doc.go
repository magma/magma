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

/*
	Package metricsd supports metrics collection, augmentation, and export,
	as well as providing REST API endpoints for viewing metrics and
	managing alerts.

	A metrics profile defines the manner and extent to which metrics are
	consumed and exported.

	Metrics profiles are registered with the profile registry. The metricsd
	service consumes profiles from the registry and kicks off their respective
	collect/export loops.

	Only one metrics profile can be active at a time. The active profile is
	determined by a service-level config value, and changing the active profile
	requires a metricsd service restart. Each exporter in the active profile
	receives metrics from every collector.

	Available exporters include export to a custom Prometheus push gateway, as
	well as support for more time-series-oriented endpoints.

	The metricsd service provides gRPC endpoints to accept metrics pushed
	via the REST API, as well as metrics collected (reported) from non-local
	services such as those on AGWs.

	The obsidian REST API handlers provide endpoints to query metrics, as well
	as view, configure, and silence alerts.
*/
package metricsd

const (
	ServiceName = "METRICSD"
)
