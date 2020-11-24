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

import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {
  enodeb_serials,
  gateway_dns_configs,
  gateway_epc_configs,
  gateway_id,
  gateway_ran_configs,
  generic_command_params,
  lte_gateway,
  magmad_gateway_configs,
  mutable_lte_gateway,
  network_id,
  ping_request,
  tier,
  tier_id,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

/************************** Gateway Tier State *******************************/
type InitTierStateProps = {
  networkId: network_id,
  setTiers: ({[string]: tier}) => void,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
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
type InitEnodeStateProps = {
  networkId: network_id,
  setEnbInfo: ({[string]: EnodebInfo}) => void,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
};

export async function InitEnodeState(props: InitEnodeStateProps) {
  const {networkId, setEnbInfo, enqueueSnackbar} = props;
  let enb = {};
  try {
    enb = await MagmaV1API.getLteByNetworkIdEnodebs({networkId});
  } catch (e) {
    enqueueSnackbar?.('failed fetching enodeb information', {
      variant: 'error',
    });
    return;
  }

  if (!enb) {
    return;
  }

  let err = false;
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
      err = true;
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
  if (err) {
    enqueueSnackbar?.('failed fetching enodeb state information', {
      variant: 'error',
    });
  }
  setEnbInfo(enbInfo);
}

type EnodebStateProps = {
  networkId: network_id,
  enbInfo: {[string]: EnodebInfo},
  setEnbInfo: ({[string]: EnodebInfo}) => void,
  key: string,
  value?: EnodebInfo,
};

export async function SetEnodebState(props: EnodebStateProps) {
  const {networkId, enbInfo, setEnbInfo, key, value} = props;
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
    const newEnb = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerial({
      networkId: networkId,
      enodebSerial: key,
    });
    if (newEnb) {
      const newEnbSt = await MagmaV1API.getLteByNetworkIdEnodebsByEnodebSerialState(
        {
          networkId: networkId,
          enodebSerial: key,
        },
      );
      setEnbInfo({...enbInfo, [key]: {enb_state: newEnbSt, enb: newEnb}});
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdEnodebsByEnodebSerial({
      networkId: networkId,
      enodebSerial: key,
    });
    const newEnbInfo = {...enbInfo};
    delete newEnbInfo[key];
    setEnbInfo(newEnbInfo);
  }
}

/**************************** Gateway State **********************************/
type GatewayStateProps = {
  networkId: network_id,
  lteGateways: {[string]: lte_gateway},
  setLteGateways: ({[string]: lte_gateway}) => void,
  key: gateway_id,
  value?: mutable_lte_gateway,
};

export async function SetGatewayState(props: GatewayStateProps) {
  const {networkId, lteGateways, setLteGateways, key, value} = props;
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
    const gateway = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayId({
      networkId: networkId,
      gatewayId: key,
    });
    if (gateway) {
      const newLteGateways = {...lteGateways, [key]: gateway};
      setLteGateways(newLteGateways);
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
  await Promise.all(requests);
  const gateways = await MagmaV1API.getLteByNetworkIdGateways({
    networkId,
  });
  setLteGateways(gateways);
}

export type GatewayCommandProps = {
  networkId: network_id,
  gatewayId: gateway_id,
  command: 'reboot' | 'ping' | 'generic',
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
