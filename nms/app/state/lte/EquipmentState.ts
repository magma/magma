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
  CellularGatewayPool,
  GenericCommandParams,
  MutableCellularGatewayPool,
  PingRequest,
  Tier,
} from '../../../generated';
import type {
  GatewayId,
  GatewayPoolId,
  NetworkId,
} from '../../../shared/types/network';
import type {
  GatewayPoolRecordsType,
  gatewayPoolsStateType,
} from '../../components/context/GatewayPoolsContext';
import type {OptionsObject} from 'notistack';

import MagmaAPI from '../../api/MagmaAPI';

/************************** Gateway Tier State *******************************/
type InitTierStateProps = {
  networkId: NetworkId;
  setTiers: (tiers: Record<string, Tier>) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

export async function InitTierState(props: InitTierStateProps) {
  const {networkId, setTiers, enqueueSnackbar} = props;
  let tierIdList: Array<string> = [];
  try {
    tierIdList = (
      await MagmaAPI.upgrades.networksNetworkIdTiersGet({
        networkId,
      })
    ).data;
  } catch (e) {
    enqueueSnackbar?.('failed fetching tier information', {
      variant: 'error',
    });
  }

  const requests = tierIdList.map(tierId => {
    try {
      return MagmaAPI.upgrades.networksNetworkIdTiersTierIdGet({
        networkId,
        tierId,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching tier information for ' + tierId, {
        variant: 'error',
      });
      return;
    }
  });

  const tierResponse = await Promise.all(requests);
  const tiers: Record<string, Tier> = {};
  tierResponse
    .filter(Boolean)
    .map(res => res!.data)
    .forEach(item => {
      tiers[item.id] = item;
    });
  setTiers(tiers);
}

type TierStateProps = {
  networkId: NetworkId;
  tiers: Record<string, Tier>;
  setTiers: (tiers: Record<string, Tier>) => void;
  key: string;
  value?: Tier;
};

export async function SetTierState(props: TierStateProps) {
  const {networkId, tiers, setTiers, key, value} = props;

  if (value != null) {
    if (!(key in tiers)) {
      await MagmaAPI.upgrades.networksNetworkIdTiersPost({
        networkId: networkId,
        tier: value,
      });
    } else {
      await MagmaAPI.upgrades.networksNetworkIdTiersTierIdPut({
        networkId: networkId,
        tierId: key,
        tier: value,
      });
    }
    setTiers({...tiers, [key]: value});
  } else {
    await MagmaAPI.upgrades.networksNetworkIdTiersTierIdDelete({
      networkId: networkId,
      tierId: key,
    });
    const newTiers = {...tiers};
    delete newTiers[key];
    setTiers(newTiers);
  }
}

/**************************** Enode State ************************************/
type FetchProps = {
  networkId: string;
  id?: string;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

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

/**************************** Gateway Pools State **********************************/

export async function FetchGatewayPools(props: FetchProps) {
  const {networkId, id, enqueueSnackbar} = props;
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

type GatewayPoolsStateProps = {
  networkId: NetworkId;
  gatewayPools: Record<string, gatewayPoolsStateType>;
  setGatewayPools: (
    gatewayPools: Record<string, gatewayPoolsStateType>,
  ) => void;
  key: GatewayPoolId;
  value?: MutableCellularGatewayPool;
  resources?: Array<GatewayPoolRecordsType>;
};
// update gateway pool config
export async function SetGatewayPoolsState(props: GatewayPoolsStateProps) {
  const {networkId, gatewayPools, setGatewayPools, key, value} = props;
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
      const newGwPool = await FetchGatewayPools({
        networkId,
        id: key,
      });
      setGatewayPools({
        ...gatewayPools,
        [key]: {
          gatewayPool: newGwPool,
          gatewayPoolRecords: gatewayPools[key].gatewayPoolRecords,
        } as gatewayPoolsStateType,
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
export async function UpdateGatewayPoolRecords(props: GatewayPoolsStateProps) {
  const {networkId, gatewayPools, setGatewayPools, key, resources} = props;

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
    const newGwPool = await FetchGatewayPools({
      networkId: networkId,
      id: key,
    });
    setGatewayPools({
      ...gatewayPools,
      [key]: {
        gatewayPool: newGwPool,
        gatewayPoolRecords: resources,
      } as gatewayPoolsStateType,
    });
    return;
  }
}
type InitGatewayPoolStateType = {
  setGatewayPools: (
    gatewayPools: Record<string, gatewayPoolsStateType>,
  ) => void;
  networkId: NetworkId;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

export async function InitGatewayPoolState(props: InitGatewayPoolStateType) {
  const {networkId, setGatewayPools, enqueueSnackbar} = props;
  const pools = (await FetchGatewayPools({
    networkId: networkId,
  })) as Record<string, CellularGatewayPool>;

  if (pools) {
    const poolGatewayState: Record<string, gatewayPoolsStateType> = {};
    Object.keys(pools).map(async poolId => {
      const pool = pools[poolId];
      try {
        // get primary/secondary gateways for each gateway pool
        const records = pool.gateway_ids.map(async id => {
          const gatewayRecords = (
            await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPoolingGet(
              {
                networkId,
                gatewayId: id,
              },
            )
          ).data;
          return gatewayRecords.map(record => {
            return {...record, gateway_id: id};
          });
        });
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
    });
    setGatewayPools(poolGatewayState);
  }
}
