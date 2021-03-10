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

/*
	Package configurator supports configuration management and dynamic
	generation by manipulating a graph of network entities.

	Entity graph

	Configurator manages a directed acyclic graph (DAG) of network-partitioned
	entities. Callers can define their graph via the following three types
		- Network
			- Network-level configs
			- Network-level metadata
		- Network entity (vertices)
			- Type
			- Key (ID)
			- Config (serialized)
		- Edge (directed edge)
			- Connect two network entities
	Each network provides an isolated graph, with network-level configs and
	metadata. Network entities represent logical entities within a network
	such as a gateway, subscriber, or APN. Edges connect two network
	entities with a directed edge.

	Code calling configurator should take care to define entity types and
	relations which would invariably result in an acyclic graph.

	Configurator supports token based pagination for entity loads. Clients can
	specify a page size and token for a given entity type. The load response
	will contain the entity page and a page token to be used on subsequent
	requests. If a page size is not specified, configurator defaults to the
	maximum supported page size (15k). The max size is configurable via
	configurator's service config.

	Generating configs

	Configurator provides two interfaces: northbound and southbound. The
	northbound interface manipulates the entity graphs according to requests
	from the Orchestrator REST API and from other Orchestrator services. The
	southbound interface synthesizes the entity graph into mconfigs for
	particular gateways.

	Configurator generates these gateway mconfigs ("Magma" configs) by
	outsourcing config generation to a dynamic set of mconfig builders.
	Configurator sends each registered builder the gateway ID for which to
	build a config, along with the encompassing entity graph and network.
	With this information, each builder can traverse the graph as-necessary to
	dynamically build a config for the requesting gateway.

	Configurator assembles the set of partial configs from each mconfig builder
	into a complete config. Before returning, it adds metadata such as
	time of creation and hash/digest of the configs.

	Mconfig builders are Orchestrator services registering an MconfigBuilder
	under their gRPC endpoint. Any Orchestrator service can provide its own
	builder servicer. Configurator discovers mconfig builders using K8s labels.
	Any service with the label "orc8r.io/mconfig_builder" will be assumed to
	provide an mconfig builder servicer.
*/
package configurator

const (
	// ServiceName is the name of this service.
	ServiceName = "CONFIGURATOR"

	// NetworkConfigSerdeDomain is the Serde domain for network configs.
	NetworkConfigSerdeDomain = "configurator_network_configs"

	// NetworkEntitySerdeDomain is the Serde domain for network entity configs.
	NetworkEntitySerdeDomain = "configurator_entity_configs"
)
