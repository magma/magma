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
import type {
  gateway_he_config,
  gateway_id,
  lte_gateway,
  mutable_cellular_gateway_pool,
  network_dns_config,
  network_id,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

export const toString = (input: ?number | ?string): string => {
  return input !== null && input !== undefined ? input + '' : '';
};

type GatewaySharedFields = {
  hardware_id: string,
  name: string,
  logicalID: string,
  challengeType: string,
  enodebRFTXEnabled: boolean,
  enodebRFTXOn: boolean,
  latLon: {lat: number, lon: number},
  version: string,
  vpnIP: string,
  enodebConnected: boolean,
  gpsConnected: boolean,
  isBackhaulDown: boolean,
  lastCheckin: string,
  mmeConnected: boolean,
  autoupgradePollInterval: ?number,
  checkinInterval: ?number,
  checkinTimeout: ?number,
  tier: ?string,
  autoupgradeEnabled: boolean,
  attachedEnodebSerials: Array<string>,
  ran: {pci: ?number, transmitEnabled: boolean},
  epc: {ipBlock: string, natEnabled: boolean},
  nonEPSService: {
    control: number,
    csfbRAT: number,
    csfbMCC: ?string,
    csfbMNC: ?string,
    lac: ?number,
  },
};

export type GatewayV1 = {
  ...GatewaySharedFields,
  rawGateway: lte_gateway,
};

export type GatewayPayload = {
  gateway_id: GatewayId,
  status?: GatewayStatusPayload,
  record?: AccessGatewayRecord,
  name?: GatewayName,
};

type SystemStatus = {
  time?: number,
  uptime_secs?: number,
  cpu_user?: number,
  cpu_system?: number,
  cpu_idle?: number,
  mem_total?: number,
  mem_available?: number,
  mem_used?: number,
  mem_free?: number,
  swap_total?: number,
  swap_used?: number,
  swap_free?: number,
  disk_partitions?: Array<DiskPartition>,
};

type PlatformInfo = {
  vpn_ip?: string,
  packages?: Array<SoftwarePackage>,
  kernel_version?: string,
  kernel_versions_installed?: Array<string>,
  config_info?: ConfigInfo,
};

type MachineInfo = {
  cpu_info?: {
    core_count?: number,
    threads_per_core?: number,
    architecture?: string,
    model_name?: string,
  },
  network_info?: {
    network_interfaces?: Array<NetworkInterface>,
    routing_table?: Array<Route>,
  },
};

type NetworkInterface = {
  network_interface_id?: string,
  status?: 'UP' | 'DOWN' | 'UNKNOWN',
  mac_address?: string,
  ip_addresses?: Array<string>,
  ipv6_addresses?: Array<string>,
};

type DiskPartition = {
  device?: string,
  mount_point?: string,
  total?: number,
  used?: number,
  free?: number,
};

type SoftwarePackage = {
  name?: string,
  version?: string,
};

type ConfigInfo = {
  mconfig_created_at?: number,
};

type Route = {
  destination_ip?: string,
  gateway_ip?: string,
  genmask?: string,
  network_interface_id?: string,
};

type GatewayName = string;

type ChallengeKey = {
  key_type: 'ECHO' | 'SOFTWARE_ECDSA_SHA256',
  key?: string,
};

type AccessGatewayRecord = {hardware_id: string, key: ChallengeKey};

type GatewayId = string;

// TODO: strip out devmand related fields and put them into a separate file
type GatewayMeta = {
  gps_latitude: number,
  gps_longitude: number,
  rf_tx_on: boolean,
  enodeb_connected: number,
  gps_connected: number,
  mme_connected: number,
  devmand: ?string,
  status: ?string,
};

type GatewayStatusPayload = {
  checkin_time?: number,
  hardware_id?: string,
  version?: string,
  system_status?: SystemStatus,
  platform_info?: PlatformInfo,
  machine_info?: MachineInfo,
  cert_expiration_time?: number,
  meta?: GatewayMeta,
  vpn_ip?: string,
  kernel_version?: string,
  kernel_versions_installed?: Array<string>,
};

export type FederationGatewayHealthStatus = {
  status: string,
};

const GATEWAY_KEEPALIVE_TIMEOUT_MS = 1000 * 5 * 60;

export const GatewayTypeEnum = Object.freeze({
  HEALTHY_GATEWAY: 'Good',
  UNHEALTHY_GATEWAY: 'Bad',
  UNKNOWN: '-',
});

// health status used for federation gateways
export const HEALTHY_STATUS = 'HEALTHY';

export const UNHEALTHY_STATUS = 'UNHEALTHY';

export default function isGatewayHealthy({status}: lte_gateway) {
  if (status != null) {
    const checkin = status.checkin_time;
    if (checkin != null) {
      return Date.now() - checkin < GATEWAY_KEEPALIVE_TIMEOUT_MS;
    }
  }
  return false;
}

/**
 * Returns health status of the federation gateway.
 *
 * @param {network_id} networkId: Id of the federation network.
 * @param {gateway_id} gatewayId: Id of the gateway
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 * @returns the health status of the gateway or an empty string for an error.
 */
export async function getFederationGatewayHealthStatus(
  networkId: network_id,
  gatewayId: gateway_id,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
): Promise<FederationGatewayHealthStatus> {
  try {
    const gwHealthStatus = await MagmaV1API.getFegByNetworkIdGatewaysByGatewayIdHealthStatus(
      {
        networkId,
        gatewayId,
      },
    );
    return {status: gwHealthStatus.status};
  } catch (e) {
    enqueueSnackbar?.(
      'failed fetching health status information for federation gateway with id ' +
        gatewayId,
      {
        variant: 'error',
      },
    );
    return {status: ''};
  }
}

export const DynamicServices = Object.freeze({
  MONITORD: 'monitord',
  EVENTD: 'eventd',
  TD_AGENT_BIT: 'td-agent-bit',
});

export const DEFAULT_GATEWAY_CONFIG: lte_gateway = {
  apn_resources: {},
  cellular: {
    epc: {
      ip_block: '192.168.128.0/24',
      nat_enabled: true,
      dns_primary: '',
      dns_secondary: '',
      sgi_management_iface_gw: '',
      sgi_management_iface_static_ip: '',
      sgi_management_iface_vlan: '',
    },
    ran: {
      pci: 260,
      transmit_enabled: true,
    },
  },
  connected_enodeb_serials: [],
  description: '',
  device: {
    hardware_id: '',
    key: {
      key: '',
      key_type: 'SOFTWARE_ECDSA_SHA256',
    },
  },
  id: '',
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 60,
    checkin_interval: 60,
    checkin_timeout: 30,
    dynamic_services: [DynamicServices.EVENTD, DynamicServices.TD_AGENT_BIT],
  },
  name: '',
  status: {
    platform_info: {
      packages: [
        {
          version: '',
        },
      ],
    },
  },
  tier: 'default',
};

