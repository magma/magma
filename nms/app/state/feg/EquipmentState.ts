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
 */

import type {
  FederationGateway,
  MutableFederationGateway,
} from '../../../generated';
import type {FederationGatewayHealthStatus} from '../../components/GatewayUtils';
import type {OptionsObject} from 'notistack';

import MagmaAPI from '../../api/MagmaAPI';
import {EnqueueSnackbar} from '../../hooks/useSnackbar';
import {GatewayId, NetworkId} from '../../../shared/types/network';
import {getFederationGatewayHealthStatus} from '../../components/GatewayUtils';

type InitGatewayStateProps = {
  networkId: NetworkId;
  setFegGateways: (fegGateways: Record<string, FederationGateway>) => void;
  setFegGatewaysHealthStatus: (
    gatewayHealthStatuses: Record<string, FederationGatewayHealthStatus>,
  ) => void;
  setActiveFegGatewayId: (gatewayId: GatewayId) => void;
  enqueueSnackbar: EnqueueSnackbar;
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
  const result = await FetchFegGateways({networkId, enqueueSnackbar});
  if (result) {
    setFegGateways(result.fegGateways);
    setFegGatewaysHealthStatus(result.fegGatewaysHealthStatus);
    setActiveFegGatewayId(result.activeFegGatewayId);
  }
}

/**
 * A prop passed when setting the gateway state.
 *
 * @property {network_id} networkId Id of the federation network
 * @property {{[gateway_id]: federation_gateway}} fegGateways Federation gateways of the network.
 * @property {{[gateway_id]: FederationGatewayHealthStatus}} fegGatewaysHealthStatus Health status of the federation gateways.
 * @property {({[string]: federation_gateway}) => void} setFegGateways Sets federation gateways.
 * @property {({[string]: FederationGatewayHealthStatus}) => void} setFegGatewaysHealthStatus Sets federation gateways health status.
 * @property {(gatewayId:gateway_id) => void} setActiveFegGatewayId Sets the active gateway id.
 * @property {gateway_id} key Id of the gateway to be added, deleted or edited.
 * @property {mutable_federation_gateway} value New Value for the gateway with the id: key.
 * @property {[gateway_id]: federation_gateway} newState New State of the Federation Gateway.
 * @property {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 */
type GatewayStateProps = {
  networkId: NetworkId;
  fegGateways: Record<GatewayId, FederationGateway>;
  fegGatewaysHealthStatus: Record<GatewayId, FederationGatewayHealthStatus>;
  setFegGateways: (fegGateways: Record<GatewayId, FederationGateway>) => void;
  setFegGatewaysHealthStatus: (
    gatewayHealthStatus: Record<GatewayId, FederationGatewayHealthStatus>,
  ) => void;
  setActiveFegGatewayId: (activeGwId: GatewayId) => void;
  key: GatewayId;
  value?: MutableFederationGateway;
  enqueueSnackbar: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

/**
 * Adds, edits, or deletes a federation gateway or sets the gateway state to a new state. It
 * then makes sure to sync the health status of the gateways and update the active gateway id
 * in case it changed.
 *
 * @param {GatewayStateProps} props an object containing the neccessary values to change the gateway state
 */
export async function SetGatewayState(props: GatewayStateProps) {
  const {
    networkId,
    fegGateways,
    fegGatewaysHealthStatus,
    setFegGateways,
    setFegGatewaysHealthStatus,
    setActiveFegGatewayId,
    key,
    value,
    enqueueSnackbar,
  } = props;
  if (value) {
    if (!(key in fegGateways)) {
      await MagmaAPI.federationGateways.fegNetworkIdGatewaysPost({
        networkId: networkId,
        gateway: value,
      });
      setFegGateways({...fegGateways, [key]: value as FederationGateway});
    } else {
      await MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdPut({
        networkId: networkId,
        gatewayId: key,
        gateway: value,
      });
      setFegGateways({...fegGateways, [key]: value as FederationGateway});
    }
    const newFegGatewaysHealthStatus = {...fegGatewaysHealthStatus};
    newFegGatewaysHealthStatus[key] = await getFederationGatewayHealthStatus(
      networkId,
      key,
      enqueueSnackbar,
    );
    setFegGatewaysHealthStatus(newFegGatewaysHealthStatus);
  } else {
    await MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdDelete({
      networkId: networkId,
      gatewayId: key,
    });
    const newFegGateways = {...fegGateways};
    const newFegGatewaysHealthStatus = {...fegGatewaysHealthStatus};
    delete newFegGateways[key];
    delete newFegGatewaysHealthStatus[key];
    setFegGateways(newFegGateways);
    setFegGatewaysHealthStatus(newFegGatewaysHealthStatus);
  }
  setActiveFegGatewayId(
    await getActiveFegGatewayId(networkId, fegGateways, enqueueSnackbar),
  );
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
  networkId: NetworkId,
  fegGateways: Record<GatewayId, FederationGateway>,
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined,
): Promise<Record<GatewayId, FederationGatewayHealthStatus>> {
  const fegGatewaysHealthStatus: Record<
    GatewayId,
    FederationGatewayHealthStatus
  > = {};
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
  networkId: NetworkId,
  fegGateways: Record<GatewayId, FederationGateway>,
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined,
): Promise<string> {
  try {
    const response = (
      await MagmaAPI.federationNetworks.fegNetworkIdClusterStatusGet({
        networkId,
      })
    ).data;
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

interface FetchAllGatewaysProps {
  networkId: NetworkId;
  enqueueSnackbar: EnqueueSnackbar;
}

/**
 * Fetches and returns the list of gateways under the federation network or
 * the specific gateway if the id is provided.
 *
 * @param {network_id} networkId: Id of the federation network.
 * @param {gateway_id} id id of the federation gateway.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 * @returns {{[string]: federation_gateway}} returns an object containing the federation
 *   gateways in the network or the federation gateway with the given id. It returns an empty
 *   object and displays any error encountered on the snackbar when it fails to fetch the gateways.
 */
export async function FetchFegGateways(props: FetchAllGatewaysProps) {
  const {networkId, enqueueSnackbar} = props;
  try {
    const fegGateways = (
      await MagmaAPI.federationGateways.fegNetworkIdGatewaysGet({
        networkId: networkId,
      })
    ).data;

    const [fegGatewaysHealthStatus, activeFegGatewayId] = await Promise.all([
      getFegGatewaysHealthStatus(networkId, fegGateways, enqueueSnackbar),
      getActiveFegGatewayId(networkId, fegGateways, enqueueSnackbar),
    ]);

    return {fegGateways, fegGatewaysHealthStatus, activeFegGatewayId};
  } catch (e) {
    enqueueSnackbar?.('Failed fetching gateway information', {
      variant: 'error',
    });
  }
}

interface FetchSingleGatewayProps {
  networkId: NetworkId;
  enqueueSnackbar: EnqueueSnackbar;
  id: GatewayId;
}

export async function FetchFegGateway(props: FetchSingleGatewayProps) {
  const {networkId, enqueueSnackbar, id} = props;

  try {
    const [gatewayResponse, healthStatus] = await Promise.all([
      MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdGet({
        networkId: networkId,
        gatewayId: id,
      }),
      getFederationGatewayHealthStatus(networkId, id, enqueueSnackbar),
    ]);
    if (gatewayResponse) {
      return {fegGateway: gatewayResponse.data, healthStatus};
    }
  } catch (e) {
    enqueueSnackbar(
      `Failed fetching gateway information for the gateway with id: ${id}`,
      {variant: 'error'},
    );
  }
}
