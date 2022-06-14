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

import MagmaAPI from '../../../api/MagmaAPI';
import type {
  LteNetwork,
  NetworkDnsConfig,
  NetworkEpcConfigs,
  NetworkRanConfigs,
  NetworkSubscriberConfig,
} from '../../../generated-ts';
import type {NetworkId} from '../../../shared/types/network';

export type UpdateNetworkProps = {
  networkId: NetworkId;
  lteNetwork?: LteNetwork;
  epcConfigs?: NetworkEpcConfigs;
  lteRanConfigs?: NetworkRanConfigs;
  lteDnsConfig?: NetworkDnsConfig;
  subscriberConfig?: NetworkSubscriberConfig;
  setLteNetwork: (lteNetwork: LteNetwork) => void;
  refreshState: boolean;
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setLteNetwork} = props;
  const requests = [];

  if (props.lteNetwork) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdPut({
        networkId: networkId,
        lteNetwork: {...props.lteNetwork},
      }),
    );
  }

  if (props.epcConfigs) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdCellularEpcPut({
        networkId: props.networkId,
        config: props.epcConfigs,
      }),
    );
  }
  if (props.lteRanConfigs) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdCellularRanPut({
        networkId: props.networkId,
        config: props.lteRanConfigs,
      }),
    );
  }
  if (props.lteDnsConfig) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdDnsPut({
        networkId: props.networkId,
        config: props.lteDnsConfig,
      }),
    );
  }
  if (props.subscriberConfig) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdSubscriberConfigPut({
        networkId: props.networkId,
        record: props.subscriberConfig,
      }),
    );
  }
  // TODO(andreilee): Provide a way to handle errors here
  await Promise.all(requests);
  if (props.refreshState) {
    setLteNetwork(
      (
        await MagmaAPI.lteNetworks.lteNetworkIdGet({
          networkId,
        })
      ).data,
    );
  }
}