export const DEFAULT_DNS_CONFIG: network_dns_config = {
  dhcp_server_enabled: false,
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

export const DEFAULT_HE_CONFIG: gateway_he_config = {
  enable_encryption: false,
  encryption_key: '',
  enable_header_enrichment: false,
  he_encoding_type: 'BASE64',
  he_encryption_algorithm: 'RC4',
  he_hash_function: 'MD5',
};

export const DEFAULT_GW_POOL_CONFIG: mutable_cellular_gateway_pool = {
  config: {mme_group_id: 1},
  gateway_pool_id: '',
  gateway_pool_name: '',
};

export const DEFAULT_GW_PRIMARY_CONFIG = {
  gateway_id: '',
  gateway_pool_id: '',
  mme_code: 1,
  mme_relative_capacity: 255,
};

export const DEFAULT_GW_SECONDARY_CONFIG = {
  gateway_id: '',
  gateway_pool_id: '',
  mme_code: 1,
  mme_relative_capacity: 1,
};
// services running on the LTE AGWq
export const RUNNING_SERVICES = [
  'policydb',
  'control_proxy',
  'mobilityd',
  'smsd',
  'pipelined',
  'sessiond',
  'redis',
  'dnsd',
  'mme',
  'directoryd',
  'eventd',
  'enodebd',
  'state',
  'subscriberdb',
  'magmad',
  'health',
  'ctraced',
];
