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

import * as React from 'react';
import LoadingFiller from '../components/LoadingFiller';
import MagmaAPI from '../api/MagmaAPI';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {FederationGateway, MutableFederationGateway} from '../../generated';
import {
  FederationGatewayHealthStatus,
  getFederationGatewayHealthStatus,
} from '../components/GatewayUtils';
import {GatewayId, NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';

export type FEGGatewayContextType = {
  state: Record<string, FederationGateway>;
  setState: (key: GatewayId, val?: MutableFederationGateway) => Promise<void>;
  updateGateway: (props: UpdateFegGatewayParams) => Promise<void>;
  refetch: (id?: GatewayId) => void;
  health: Record<GatewayId, FederationGatewayHealthStatus>;
  activeFegGatewayId: GatewayId;
};

const FEGGatewayContext = React.createContext<FEGGatewayContextType>(
  {} as FEGGatewayContextType,
);

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
async function initGatewayState(params: {
  networkId: NetworkId;
  setFegGateways: (fegGateways: Record<string, FederationGateway>) => void;
  setFegGatewaysHealthStatus: (
    gatewayHealthStatuses: Record<string, FederationGatewayHealthStatus>,
  ) => void;
  setActiveFegGatewayId: (gatewayId: GatewayId) => void;
  enqueueSnackbar: EnqueueSnackbar;
}) {
  const {
    networkId,
    setFegGateways,
    setFegGatewaysHealthStatus,
    setActiveFegGatewayId,
    enqueueSnackbar,
  } = params;
  const result = await fetchFegGateways({networkId, enqueueSnackbar});
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
 * @property {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 */
type GatewayStateParams = {
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
  enqueueSnackbar: EnqueueSnackbar;
};

/**
 * Adds, edits, or deletes a federation gateway or sets the gateway state to a new state. It
 * then makes sure to sync the health status of the gateways and update the active gateway id
 * in case it changed.
 *
 * @param {GatewayStateParams} params an object containing the necessary values to change the gateway state
 */
async function setGatewayState(params: GatewayStateParams) {
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
  } = params;
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

export type UpdateFegGatewayParams = {
  gatewayId: GatewayId;
  tierId: string;
};

async function updateGateway(
  params: {networkId: NetworkId} & UpdateFegGatewayParams,
) {
  const {networkId, gatewayId, tierId} = params;
  await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdTierPut({
    networkId,
    gatewayId: gatewayId,
    tierId: JSON.stringify(`"${tierId}"`),
  });
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
  enqueueSnackbar?: EnqueueSnackbar,
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
  enqueueSnackbar?: EnqueueSnackbar,
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
async function fetchFegGateways(params: {
  networkId: NetworkId;
  enqueueSnackbar: EnqueueSnackbar;
}) {
  const {networkId, enqueueSnackbar} = params;
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

async function fetchFegGateway(params: {
  networkId: NetworkId;
  enqueueSnackbar: EnqueueSnackbar;
  id: GatewayId;
}) {
  const {networkId, enqueueSnackbar, id} = params;

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

/**
 * Fetches and returns the federation gateways, their health status and
 * the active federation gateway id.
 * @param {network_id} networkId Id of the network
 */
export function FEGGatewayContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
  const {networkId} = props;
  const [fegGateways, setFegGateways] = useState<
    Record<GatewayId, FederationGateway>
  >({});
  const [fegGatewaysHealthStatus, setFegGatewaysHealthStatus] = useState<
    Record<GatewayId, FederationGatewayHealthStatus>
  >({});
  const [activeFegGatewayId, setActiveFegGatewayId] = useState<GatewayId>('');
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      await initGatewayState({
        networkId,
        setFegGateways,
        setFegGatewaysHealthStatus,
        setActiveFegGatewayId,
        enqueueSnackbar,
      });
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <FEGGatewayContext.Provider
      value={{
        state: fegGateways,
        setState: (key, value?) => {
          return setGatewayState({
            networkId,
            fegGateways,
            fegGatewaysHealthStatus,
            setFegGateways,
            setFegGatewaysHealthStatus,
            setActiveFegGatewayId,
            key,
            value,
            enqueueSnackbar,
          });
        },
        updateGateway: props =>
          updateGateway({
            networkId,
            ...props,
          }),
        refetch: (id?: GatewayId) => {
          if (id) {
            void fetchFegGateway({id, networkId, enqueueSnackbar}).then(
              response => {
                if (response) {
                  setFegGateways(gateways => ({
                    ...gateways,
                    [id]: response.fegGateway,
                  }));
                  setFegGatewaysHealthStatus(healthStatus => ({
                    ...healthStatus,
                    [id]: response.healthStatus,
                  }));
                }
              },
            );
          } else {
            void fetchFegGateways({networkId, enqueueSnackbar}).then(
              response => {
                if (response) {
                  setFegGateways(response.fegGateways);
                  setFegGatewaysHealthStatus(response.fegGatewaysHealthStatus);
                  setActiveFegGatewayId(response.activeFegGatewayId);
                }
              },
            );
          }
        },
        health: fegGatewaysHealthStatus,
        activeFegGatewayId,
      }}>
      {props.children}
    </FEGGatewayContext.Provider>
  );
}

export default FEGGatewayContext;
