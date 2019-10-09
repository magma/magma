/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

// TODO: Coordinate with murtadha on type
export type CheckindGateway = {
  gateway_id: string,
  config: {[string]: any},
  meta: {
    enodeb_connected: ?boolean,
    rf_tx_on: ?boolean,
  },
  name: string,
  status: ?{
    checkin_time: number,
    hardware_id: string,
    version: string,
    system_status: {
      time: number,
      cpu_user: number,
      cpu_system: number,
      cpu_idle: number,
      mem_total: number,
      mem_available: number,
      mem_used: number,
      mem_free: number,
      uptime_secs: number,
    },
    cert_expiration_time: number,
    meta: {[string]: string},
    vpn_ip: string,
    kernel_version: string,
  },
  record: {
    hardware_id: string,
    key: {
      key_type: string,
      key: string,
    },
  },
  offset: number,
};

export type WifiConfig = {
  mesh_id?: string,
  info?: string,
  longitude?: number,
  latitude?: number,
  client_channel: string,
  is_production: boolean,
  additional_props: ?{[string]: string},
};

export type CellularConfig = {
  attached_enodeb_serials: Array<string>,
  epc: {
    nat_enabled: boolean,
    ip_block: string,
  },
  ran: {
    pci: string | number,
    transmit_enabled: boolean,
  },
  non_eps_service: ?{
    non_eps_service_control: string | number,
    csfb_rat: string | number,
    csfb_mcc: string | number,
    csfb_mnc: string | number,
    lac: string | number,
  },
};

export type MagmadConfig = {
  autoupgrade_enabled: boolean,
  autoupgrade_poll_interval: number,
  checkin_interval: number,
  checkin_timeout: number,
  tier: string,
};

export type Record = {
  hardware_id: string,
};

export type NetworkUpgradeImage = {
  name: string,
  order: number,
};

export type UpgradeReleaseChannel = {
  name: string,
  supported_versions: Array<string>,
};

export type DevmandConfig = {
  managed_devices: string[],
};
