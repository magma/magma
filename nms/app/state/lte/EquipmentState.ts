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
  Enodeb,
  EnodebState as EnodebStateResponse,
  GatewayCellularConfigs,
  GatewayDnsConfigs,
  GatewayEpcConfigs,
  GatewayRanConfigs,
  GenericCommandParams,
  LteGateway,
  MagmadGatewayConfigs,
  MutableCellularGatewayPool,
  MutableLteGateway,
  PingRequest,
  Tier,
} from '../../../generated-ts';
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {
  EnodebSerial,
  GatewayId,
  GatewayPoolId,
  NetworkId,
  TierId,
} from '../../../shared/types/network';
import type {EnodebState} from '../../components/context/EnodebContext';
import type {
  GatewayPoolRecordsType,
  gatewayPoolsStateType,
} from '../../components/context/GatewayPoolsContext';
import type {OptionsObject} from 'notistack';

import MagmaAPI from '../../../api/MagmaAPI';

/************************** Gateway Tier State *******************************/
type InitTierStateProps = {
  networkId: NetworkId;
  setTiers: (arg0: Record<string, Tier>) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined;
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
  // reduce function gives a flow lint, hence using forEach instead
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
  setTiers: (arg0: Record<string, Tier>) => void;
  key: TierId;
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
  ) => (string | number) | null | undefined;
};

export async function FetchEnodebs(
  props: FetchProps,
): Promise<Record<string, EnodebInfo> | undefined> {
  const {networkId, id} = props;
  let enbs: Record<string, Enodeb> = {};
  if (id !== undefined && id !== null) {
    let enb: Enodeb;
    try {
      enb = (
        await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialGet({
          networkId: networkId,
          enodebSerial: id,
        })
      ).data;
      if (enb) {
        const newEnbSt = (
          await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialStateGet({
            networkId: networkId,
            enodebSerial: id,
          })
        ).data;
        const newEnb = {
          [id]: {
            enb_state: newEnbSt,
            enb: enb,
          },
        };
        return newEnb;
      }
    } catch (e) {
      return {
        [id]: {
          enb_state: {},
          enb: enb!,
        } as EnodebInfo,
      };
    }
  } else {
    const resp = (
      await MagmaAPI.enodebs.lteNetworkIdEnodebsGet({
        networkId,
      })
    ).data;
    enbs = resp.enodebs;
    if (!enbs) {
      return;
    }

    const requests = Object.keys(enbs).map(async k => {
      try {
        const {serial} = enbs[k];
        // eslint-disable-next-line max-len
        const enbSt = (
          await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialStateGet({
            networkId: networkId,
            enodebSerial: serial,
          })
        ).data;
        return [enbs[k], enbSt ?? {}] as const;
      } catch (e) {
        return [enbs[k], {}] as const;
      }
    });

    const enbResp = await Promise.all(requests);
    const enbInfo: Record<string, EnodebInfo> = {};
    enbResp.filter(Boolean).forEach(r => {
      if (r.length > 0) {
        const [enb, enbSt] = r;
        if (enb != null && enbSt != null) {
          enbInfo[enb.serial] = {
            enb,
            enb_state: enbSt as EnodebStateResponse,
          };
        }
      }
    });
    return enbInfo;
  }
}

type InitEnodeStateProps = {
  networkId: NetworkId;
  setEnbInfo: (arg0: Record<string, EnodebInfo>) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined;
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
  networkId: NetworkId;
  enbInfo: Record<string, EnodebInfo>;
  setEnbInfo: (arg0: Record<string, EnodebInfo>) => void;
  key: string;
  value?: EnodebInfo;
  newState?: EnodebState;
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
      await MagmaAPI.enodebs.lteNetworkIdEnodebsPost({
        networkId: networkId,
        enodeb: value.enb,
      });
      setEnbInfo({...enbInfo, [key]: value});
    } else {
      await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut({
        networkId: networkId,
        enodebSerial: key,
        enodeb: value.enb,
      });
      const prevEnbSt = enbInfo[key].enb_state;
      setEnbInfo({
        ...enbInfo,
        [key]: {
          enb_state: prevEnbSt,
          enb: value.enb,
        },
      });
    }
  } else {
    await MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialDelete({
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

type GatewayStateProps = {
  networkId: NetworkId;
  lteGateways: Record<string, LteGateway>;
  setLteGateways: (arg0: Record<string, LteGateway>) => void;
  key: GatewayId;
  value?: MutableLteGateway;
  newState?: Record<string, LteGateway>;
};

export async function SetGatewayState(props: GatewayStateProps) {
  const {networkId, lteGateways, setLteGateways, key, value, newState} = props;
  if (newState) {
    setLteGateways(newState);
    return;
  }
  if (value != null) {
    if (!(key in lteGateways)) {
      await MagmaAPI.lteGateways.lteNetworkIdGatewaysPost({
        networkId: networkId,
        gateway: value,
      });
      setLteGateways({...lteGateways, [key]: value});
    } else {
      await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdPut({
        networkId: networkId,
        gatewayId: key,
        gateway: value,
      });
      setLteGateways({...lteGateways, [key]: value});
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

export type UpdateGatewayProps = {
  gatewayId: GatewayId;
  tierId?: TierId;
  magmadConfigs?: MagmadGatewayConfigs;
  epcConfigs?: GatewayEpcConfigs;
  ranConfigs?: GatewayRanConfigs;
  dnsConfig?: GatewayDnsConfigs;
  cellularConfigs?: GatewayCellularConfigs;
  enbs?: Array<EnodebSerial>;
  networkId: NetworkId;
  setLteGateways: (arg0: Record<string, LteGateway>) => void;
};

export async function UpdateGateway(props: UpdateGatewayProps) {
  const {networkId, gatewayId, setLteGateways} = props;
  if (props.tierId !== undefined && props.tierId !== '') {
    await MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdTierPut({
      networkId,
      gatewayId: gatewayId,
      tierId: JSON.stringify(`"${props.tierId}"`),
    });
  }
  const requests = [];
  if (props.magmadConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdMagmadPut({
        networkId,
        gatewayId: gatewayId,
        magmad: props.magmadConfigs,
      }),
    );
  }
  if (props.epcConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularEpcPut({
        networkId,
        gatewayId: gatewayId,
        config: props.epcConfigs,
      }),
    );
  }
  if (props.ranConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularRanPut({
        networkId,
        gatewayId: gatewayId,
        config: props.ranConfigs,
      }),
    );
  }

  if (props.enbs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdConnectedEnodebSerialsPut(
        {
          networkId,
          gatewayId: gatewayId,
          serials: props.enbs,
        },
      ),
    );
  }

  if (props.dnsConfig) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularDnsPut({
        networkId,
        gatewayId: gatewayId,
        config: props.dnsConfig,
      }),
    );
  }

  if (props.cellularConfigs) {
    requests.push(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPut({
        networkId,
        gatewayId: gatewayId,
        config: props.cellularConfigs,
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
  setGatewayPools: (arg0: Record<string, gatewayPoolsStateType>) => void;
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
    const deletedGateways = Array.from(
      gatewayPools[key].gatewayPool.gateway_ids,
    ).filter(gwId => !resourcesIds.includes(gwId));
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
  setGatewayPools: (arg0: Record<string, gatewayPoolsStateType>) => void;
  networkId: NetworkId;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined;
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
        const records = Array.from(pool.gateway_ids).map(async id => {
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
