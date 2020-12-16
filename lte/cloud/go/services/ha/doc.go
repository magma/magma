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

/* Package ha provides the LTE HA orc8r service.

This service has a single RPC endpoint. The RPC endpoint will be used by
secondary gateways' MME to know when to offload its users for a given ENB
back to the primary.

To gather this state, this service looks at the primary gateways in the
gateway pool(s) of the calling gateway. For each primary it fetches the
configured ENBs and then checks the following state:

1. Has the primary checked in within last 3 mins
2. Is the ENB connected to the primary within last 3 mins
3. Does the ENB have throughput on the primary

The service then sends back the ENB ID -> offload state
for all of these.
*/
package ha

const (
	ServiceName = "ha"
)
