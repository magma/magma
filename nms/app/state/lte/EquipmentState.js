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
// $FlowFixMe migrated to typescript
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
// $FlowFixMe migrated to typescript
import type {EnodebState} from '../../components/context/EnodebContext';
import type {EnqueueSnackbarOptions} from 'notistack';
import type {
  GatewayPoolRecordsType,
  gatewayPoolsStateType,
  // $FlowFixMe migrated to typescript
} from '../../components/context/GatewayPoolsContext';
import type {
  enodeb_serials,
  gateway_cellular_configs,
  gateway_dns_configs,
  gateway_epc_configs,
  gateway_id,
  gateway_pool_id,
  gateway_ran_configs,
  generic_command_params,
  lte_gateway,
  magmad_gateway_configs,
  mutable_cellular_gateway_pool,
  mutable_lte_gateway,
  network_id,
  ping_request,
  tier,
  tier_id,
} from '../../../generated/MagmaAPIBindings';

import MagmaV1API from '../../../generated/WebClient';

/************************** Gateway Tier State *******************************/
type InitTierStateProps = {
  networkId: network_id,
  setTiers: ({[string]: tier}) => void,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function InitTierState(props: InitTierStateProps) {
  const {networkId, setTiers, enqueueSnackbar} = props;
  let tierIdList = [];
  try {
    tierIdList = await MagmaV1API.getNetworksByNetworkIdTiers({networkId});
  } catch (e) {
    enqueueSnackbar?.('failed fetching tier information', {variant: 'error'});
  }

  const requests = tierIdList.map(tierId => {
    try {
      return MagmaV1API.getNetworksByNetworkIdTiersByTierId({
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
  const tiers = {};
  // reduce function gives a flow lint, hence using forEach instead
  tierResponse.filter(Boolean).forEach(item => {
    tiers[item.id] = item;
  });
  setTiers(tiers);
}

type TierStateProps = {
  networkId: network_id,
  tiers: {[string]: tier},
  setTiers: ({[string]: tier}) => void,
  key: tier_id,
  value?: tier,
};

export async function SetTierState(props: TierStateProps) {
  const {networkId, tiers, setTiers, key, value} = props;
  if (value != null) {
    if (!(key in tiers)) {
      await MagmaV1API.postNetworksByNetworkIdTiers({
        networkId: networkId,
        tier: value,
      });
    } else {
      await MagmaV1API.putNetworksByNetworkIdTiersByTierId({
        networkId: networkId,
        tierId: key,
        tier: value,
      });
    }
    setTiers({...tiers, [key]: value});
  } else {
    await MagmaV1API.deleteNetworksByNetworkIdTiersByTierId({
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
  networkId: string,
  id?: string,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function FetchEnodebs(props: FetchProps) {
  const {networkId, id} = props;
  let enb = {};
  if (id !== undefined && id !== null) {
    try {
      enb = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerial({
        networkId: networkId,
        enodebSerial: id,
      });
      if (enb) {
        const newEnbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
          {
            networkId: networkId,
            enodebSerial: id,
          },
        );
        const newEnb = {[id]: {enb_state: newEnbSt, enb: enb}};
        return newEnb;
      }
    } catch (e) {
      return {[id]: {enb_state: {}, enb: enb}};
    }
  } else {
    let resp = {};
    resp = await MagmaV1API.getLteByNetworkIdEnodebs({networkId});
    enb = resp['enodebs'];
    if (!enb) {
      return;
    }

    const requests = Object.keys(enb).map(async k => {
      try {
        const {serial} = enb[k];
        // eslint-disable-next-line max-len
        const enbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
          {
            networkId: networkId,
            enodebSerial: serial,
          },
        );
        return [enb[k], enbSt ?? {}];
      } catch (e) {
        return [enb[k], {}];
      }
    });

    const enbResp = await Promise.all(requests);
    const enbInfo = {};
    enbResp.filter(Boolean).forEach(r => {
      if (r.length > 0) {
        const [enb, enbSt] = r;
        if (enb != null && enbSt != null) {
          enbInfo[enb.serial] = {
            enb: enb,
            enb_state: enbSt,
          };
        }
      }
    });
    return enbInfo;
  }
}

type InitEnodeStateProps = {
  networkId: network_id,
  setEnbInfo: ({[string]: EnodebInfo}) => void,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function InitEnodeState(props: InitEnodeStateProps) {
  const enodebInfo = await FetchEnodebs({
    networkId: props.networkId,
    enqueueSnackbar: props.enqueueSnackbar,
  });
  if (enodebInfo) {
    props.setEnbInfo(enodebInfo);
  }
}

type EnodebStateProps = {
  networkId: network_id,
  enbInfo: {[string]: EnodebInfo},
  setEnbInfo: ({[string]: EnodebInfo}) => void,
  key: string,
  value?: EnodebInfo,
  newState?: EnodebState,
};

export async function SetEnodebState(props: EnodebStateProps) {
  const {networkId, enbInfo, setEnbInfo, key, value, newState} = props;
  if (newState) {
    setEnbInfo(newState.enbInfo);
    return;
  }
  if (value != null) {
    // remove attached gateway id read only property
    if (value.enb.hasOwnProperty('attached_gateway_id')) {
      delete value.enb['attached_gateway_id'];
    }
    if (!(key in enbInfo)) {
      await MagmaV1API.postLteByNetworkIdEnodebs({
        networkId: networkId,
        enodeb: value.enb,
      });
      setEnbInfo({...enbInfo, [key]: value});
    } else {
      await MagmaV1API.putLteByNetworkIdEnodebsByEnodebSerial({
        networkId: networkId,
        enodebSerial: key,
        enodeb: value.enb,
      });
      const prevEnbSt = enbInfo[key].enb_state;
      setEnbInfo({...enbInfo, [key]: {enb_state: prevEnbSt, enb: value.enb}});
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdEnodebsByEnodebSerial({
      networkId: networkId,
      enodebSerial: key,
    });
    const newEnbInfo = {...enbInfo};
    delete newEnbInfo[key];
    setEnbInfo(newEnbInfo);
    return;
  }
}

/**************************** Gateway State **********************************/

export async function FetchGateways(props: FetchProps) {
  const {networkId, id, enqueueSnackbar} = props;
  if (id !== undefined && id !== null) {
    try {
      const gateway = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayId({
        networkId: networkId,
        gatewayId: id,
      });
      if (gateway) {
        return {[id]: gateway};
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway information', {
        variant: 'error',
      });
    }
  } else {
    try {
      return await MagmaV1API.getLteByNetworkIdGateways({
        networkId: networkId,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway information', {
        variant: 'error',
      });
    }
  }
}

type GatewayStateProps = {
  networkId: network_id,
  lteGateways: {[string]: lte_gateway},
  setLteGateways: ({[string]: lte_gateway}) => void,
  key: gateway_id,
  value?: mutable_lte_gateway,
  newState?: {[string]: lte_gateway},
};

export async function SetGatewayState(props: GatewayStateProps) {
  const {networkId, lteGateways, setLteGateways, key, value, newState} = props;
  if (newState) {
    setLteGateways(newState);
    return;
  }
  if (value != null) {
    if (!(key in lteGateways)) {
      await MagmaV1API.postLteByNetworkIdGateways({
        networkId: networkId,
        gateway: value,
      });
      setLteGateways({...lteGateways, [key]: value});
    } else {
      await MagmaV1API.putLteByNetworkIdGatewaysByGatewayId({
        networkId: networkId,
        gatewayId: key,
        gateway: value,
      });
      setLteGateways({...lteGateways, [key]: value});
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdGatewaysByGatewayId({
      networkId: networkId,
      gatewayId: key,
    });
    const newLteGateways = {...lteGateways};
    delete newLteGateways[key];
    setLteGateways(newLteGateways);
  }
}

export type UpdateGatewayProps = {
  gatewayId: gateway_id,
  tierId?: tier_id,
  magmadConfigs?: magmad_gateway_configs,
  epcConfigs?: gateway_epc_configs,
  ranConfigs?: gateway_ran_configs,
  dnsConfig?: gateway_dns_configs,
  cellularConfigs?: gateway_cellular_configs,
  enbs?: enodeb_serials,
  networkId: network_id,
  setLteGateways: ({[string]: lte_gateway}) => void,
};

export async function UpdateGateway(props: UpdateGatewayProps) {
  const {networkId, gatewayId, setLteGateways} = props;
  if (props.tierId !== undefined && props.tierId !== '') {
    await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdTier({
      networkId,
      gatewayId: gatewayId,
      tierId: JSON.stringify(`"${props.tierId}"`),
    });
  }
  const requests = [];
  if (props.magmadConfigs) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdMagmad({
        networkId,
        gatewayId: gatewayId,
        magmad: props.magmadConfigs,
      }),
    );
  }
  if (props.epcConfigs) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellularEpc({
        networkId,
        gatewayId: gatewayId,
        config: props.epcConfigs,
      }),
    );
  }
  if (props.ranConfigs) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellularRan({
        networkId,
        gatewayId: gatewayId,
        config: props.ranConfigs,
      }),
    );
  }

  if (props.enbs) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdConnectedEnodebSerials({
        networkId,
        gatewayId: gatewayId,
        serials: props.enbs,
      }),
    );
  }

  if (props.dnsConfig) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellularDns({
        networkId,
        gatewayId: gatewayId,
        config: props.dnsConfig,
      }),
    );
  }

  if (props.cellularConfigs) {
    requests.push(
      MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellular({
        networkId,
        gatewayId: gatewayId,
        config: props.cellularConfigs,
      }),
    );
  }
  await Promise.all(requests);
  const gateways = await MagmaV1API.getLteByNetworkIdGateways({
    networkId,
  });
  setLteGateways(gateways);
}
export type GatewayCommandProps = {
  networkId: network_id,
  gatewayId: gateway_id,
  command: 'reboot' | 'ping' | 'restartServices' | 'generic',
  pingRequest?: ping_request,
  params?: generic_command_params,
};

