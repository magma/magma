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
import type {
  FederatedNetworkConfigs,
  FegLteNetwork,
  NetworkEpcConfigs,
  NetworkRanConfigs,
  NetworkSubscriberConfig,
} from '../../../generated-ts';
import type {NetworkId} from '../../../shared/types/network';

import MagmaAPI from '../../../api/MagmaAPI';
import {UpdateNetworkState as UpdateLteNetworkState} from '../lte/NetworkState';

export type UpdateNetworkProps = {
  networkId: NetworkId;
  lteNetwork?: FegLteNetwork & {
    subscriber_config: NetworkSubscriberConfig;
  };
  federation?: FederatedNetworkConfigs;
  epcConfigs?: NetworkEpcConfigs;
  lteRanConfigs?: NetworkRanConfigs;
  subscriberConfig?: NetworkSubscriberConfig;
  setLteNetwork: (
    arg0: FegLteNetwork & {
      subscriber_config: NetworkSubscriberConfig;
    },
  ) => void;
  refreshState: boolean;
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setLteNetwork} = props;
  const requests = [];
  if (props.lteNetwork) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdPut({
        networkId: networkId,
        lteNetwork: {...props.lteNetwork},
      }),
    );
  }
  if (props.federation) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdFederationPut({
        networkId: networkId,
        config: props.federation,
      }),
    );
  }
  if (props.subscriberConfig) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdSubscriberConfigPut({
        networkId: props.networkId,
        record: props.subscriberConfig,
      }),
    );
  }
  if (props.epcConfigs != null || props.lteRanConfigs != null) {
    await UpdateLteNetworkState({
      networkId,
      setLteNetwork: () => {},
      epcConfigs: props.epcConfigs,
      lteRanConfigs: props.lteRanConfigs,
      refreshState: false,
    });
  }
  await Promise.all(requests);
  if (props.refreshState) {
    const [fegLteResp, fegLteSubscriberConfigResp] = await Promise.allSettled([
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet({
        networkId,
      }),
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdSubscriberConfigGet({
        networkId,
      }),
    ]);
    if (fegLteResp.status === 'fulfilled') {
      let subscriber_config = {};
      if (fegLteSubscriberConfigResp.status === 'fulfilled') {
        subscriber_config = fegLteSubscriberConfigResp.value.data;
      }
      setLteNetwork({...fegLteResp.value.data, subscriber_config});
    }
  }
}
