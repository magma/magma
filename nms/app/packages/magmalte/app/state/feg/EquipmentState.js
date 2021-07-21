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

import type {EnqueueSnackbarOptions} from 'notistack';
import type {FederationGatewayHealthStatus} from '../../components/GatewayUtils';
import type {
  federation_gateway,
  gateway_id,
  network_id,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import {getFederationGatewayHealthStatus} from '../../components/GatewayUtils';

type InitGatewayStateProps = {
  networkId: network_id,
  setFegGateways: ({[string]: federation_gateway}) => void,
  setFegGatewaysHealthStatus: ({
    [string]: FederationGatewayHealthStatus,
  }) => void,
  setActiveFegGatewayId: (gatewayId: gateway_id) => void,
  enqueueSnackbar: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

/**
 * Initializes the federation gateway state which is going to have a maximum of
 * 2 federation gateways, their health status, and the gateway id of the active
 * federation gateway.
 * @param {network_id} networkId Id of the federation network
 * @param {({[string]: federation_gateway}) => void} setFegGateways Sets federation gateways.
 * @param {({[string]: FederationGatewayHealthStatus}) => void} setFegGatewaysHealthStatus Sets federation gateways health status.
 * @param {(gatewayId:gateway_id) => void} setActiveFegGatewayId Sets the active gateway id.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 */
export async function InitGatewayState(props: InitGatewayStateProps) {
  const {
    networkId,
    setFegGateways,
    setFegGatewaysHealthStatus,
    setActiveFegGatewayId,
    enqueueSnackbar,
  } = props;
  try {
    const fegGateways = await MagmaV1API.getFegByNetworkIdGateways({
      networkId: networkId,
    });
    const [fegGatewaysHealthStatus, activeFegGatewayId] = await Promise.all([
      getFegGatewaysHealthStatus(networkId, fegGateways, enqueueSnackbar),
      getActiveFegGatewayId(networkId, fegGateways, enqueueSnackbar),
    ]);
    setFegGateways(fegGateways);
    setFegGatewaysHealthStatus(fegGatewaysHealthStatus);
    setActiveFegGatewayId(activeFegGatewayId);
  } catch (e) {
    enqueueSnackbar?.('failed fetching federation gateway information', {
      variant: 'error',
    });
  }
}

/**
 * Returns an object containing the IDs of the federation gateways mapped to
 * a boolean value showing each gateway's health status. A boolean value of
 * true shows that the gateway is healthy.
 *
 * @param {network_id} networkId: Id of the federation network.
 * @param {{[gateway_id]: federation_gateway}} fegGateways Federation gateways of the network.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 * @returns an object containing the IDs of the federation gateways mapped to their health status.
 */
async function getFegGatewaysHealthStatus(
  networkId: network_id,
  fegGateways: {[gateway_id]: federation_gateway},
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
): Promise<{[gateway_id]: FederationGatewayHealthStatus}> {
  const fegGatewaysHealthStatus = {};
  const fegGatewaysId = Object.keys(fegGateways);
  for (const fegGatewayId of fegGatewaysId) {
    const healthStatus = await getFederationGatewayHealthStatus(
      networkId,
      fegGatewayId,
      enqueueSnackbar,
    );
    fegGatewaysHealthStatus[fegGatewayId] = healthStatus;
  }
  return fegGatewaysHealthStatus;
}

/**
 * Fetches and returns the active federation gateway id. If it doesn't
 * have one, then it returns an empty string.
 *
 * @param {network_id} networkId: Id of the federation network.
 * @param {{[gateway_id]: federation_gateway}} fegGateways Federation gateways of the network.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 * @returns returns the active federation gateway id or an empty string.
 */
async function getActiveFegGatewayId(
  networkId: network_id,
  fegGateways: {[gateway_id]: federation_gateway},
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
): Promise<string> {
  try {
    const response = await MagmaV1API.getFegByNetworkIdClusterStatus({
      networkId,
    });
    const activeFegGatewayId = response?.active_gateway;
    // make sure active gateway id is not a dummy id
    return fegGateways[activeFegGatewayId] ? activeFegGatewayId : '';
  } catch (e) {
    enqueueSnackbar?.('failed fetching active federation gateway id', {
      variant: 'error',
    });
    return '';
  }
}
