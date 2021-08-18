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
 * @flow strict-local
 * @format
 */

import type {enodeb, enodeb_state} from '@fbcnms/magma-api';
import type {enodeb_configuration} from '@fbcnms/magma-api';

export const EnodebDeviceClass: {
  [string]: $PropertyType<enodeb_configuration, 'device_class'>,
} = Object.freeze({
  BAICELLS_NOVA_233_2_OD_FDD: 'Baicells Nova-233 G2 OD FDD',
  BAICELLS_NOVA_243_OD_TDD: 'Baicells Nova-243 OD TDD',
  BAICELLS_NEUTRINO_224_ID_FDD: 'Baicells Neutrino 224 ID FDD',
  BAICELLS_ID: 'Baicells ID TDD/FDD',
  NURAN_CAVIUM_OC_LTE: 'NuRAN Cavium OC-LTE',
});

export const EnodebBandwidthOption: {
  [string]: $NonMaybeType<$PropertyType<enodeb_configuration, 'bandwidth_mhz'>>,
} = Object.freeze({
  '3': 3,
  '5': 5,
  '10': 10,
  '15': 15,
  '20': 20,
});

export type EnodebInfo = {
  enb: enodeb,
  enb_state: enodeb_state,
};

export function isEnodebHealthy(enbInfo: EnodebInfo) {
  return (
    enbInfo.enb.enodeb_config?.managed_config?.transmit_enabled ===
    enbInfo.enb_state.rf_tx_on
  );
}
