/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {lte_gateway} from '@fbcnms/magma-api';

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
