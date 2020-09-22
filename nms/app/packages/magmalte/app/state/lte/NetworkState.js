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
 *
 * @flow strict-local
 * @format
 */

import type {
  lte_network,
  network_dns_config,
  network_epc_configs,
  network_id,
  network_ran_configs,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

export type UpdateNetworkProps = {
  networkId: network_id,
  lteNetwork?: lte_network,
  epcConfigs?: network_epc_configs,
  lteRanConfigs?: network_ran_configs,
  lteDnsConfig?: network_dns_config,
  setLteNetwork: lte_network => void,
  refreshState: boolean,
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setLteNetwork} = props;
  const requests = [];
  if (props.lteNetwork !== undefined) {
    requests.push(
      await MagmaV1API.putLteByNetworkId({
        networkId: networkId,
        lteNetwork: {
          ...props.lteNetwork,
        },
      }),
    );
  }

  if (props.epcConfigs !== undefined) {
    requests.push(
      await MagmaV1API.putLteByNetworkIdCellularEpc({
        networkId: props.networkId,
        config: props.epcConfigs,
      }),
    );
  }
  if (props.lteRanConfigs !== undefined) {
    requests.push(
      await MagmaV1API.putLteByNetworkIdCellularRan({
        networkId: props.networkId,
        config: props.lteRanConfigs,
      }),
    );
  }
  if (props.lteDnsConfig !== undefined) {
    requests.push(
      await MagmaV1API.putLteByNetworkIdDns({
        networkId: props.networkId,
        config: props.lteDnsConfig,
      }),
    );
  }

  await Promise.all(requests);
  if (props.refreshState) {
    setLteNetwork(await MagmaV1API.getLteByNetworkId({networkId}));
  }
}
