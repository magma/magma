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
  federated_network_configs,
  feg_lte_network,
  network_epc_configs,
  network_id,
  network_ran_configs,
  network_subscriber_config,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

import {UpdateNetworkState as UpdateLteNetworkState} from '../lte/NetworkState';
export type UpdateNetworkProps = {
  networkId: network_id,
  lteNetwork?: feg_lte_network & {subscriber_config: network_subscriber_config},
  federation?: federated_network_configs,
  epcConfigs?: network_epc_configs,
  lteRanConfigs?: network_ran_configs,
  subscriberConfig?: network_subscriber_config,
  setLteNetwork: feg_lte_network => void,
  refreshState: boolean,
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setLteNetwork} = props;
  const requests = [];
  if (props.lteNetwork) {
    requests.push(
      await MagmaV1API.putFegLteByNetworkId({
        networkId: networkId,
        lteNetwork: {
          ...props.lteNetwork,
        },
      }),
    );
  }
  if (props.federation) {
    requests.push(
      await MagmaV1API.putFegLteByNetworkIdFederation({
        networkId: networkId,
        config: props.federation,
      }),
    );
  }
  if (props.subscriberConfig) {
    requests.push(
      await MagmaV1API.putFegLteByNetworkIdSubscriberConfig({
        networkId: props.networkId,
        record: props.subscriberConfig,
      }),
    );
  }
  if (props.epcConfigs != null || props.lteRanConfigs != null) {
    await UpdateLteNetworkState({
      networkId,
      setLteNetwork: _ => {},
      epcConfigs: props.epcConfigs,
      lteRanConfigs: props.lteRanConfigs,
      refreshState: false,
    });
  }
  await Promise.all(requests);
  if (props.refreshState) {
    const [fegLteResp, fegLteSubscriberConfigResp] = await Promise.allSettled([
      MagmaV1API.getFegLteByNetworkId({networkId}),
      MagmaV1API.getFegLteByNetworkIdSubscriberConfig({networkId}),
    ]);
    if (fegLteResp.value) {
      let subscriber_config = {};
      if (fegLteSubscriberConfigResp.value) {
        subscriber_config = fegLteSubscriberConfigResp.value;
      }
      setLteNetwork({...fegLteResp.value, subscriber_config});
    }
  }
}
