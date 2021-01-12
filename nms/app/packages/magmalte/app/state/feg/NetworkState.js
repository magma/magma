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
  feg_network,
  network_id,
  network_subscriber_config,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

export type UpdateNetworkProps = {
  networkId: network_id,
  subscriberConfig?: network_subscriber_config,
  setFegNetwork: feg_network => void,
  refreshState: boolean,
};

export async function UpdateNetworkState(props: UpdateNetworkProps) {
  const {networkId, setFegNetwork} = props;
  const requests = [];
  if (props.subscriberConfig) {
    requests.push(
      await MagmaV1API.putFegByNetworkIdSubscriberConfig({
        networkId: props.networkId,
        record: props.subscriberConfig,
      }),
    );
  }
  await Promise.all(requests);
  if (props.refreshState) {
    setFegNetwork(await MagmaV1API.getFegByNetworkId({networkId}));
  }
}
