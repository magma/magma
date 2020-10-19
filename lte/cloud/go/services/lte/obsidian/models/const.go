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

package models

const (
	// NATAllocationMode NAT IP allocation mode
	NATAllocationMode = "NAT"
	// StaticAllocationMode Static IP allocation mode
	StaticAllocationMode = "STATIC"
	// DHCPPassthroughAllocationMode DHCP Passthrough (carrier wifi) IP allocation mode
	DHCPPassthroughAllocationMode = "DHCP_PASSTHROUGH"
	// DHCPBroadcastAllocationMode DHCP Broadcast IP allocation mode
	DHCPBroadcastAllocationMode = "DHCP_BROADCAST"

	// ManagedConfigType Configuration type for managed radios
	ManagedConfigType = "MANAGED"
	// UnmanagedConfigType Configuration type for externally managed radios
	UnmanagedConfigType = "UNMANAGED"
)
