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
import type {OptionsObject} from 'notistack';

import MagmaAPI from '../../api/MagmaAPI';
import type {
  GatewayHeConfig,
  LteGateway,
  MutableCellularGatewayPool,
  NetworkDnsConfig,
  ServiceStatusHealth,
} from '../../generated-ts';
import type {GatewayId, NetworkId} from '../../shared/types/network';

export const toString = (input: number | null | undefined): string => {
  return input !== null && input !== undefined ? input.toString() : '';
};

type GatewaySharedFields = {
  hardware_id: string;
  name: string;
  logicalID: string;
  challengeType: string;
  enodebRFTXEnabled: boolean;
  enodebRFTXOn: boolean;
  latLon: {
    lat: number;
    lon: number;
  };
  version: string;
  vpnIP: string;
  enodebConnected: boolean;
  gpsConnected: boolean;
  isBackhaulDown: boolean;
  lastCheckin: string;
  mmeConnected: boolean;
  autoupgradePollInterval: number | null | undefined;
  checkinInterval: number | null | undefined;
  checkinTimeout: number | null | undefined;
  tier: string | null | undefined;
  autoupgradeEnabled: boolean;
  attachedEnodebSerials: Array<string>;
  ran: {
    pci: number | null | undefined;
    transmitEnabled: boolean;
  };
  epc: {
    ipBlock: string;
    natEnabled: boolean;
  };
  nonEPSService: {
    control: number;
    csfbRAT: number;
    csfbMCC: string | null | undefined;
    csfbMNC: string | null | undefined;
    lac: number | null | undefined;
  };
};

export type GatewayV1 = GatewaySharedFields & {
  rawGateway: LteGateway;
};

export type GatewayPayload = {
  gateway_id: GatewayId;
  status?: GatewayStatusPayload;
  record?: AccessGatewayRecord;
  name?: GatewayName;
};

type SystemStatus = {
  time?: number;
  uptime_secs?: number;
  cpu_user?: number;
  cpu_system?: number;
  cpu_idle?: number;
  mem_total?: number;
  mem_available?: number;
  mem_used?: number;
  mem_free?: number;
  swap_total?: number;
  swap_used?: number;
  swap_free?: number;
  disk_partitions?: Array<DiskPartition>;
};

type PlatformInfo = {
  vpn_ip?: string;
  packages?: Array<SoftwarePackage>;
  kernel_version?: string;
  kernel_versions_installed?: Array<string>;
  config_info?: ConfigInfo;
};

type MachineInfo = {
  cpu_info?: {
    core_count?: number;
    threads_per_core?: number;
    architecture?: string;
    model_name?: string;
  };
  network_info?: {
    network_interfaces?: Array<NetworkInterface>;
    routing_table?: Array<Route>;
  };
};

type NetworkInterface = {
  network_interface_id?: string;
  status?: 'UP' | 'DOWN' | 'UNKNOWN';
  mac_address?: string;
  ip_addresses?: Array<string>;
  ipv6_addresses?: Array<string>;
};

type DiskPartition = {
  device?: string;
  mount_point?: string;
  total?: number;
  used?: number;
  free?: number;
};

type SoftwarePackage = {
  name?: string;
  version?: string;
};

type ConfigInfo = {
  mconfig_created_at?: number;
};

type Route = {
  destination_ip?: string;
  gateway_ip?: string;
  genmask?: string;
  network_interface_id?: string;
};

type GatewayName = string;

type ChallengeKey = {
  key_type: 'ECHO' | 'SOFTWARE_ECDSA_SHA256';
  key?: string;
};

type AccessGatewayRecord = {
  hardware_id: string;
  key: ChallengeKey;
};

// TODO: strip out devmand related fields and put them into a separate file
type GatewayMeta = {
  gps_latitude: number;
  gps_longitude: number;
  rf_tx_on: boolean;
  enodeb_connected: number;
  gps_connected: number;
  mme_connected: number;
  devmand: string | null | undefined;
  status: string | null | undefined;
};

type GatewayStatusPayload = {
  checkin_time?: number;
  hardware_id?: string;
  version?: string;
  system_status?: SystemStatus;
  platform_info?: PlatformInfo;
  machine_info?: MachineInfo;
  cert_expiration_time?: number;
  meta?: GatewayMeta;
  vpn_ip?: string;
  kernel_version?: string;
  kernel_versions_installed?: Array<string>;
};

export type FederationGatewayHealthStatus = {
  status: string;
  service_status: Record<string, ServiceStatusHealth>;
};

export const GatewayTypeEnum = Object.freeze({
  HEALTHY_GATEWAY: 'Good',
  UNHEALTHY_GATEWAY: 'Bad',
  UNKNOWN: '-',
});

// health status used for federation gateways
export const HEALTHY_STATUS = 'HEALTHY';
export const UNHEALTHY_STATUS = 'UNHEALTHY';

// availability status of federation gateway health service
export const AVAILABLE_STATUS = 'AVAILABLE';

export const ServiceTypeEnum = Object.freeze({
  HEALTHY_SERVICE: 'Up',
  UNHEALTHY_SERVICE: 'Down',
  UNENABLED_SERVICE: 'Not Enabled',
  UNAVAILABLE_SERVICE: 'N/A',
});

/**
 * Returns health status of the federation gateway.
 *
 * @param {network_id} networkId: Id of the federation network.
 * @param {gateway_id} gatewayId: Id of the gateway
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error
 * @returns the health status of the gateway or an empty string for an error.
 */
export async function getFederationGatewayHealthStatus(
  networkId: NetworkId,
  gatewayId: GatewayId,
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined,
): Promise<FederationGatewayHealthStatus> {
  try {
    const gwHealthStatus = (
      await MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdHealthStatusGet(
        {networkId, gatewayId},
      )
    ).data;
    return {
      status: gwHealthStatus.status,
      service_status: gwHealthStatus.service_status ?? {},
    };
  } catch (e) {
    enqueueSnackbar?.(
      'failed fetching health status information for federation gateway with id ' +
        gatewayId,
      {
        variant: 'error',
      },
    );
    return {status: '', service_status: {}};
  }
}

export const DynamicServices = Object.freeze({
  MONITORD: 'monitord',
  EVENTD: 'eventd',
  TD_AGENT_BIT: 'td-agent-bit',
});

export const DEFAULT_GATEWAY_CONFIG: LteGateway = {
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
  checked_in_recently: false,
  tier: 'default',
};

export const DEFAULT_DNS_CONFIG: NetworkDnsConfig = {
  dhcp_server_enabled: false,
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

export const DEFAULT_HE_CONFIG: GatewayHeConfig = {
  enable_encryption: false,
  encryption_key: '',
  enable_header_enrichment: false,
  he_encoding_type: 'BASE64',
  he_encryption_algorithm: 'RC4',
  he_hash_function: 'MD5',
};

export const DEFAULT_GW_POOL_CONFIG: MutableCellularGatewayPool = {
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