export async function RunGatewayCommands(props: GatewayCommandProps) {
  const {networkId, gatewayId} = props;

  switch (props.command) {
    case 'reboot':
      return await MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot(
        {networkId, gatewayId},
      );

    case 'restartServices':
      return await MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
        {networkId, gatewayId, services: []},
      );

    case 'ping':
      if (props.pingRequest != null) {
        return await MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandPing(
          {networkId, gatewayId, pingRequest: props.pingRequest},
        );
      }

    default:
      if (props.params != null) {
        return await MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandGeneric(
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
      const gatewayPool = await MagmaV1API.getLteByNetworkIdGatewayPoolsByGatewayPoolId(
        {
          networkId: networkId,
          gatewayPoolId: id,
        },
      );
      return gatewayPool;
    } catch (e) {
      enqueueSnackbar?.(`failed fetching gateway pool ${id} information`, {
        variant: 'error',
      });
    }
  } else {
    try {
      return await MagmaV1API.getLteByNetworkIdGatewayPools({
        networkId: networkId,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching gateway pools information', {
        variant: 'error',
      });
    }
  }
}
type GatewayPoolsStateProps = {
  networkId: network_id,
  gatewayPools: {[string]: gatewayPoolsStateType},
  setGatewayPools: ({[string]: gatewayPoolsStateType}) => void,
  key: gateway_pool_id,
  value?: mutable_cellular_gateway_pool,
  resources?: Array<GatewayPoolRecordsType>,
};
// update gateway pool config
export async function SetGatewayPoolsState(props: GatewayPoolsStateProps) {
  const {networkId, gatewayPools, setGatewayPools, key, value} = props;
  if (value != null) {
    if (!(key in gatewayPools)) {
      await MagmaV1API.postLteByNetworkIdGatewayPools({
        networkId: networkId,
        haGatewayPool: value,
      });
      setGatewayPools({
        ...gatewayPools,
        [key]: {
          gatewayPool: {...value, gateway_ids: []},
          gatewayPoolRecords: [],
        },
      });
    } else {
      await MagmaV1API.putLteByNetworkIdGatewayPoolsByGatewayPoolId({
        networkId,
        gatewayPoolId: key,
        haGatewayPool: value,
      });
      const newGwPool = await FetchGatewayPools({networkId, id: key});
      setGatewayPools({
        ...gatewayPools,
        [key]: {
          gatewayPool: newGwPool,
          gatewayPoolRecords: gatewayPools[key].gatewayPoolRecords,
        },
      });
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdGatewayPoolsByGatewayPoolId({
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
        return await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellularPooling(
          {
            networkId: networkId,
            gatewayId: gateway_id,
            resource: [gatewayConfig] || [],
          },
        );
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
        await MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdCellularPooling({
          networkId: networkId,
          gatewayId: gwId,
          resource: [],
        }),
    );
    await Promise.all(deleteRequests);
    const newGwPool = await FetchGatewayPools({networkId: networkId, id: key});
    setGatewayPools({
      ...gatewayPools,
      [key]: {gatewayPool: newGwPool, gatewayPoolRecords: resources},
    });
    return;
  }
}
type InitGatewayPoolStateType = {
  setGatewayPools: ({[string]: gatewayPoolsStateType}) => void,
  networkId: network_id,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function InitGatewayPoolState(props: InitGatewayPoolStateType) {
  const {networkId, setGatewayPools, enqueueSnackbar} = props;
  const pools = await FetchGatewayPools({networkId: networkId});

  if (pools) {
    const poolGatewayState = {};
    Object.keys(pools).map(async poolId => {
      const pool = pools[poolId];
      try {
        // get primary/secondary gateways for each gateway pool
        const records = pool.gateway_ids?.map(async id => {
          const gatewayRecords = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayIdCellularPooling(
            {
              networkId,
              gatewayId: id,
            },
          );
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
