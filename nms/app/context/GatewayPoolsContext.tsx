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
import {
  CellularGatewayPool,
  CellularGatewayPoolRecord,
  MutableCellularGatewayPool,
} from '../../generated';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {GatewayPoolId, NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';

// add gateway ID to gateway pool records (gateway primary/secondary)
export type GatewayPoolRecordsType = {
  gateway_id: string;
} & CellularGatewayPoolRecord;
export type GatewayPoolsStateType = {
  gatewayPool: CellularGatewayPool;
  gatewayPoolRecords: Array<GatewayPoolRecordsType>;
};
/** GatewayPoolsContextType
 * state: gateway pool config and associated gateway pool records
 * setState: POST, PUT, DELETE gateway pool config
 * updateGatewayPoolRecords: POST, PUT, DELETE gateway pool records
 */
type GatewayPoolsContextType = {
  state: Record<string, GatewayPoolsStateType>;
  setState: (
    key: GatewayPoolId,
    val?: MutableCellularGatewayPool,
  ) => Promise<void>;
  updateGatewayPoolRecords: (
    key: GatewayPoolId,
    val?: MutableCellularGatewayPool,
    resources?: Array<GatewayPoolRecordsType>,
  ) => Promise<void>;
};
type GatewayToolProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

const GatewayPoolsContext = React.createContext<GatewayPoolsContextType>(
  {} as GatewayPoolsContextType,
);

async function initGatewayPoolState(params: {
  setGatewayPools: (
    gatewayPools: Record<string, GatewayPoolsStateType>,
  ) => void;
  networkId: NetworkId;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {networkId, setGatewayPools, enqueueSnackbar} = params;
  const pools = (await fetchGatewayPools({
    networkId: networkId,
  })) as Record<string, CellularGatewayPool>;

  if (pools) {
    const poolGatewayState: Record<string, GatewayPoolsStateType> = {};
    for (const poolId in pools) {
      const pool = pools[poolId];
      try {
        // get primary/secondary gateways for each gateway pool
        const records = pool.gateway_ids.map(id =>
          MagmaAPI.lteGateways
            .lteNetworkIdGatewaysGatewayIdCellularPoolingGet({
              networkId,
              gatewayId: id,
            })
            .then(({data}) => {
              return data.map(record => ({...record, gateway_id: id}));
            }),
        );
        const gwPoolRecords = await Promise.all(records);
        poolGatewayState[poolId] = {
          gatewayPool: pool,
          gatewayPoolRecords: gwPoolRecords.flat() || [],
        };
      } catch (error) {
        enqueueSnackbar?.('failed fetching gateway pool records', {
          variant: 'error',
        });
      }
    }
    setGatewayPools(poolGatewayState);
  }
}

async function fetchGatewayPools(params: {
  networkId: string;
  id?: string;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {networkId, id, enqueueSnackbar} = params;
  if (id !== undefined && id !== null) {
    try {
      const gatewayPool = (
        await MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdGet({
          networkId: networkId,
          gatewayPoolId: id,
        })
      ).data;
      return gatewayPool;
    } catch (e) {
      enqueueSnackbar?.(`failed fetching gateway pool ${id} information`, {
        variant: 'error',
      });
    }
  } else {
    try {
      return (
        await MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGet({
          networkId: networkId,
        })
      ).data;
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway pools information', {
        variant: 'error',
      });
    }
  }
}

type GatewayPoolsStateParams = {
  networkId: NetworkId;
  gatewayPools: Record<string, GatewayPoolsStateType>;
  setGatewayPools: (
    gatewayPools: Record<string, GatewayPoolsStateType>,
  ) => void;
  key: GatewayPoolId;
  value?: MutableCellularGatewayPool;
  resources?: Array<GatewayPoolRecordsType>;
};

// update gateway pool config
async function setGatewayPoolsState(params: GatewayPoolsStateParams) {
  const {networkId, gatewayPools, setGatewayPools, key, value} = params;
  if (value) {
    if (!(key in gatewayPools)) {
      await MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsPost({
        networkId: networkId,
        hAGatewayPool: value,
      });
      setGatewayPools({
        ...gatewayPools,
        [key]: {
          gatewayPool: {...value, gateway_ids: []},
          gatewayPoolRecords: [],
        },
      });
    } else {
      await MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdPut({
        networkId,
        gatewayPoolId: key,
        hAGatewayPool: value,
      });
      const newGwPool = await fetchGatewayPools({
        networkId,
        id: key,
      });
      setGatewayPools({
        ...gatewayPools,
        [key]: {
          gatewayPool: newGwPool,
          gatewayPoolRecords: gatewayPools[key].gatewayPoolRecords,
        } as GatewayPoolsStateType,
      });
    }
  } else {
    await MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdDelete({
      networkId: networkId,
      gatewayPoolId: key,
    });
    const newGatewayPools = {...gatewayPools};
    delete newGatewayPools[key];
    setGatewayPools(newGatewayPools);
  }
}

// update gateway pool primary/secondary gateways
async function updateGatewayPoolRecords(params: GatewayPoolsStateParams) {
  const {networkId, gatewayPools, setGatewayPools, key, resources} = params;

  // add primary/secondary gateways
  if (resources != null) {
    const requests = resources.map(async resource => {
      if (resource.gateway_id !== '') {
        const {gateway_id, ...gatewayConfig} = resource;
        gatewayConfig.gateway_pool_id = key;
        return (
          await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPoolingPut(
            {
              networkId: networkId,
              gatewayId: gateway_id,
              resource: [gatewayConfig] || [],
            },
          )
        ).data;
      }
    });
    await Promise.all(requests);

    // delete primary/secondary gateways
    const resourcesIds = resources.map(resource => resource.gateway_id);
    const deletedGateways = gatewayPools[key].gatewayPool.gateway_ids.filter(
      gwId => !resourcesIds.includes(gwId),
    );
    const deleteRequests = deletedGateways.map(
      async gwId =>
        await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPoolingPut(
          {
            networkId: networkId,
            gatewayId: gwId,
            resource: [],
          },
        ),
    );
    await Promise.all(deleteRequests);
    const newGwPool = await fetchGatewayPools({
      networkId: networkId,
      id: key,
    });
    setGatewayPools({
      ...gatewayPools,
      [key]: {
        gatewayPool: newGwPool,
        gatewayPoolRecords: resources,
      } as GatewayPoolsStateType,
    });
    return;
  }
}

export function GatewayPoolsContextProvider(props: GatewayToolProps) {
  const {networkId} = props;
  const [isLoading, setIsLoading] = useState(true);
  const [gatewayPools, setGatewayPools] = useState<
    Record<string, GatewayPoolsStateType>
  >({});
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        await initGatewayPoolState({
          enqueueSnackbar,
          networkId,
          setGatewayPools,
        });
      } catch (e) {
        enqueueSnackbar?.('failed fetching gateway pool information', {
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
    <GatewayPoolsContext.Provider
      value={{
        state: gatewayPools,
        setState: (key, value?) =>
          setGatewayPoolsState({
            gatewayPools,
            setGatewayPools,
            networkId,
            key,
            value,
          }),
        updateGatewayPoolRecords: (key, value?, resources?) =>
          updateGatewayPoolRecords({
            gatewayPools,
            setGatewayPools,
            networkId,
            key,
            value,
            resources,
          }),
      }}>
      {props.children}
    </GatewayPoolsContext.Provider>
  );
}

export default GatewayPoolsContext;
