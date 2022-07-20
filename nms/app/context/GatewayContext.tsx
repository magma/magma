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
import {EnodebSerial, GatewayId, NetworkId} from '../../shared/types/network';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {
  GatewayCellularConfigs,
  GatewayDnsConfigs,
  GatewayEpcConfigs,
  GatewayRanConfigs,
  LteGateway,
  MagmadGatewayConfigs,
  MutableLteGateway,
} from '../../generated';
import {useEffect, useState} from 'react';

export type GatewayContextType = {
  state: Record<string, LteGateway>;
  setState: (key: GatewayId, val?: MutableLteGateway) => Promise<void>;
  updateGateway: (props: Partial<UpdateGatewayParams>) => Promise<void>;
  refetch: (id?: string) => void;
};
type GatewayContextProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

export type UpdateGatewayParams = {
  gatewayId: GatewayId;
  tierId?: string;
  magmadConfigs?: MagmadGatewayConfigs;
  epcConfigs?: GatewayEpcConfigs;
  ranConfigs?: GatewayRanConfigs;
  dnsConfig?: GatewayDnsConfigs;
  cellularConfigs?: GatewayCellularConfigs;
  enbs?: Array<EnodebSerial>;
  networkId: NetworkId;
  setLteGateways: (lteGateways: Record<string, LteGateway>) => void;
};

const GatewayContext = React.createContext<GatewayContextType>(
  {} as GatewayContextType,
);

export async function fetchGateways(params: {
  networkId: string;
  id?: string;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {networkId, id, enqueueSnackbar} = params;
  if (id !== undefined && id !== null) {
    try {
      const gateway = (
        await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdGet({
          networkId: networkId,
          gatewayId: id,
        })
      ).data;
      if (gateway) {
        return {
          [id]: gateway,
        };
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway information', {
        variant: 'error',
      });
    }
  } else {
    try {
      return (
        await MagmaAPI.lteGateways.lteNetworkIdGatewaysGet({
          networkId: networkId,
        })
      ).data;
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway information', {
        variant: 'error',
      });
    }
  }
}

async function setGatewayState(params: {
  networkId: NetworkId;
  lteGateways: Record<string, LteGateway>;
  setLteGateways: (lteGateways: Record<string, LteGateway>) => void;
  key: GatewayId;
  value?: MutableLteGateway;
}) {
  const {networkId, lteGateways, setLteGateways, key, value} = params;

  if (value != null) {
    if (!(key in lteGateways)) {
      await MagmaAPI.lteGateways.lteNetworkIdGatewaysPost({
        networkId: networkId,
        gateway: value,
      });
      // TODO[TS-migration] does it make sense that value is of type MutableLteGateway?
      setLteGateways({...lteGateways, [key]: value as LteGateway});
    } else {
      await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdPut({
        networkId: networkId,
        gatewayId: key,
        gateway: value,
      });
      setLteGateways({...lteGateways, [key]: value as LteGateway});
    }
  } else {
    await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdDelete({
      networkId: networkId,
      gatewayId: key,
    });
    const newLteGateways = {...lteGateways};
    delete newLteGateways[key];
    setLteGateways(newLteGateways);
  }
}

async function updateGateway(params: UpdateGatewayParams) {
  const {networkId, gatewayId, setLteGateways} = params;
  if (params.tierId !== undefined && params.tierId !== '') {
    await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdTierPut({
      networkId,
      gatewayId: gatewayId,
      tierId: JSON.stringify(`"${params.tierId}"`),
    });
  }
  const requests = [];
  if (params.magmadConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdMagmadPut({
        networkId,
        gatewayId: gatewayId,
        magmad: params.magmadConfigs,
      }),
    );
  }
  if (params.epcConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularEpcPut({
        networkId,
        gatewayId: gatewayId,
        config: params.epcConfigs,
      }),
    );
  }
  if (params.ranConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularRanPut({
        networkId,
        gatewayId: gatewayId,
        config: params.ranConfigs,
      }),
    );
  }

  if (params.enbs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdConnectedEnodebSerialsPut(
        {
          networkId,
          gatewayId: gatewayId,
          serials: params.enbs,
        },
      ),
    );
  }

  if (params.dnsConfig) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularDnsPut({
        networkId,
        gatewayId: gatewayId,
        config: params.dnsConfig,
      }),
    );
  }

  if (params.cellularConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPut({
        networkId,
        gatewayId: gatewayId,
        config: params.cellularConfigs,
      }),
    );
  }
  await Promise.all(requests);

  const gateways = (
    await MagmaAPI.lteGateways.lteNetworkIdGatewaysGet({
      networkId,
    })
  ).data;
  setLteGateways(gateways);
}

export function GatewayContextProvider(props: GatewayContextProps) {
  const {networkId} = props;
  const [lteGateways, setLteGateways] = useState<Record<string, LteGateway>>(
    {},
  );
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        const lteGateways = (
          await MagmaAPI.lteGateways.lteNetworkIdGatewaysGet({
            networkId,
          })
        ).data;
        setLteGateways(lteGateways);
      } catch (e) {
        enqueueSnackbar?.('failed fetching gateway information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <GatewayContext.Provider
      value={{
        state: lteGateways,
        setState: (key, value?) => {
          return setGatewayState({
            lteGateways,
            setLteGateways,
            networkId,
            key,
            value,
          });
        },
        updateGateway: props =>
          updateGateway({
            networkId,
            setLteGateways,
            ...props,
          } as UpdateGatewayParams),
        refetch: id => {
          void fetchGateways({
            id: id,
            networkId,
            enqueueSnackbar,
          }).then(gateways => {
            if (gateways) {
              setLteGateways(gatewayState =>
                id ? {...gatewayState, ...gateways} : gateways,
              );
            }
          });
        },
      }}>
      {props.children}
    </GatewayContext.Provider>
  );
}

export default GatewayContext;
