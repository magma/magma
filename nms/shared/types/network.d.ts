/**
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
 *
 * @flow
 * @format
 */

const Networks = {
  carrier_wifi_network: 'carrier_wifi_network',
  xwfm: 'xwfm',
  feg: 'feg',
  feg_lte: 'feg_lte',
  lte: 'lte',
  rhino: 'rhino',
  symphony: 'symphony',
  third_party: 'third_party', // TODO: deprecate third_party in lieu of symphony
  wifi_network: 'wifi_network',
} as const;

export const CWF = Networks.carrier_wifi_network;
export const XWFM = Networks.xwfm;
export const FEG = Networks.feg;
export const LTE = Networks.lte;
export const FEG_LTE = Networks.feg_lte;
export const RHINO = Networks.rhino;
export const SYMPHONY = Networks.symphony;
export const THIRD_PARTY = Networks.third_party;
export const WIFI = Networks.wifi_network;

export const AllNetworkTypes: Array<NetworkType>;
export const V1NetworkTypes: Array<NetworkType>;

export function coalesceNetworkType(
  networkID: string,
  networkType: string | null | undefined,
): NetworkType | null | undefined;
