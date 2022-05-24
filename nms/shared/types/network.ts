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
 */

const Networks = {
  carrier_wifi_network: 'carrier_wifi_network',
  xwfm: 'xwfm',
  feg: 'feg',
  feg_lte: 'feg_lte',
  lte: 'lte',
} as const;

export const CWF = Networks.carrier_wifi_network;
export const XWFM = Networks.xwfm;
export const FEG = Networks.feg;
export const LTE = Networks.lte;
export const FEG_LTE = Networks.feg_lte;

export const AllNetworkTypes = Object.keys(Networks).sort() as Array<
  NetworkType
>;

export type NetworkType = keyof typeof Networks;

export function coalesceNetworkType(
  networkID: string,
  networkType?: string,
): NetworkType | null {
  if (networkType && networkType in Networks) {
    return networkType as NetworkType;
  }
  return null;
}

export type NetworkId = string;
export type GatewayId = string;
export type GatewayPoolId = string;
export type SubscriberId = string;
export type PolicyId = string;
