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

import type {GatewayId, NetworkId} from '../../../shared/types/network';
import type {GenericCommandParams, PingRequest} from '../../../generated';

import MagmaAPI from '../../api/MagmaAPI';

export type GatewayCommandProps = {
  networkId: NetworkId;
  gatewayId: GatewayId;
  command: 'reboot' | 'ping' | 'restartServices' | 'generic';
  pingRequest?: PingRequest;
  params?: GenericCommandParams;
};

export async function RunGatewayCommands(props: GatewayCommandProps) {
  const {networkId, gatewayId} = props;

  switch (props.command) {
    case 'reboot':
      return await MagmaAPI.commands.networksNetworkIdGatewaysGatewayIdCommandRebootPost(
        {networkId, gatewayId},
      );

    case 'restartServices':
      return await MagmaAPI.commands.networksNetworkIdGatewaysGatewayIdCommandRestartServicesPost(
        {networkId, gatewayId, services: []},
      );

    case 'ping':
      if (props.pingRequest != null) {
        return await MagmaAPI.commands.networksNetworkIdGatewaysGatewayIdCommandPingPost(
          {networkId, gatewayId, pingRequest: props.pingRequest},
        );
      }

    default:
      if (props.params != null) {
        return await MagmaAPI.commands.networksNetworkIdGatewaysGatewayIdCommandGenericPost(
          {networkId, gatewayId, parameters: props.params},
        );
      }
  }
}
