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

package orchestrator

const (
	// PrometheusPushAddresses is the orchestrator.yml key for the list of
	// addresses to which the metrics exporter will push Prometheus metrics
	PrometheusPushAddresses = "prometheusPushAddresses"

	// PrometheusGRPCPushAddress is the orchestrator.yml key for the GRPC address
	// to which the metrics exporter will push Prometheus metrics
	PrometheusGRPCPushAddress = "prometheusGRPCPushAddress"

	// UseGRPCExporter is a flag to determine which exporter to use
	UseGRPCExporter = "useGRPCExporter"
)
