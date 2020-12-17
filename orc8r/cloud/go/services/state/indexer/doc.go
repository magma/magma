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
	Package indexer provides tools to define, use, and update state indexers.

	State service

	The state service stores gateway state as keyed blobs.
	Examples
		- IMSI -> directory record blob
		- HWID -> gateway status blob

	Since state values are stored as arbitrary serialized blobs, the state
	service has no semantic understanding of stored values. This means
	searching over stored values would otherwise require an O(n) operation.
	Examples
		- Find IMSI with given IP -- must load all directory records into
		  memory
		- Find all gateways that haven't checked in recently -- must load all
		  gateway statuses into memory

	Derived state

	The solution is to provide customizable, online mechanisms for generating
	derived state based on existing state. Existing, "primary" state is stored
	in the state service, and derived, "secondary" state is stored in whichever
	service owns the derived state.
	Examples
		- Reverse map of directory records
			- Primary state: IMSI -> directory record
			- Secondary state: IP -> IMSI (stored in e.g. directoryd)
		- Reverse map of gateway checkin time
			- Primary state: HWID -> gateway status
			- Secondary state: checkin time -> HWID (stored in e.g. metricsd)
		- List all gateways with multiple kernel versions installed
			- Primary state: HWID -> gateway status
			- Secondary state: list of gateways (stored in e.g. bootstrapper)

	State indexers

	State indexers are Orchestrator services registering an IndexerServer under
	their gRPC endpoint. Any Orchestrator service can provide its own indexer
	servicer.

	The state service discovers indexers using K8s labels. Any service with the
	label "orc8r.io/state_indexer" will be assumed to provide an indexer
	servicer.

	Indexers provide two additional pieces of metadata -- version and types.
		- version: positive integer indicating when indexer requires reindexing
		- types: list of state types the indexer subscribes to
	These metadata are indicated by K8s annotations
		- orc8r.io/state_indexer_version -- positive integer
		- orc8r.io/state_indexer_types -- comma-separated list of state types

	Reindexing

	When an indexer's implementation changes, its derived state needs to be
	refreshed. This is accomplished by sending all existing state (of desired
	types) through the now-updated indexer.

	An indexer indicates it needs to undergo a reindex by incrementing its
	version (exposed via the above-mentioned annotation). From there, the state
	service automatically handles the reindexing process.

	Metrics and logging are available to track long-running reindex processes,
	as well an indexers CLI which reports desired and current indexer versions.

	Implementing a custom indexer

	To create a custom indexer, attach an IndexerServer to a new or existing
	Orchestrator service.

	A service can only attach a single indexer. However, that indexer can
	choose to multiplex its functionality over any desired number of "logical"
	indexers.

	See the orchestrator service for an example custom indexer.

	Notes

	The state indexer pattern currently provides no mechanism for connecting
	primary and secondary state. This means secondary state can go stale.
	Where relevant, consumers of secondary state should take this into account,
	generally by checking the primary state to ensure it agrees with the
	secondary state.
	Examples
		- Reverse map of directory records
			- Get IMSI from IP -> IMSI map (secondary state)
			- Ensure the directory record in the IMSI -> directory map contains
			  the desired IP (primary state)
		- Reverse map of gateway checkin time
			- Get HWIDs from checkin time -> HWID map (secondary state)
			- For each HWID, ensure the gateway status in the HWID -> gateway
			  status map contains the relevant checkin time (primary state)

	Automatic reindexing is only supported with Postgres. Deployments targeting
	Maria will need to use the indexer CLI to manually trigger reindex
	operations.

	There is a trivial but existent race condition during the reindex process.
	Since the index and reindex operations both use the Index gRPC method,
	and the index and reindex operations operate in parallel, it's possible for
	an indexer to receive an outdated piece of state from the reindexer.
	However, this requires
		- reindexer read old state
		- new state reported, indexer read new state
		- indexer Index call completed
		- reindexer Index call completed
	If this race condition is intolerable to the desired use case, the solution
	is to separate out the Index call into Index and Reindex methods. This is
	not currently implemented as we don't have a concrete use-case for it yet.
*/
package indexer
