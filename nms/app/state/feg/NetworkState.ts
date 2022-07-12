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

import MagmaAPI from '../../api/MagmaAPI';
import type {FegNetwork, NetworkSubscriberConfig} from '../../../generated';
import type {NetworkId} from '../../../shared/types/network';

export type UpdateNetworkProps = {
  networkId: NetworkId;
  fegNetwork?: FegNetwork;
  subscriberConfig?: NetworkSubscriberConfig;
  setFegNetwork: (fn: FegNetwork) => void;
  refreshState: boolean;
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setFegNetwork} = props;
  const requests = [];
  if (props.fegNetwork) {
    requests.push(
      await MagmaAPI.federationNetworks.fegNetworkIdPut({
        networkId: networkId,
        fegNetwork: {
          ...props.fegNetwork,
        },
      }),
    );
  }
  if (props.subscriberConfig) {
    requests.push(
      await MagmaAPI.federationNetworks.fegNetworkIdSubscriberConfigPut({
        networkId: props.networkId,
        record: props.subscriberConfig,
      }),
    );
  }
  await Promise.all(requests);
  if (props.refreshState) {
    setFegNetwork(
      (await MagmaAPI.federationNetworks.fegNetworkIdGet({networkId})).data,
    );
  }
}
