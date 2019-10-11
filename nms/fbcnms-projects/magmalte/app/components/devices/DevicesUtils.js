/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  DevmandConfig,
  Record,
} from '@fbcnms/magmalte/app/common/MagmaAPIType';

const MS_IN_MIN = 60 * 1000;
const MS_IN_HOUR = 60 * MS_IN_MIN;
const MS_IN_DAY = 24 * MS_IN_HOUR;

export type DevicesGatewayStatus = {
  checkin_time?: number,
  meta?: {
    [key: string]: string,
  },
};

export type DevicesGatewayPayload = {
  gateway_id: string,
  config: {
    devmand_gateway?: ?DevmandConfig,
    magmad_gateway?: {
      autoupgrade_enabled: boolean,
      autoupgrade_poll_interval: number,
      checkin_interval: number,
      checkin_timeout: number,
      tier: string,
    },
    [key: string]: string,
  },
  record: ?(Record & {
    [key: string]: mixed,
  }),
  status: ?DevicesGatewayStatus,
};

export type DevicesGateway = {
  id: string,

  devmand_config: ?DevmandConfig,

  hardware_id: string,

  readTime: number,
  checkinTime: ?number,
  lastCheckin: string,

  up: ?boolean,

  status: ?DevicesGatewayStatus,

  // TODO: deprecate this
  rawGateway: DevicesGatewayPayload,
};

export type DevicesManagedDevice = {
  device_config?: string,
  device_type?: Array<string>,
  host?: string,
  platform?: string,
  channels?: {
    snmp_channel?: {
      community: string,
      version: string,
    },
    frinx_channel?: {
      authorization: string,
      device_type: string,
      device_version: string,
      frinx_port: string,
      host: string,
      password: string,
      port: number,
      transport_type: string,
      username: string,
    },
    cambium_channel?: {
      client_id: string,
      client_ip: string,
      client_mac: string,
      client_secret: string,
    },
    other_channel?: {
      channel_props: {[string]: string},
    },
  },
};

export type DeviceStatus = {
  [key: string]: string | number | DeviceStatus | Array<DeviceStatus>,
};

export type FullDevice = {
  id: string,
  agentIds: string[], // list of agents that manage this device
  config: ?DevicesManagedDevice,
  // status meta from gateways are unstructured
  // eslint-disable-next-line flowtype/no-weak-types
  status: any,
  statusAgentId: string, // agentId that reported the status
};

export function buildDevicesGatewayFromPayload(
  gateway: DevicesGatewayPayload,
  now?: number,
): DevicesGateway {
  if (!gateway.record || !gateway.config) {
    throw Error('Cannot read gateway without `record` or `config`');
  }

  const currentTime = now === undefined ? new Date().getTime() : now;

  let lastCheckin = 'Not Reported';
  let version = 'Not Reported';
  let checkinTime = null;
  let up = null;

  const {status, config} = gateway;

  if (status) {
    checkinTime = status.checkin_time;
    if (status.meta?.version) {
      version = status.meta?.version;
    }

    if (status.checkin_time !== undefined) {
      const elapsedTime = Math.max(0, currentTime - status.checkin_time);
      if (elapsedTime > MS_IN_DAY) {
        lastCheckin = `${(elapsedTime / MS_IN_DAY).toFixed(2)} d ago`;
      } else if (elapsedTime > MS_IN_HOUR) {
        lastCheckin = `${(elapsedTime / MS_IN_HOUR).toFixed(2)} hr ago`;
      } else {
        lastCheckin = `${(elapsedTime / MS_IN_MIN).toFixed(2)} min ago`;
      }

      // up = within 2 minutes
      up = elapsedTime < 2 * MS_IN_MIN;
    }
  }

  return {
    id: gateway.gateway_id,

    hardware_id: gateway.record.hardware_id || 'Error: Missing hardware_id',

    devmand_config: config.devmand_gateway,

    readTime: currentTime,
    checkinTime,
    lastCheckin,
    version,

    up, // 2 minutes

    status,

    rawGateway: gateway,
  };
}

export const DEFAULT_DEVMAND_GATEWAY_CONFIGS = {
  autoupgrade_enabled: false,
  autoupgrade_poll_interval: 300,
  checkin_interval: 15,
  checkin_timeout: 12,
  tier: 'default',
};

/* get list of all devices from:
        devices (any and all devices that should exist)
        gateway status (device is reporting info)
        gateway devmand config (device should be managed by someone)

    In orc8r V1 API, this object would be returned natively in one API call
    and this function would be completely deleted.
  */
export function mergeGatewaysDevices(
  gateways: ?Array<DevicesGateway>,
  devices: ?Array<string>,
): {[key: string]: FullDevice} {
  const devicemap = {};
  if (devices) {
    // get list of all 'valid' devices
    devices.forEach(id => {
      devicemap[id] = {
        id,
        agentIds: [],
        config: null, // TODO: will exist after V1 of API
        statusAgentId: null,
        status: null,
      };
    });
  }

  if (gateways) {
    // gather "managed" devices
    gateways.forEach(gateway => {
      // map gateway.devmand_config.managed_devices to devicemap[id].agentId
      (gateway.devmand_config?.managed_devices || []).forEach(id => {
        if (id in devicemap) {
          devicemap[id].agentIds.push(gateway.id);
        } else {
          // If a device does not exist in devicemap,
          //    then we're in an inconsistent state.
          console.error(
            `Warning gateway ${gateway.id} is configured to manage non-existing device ${id}`,
          );
        }
      });
    });

    // gather "managed" devices
    gateways.forEach(gateway => {
      // status meta from gateways are unstructured
      // eslint-disable-next-line flowtype/no-weak-types
      let devmand: {[key: string]: any} = {};

      try {
        devmand = JSON.parse(gateway.status?.meta?.devmand || '{}');
      } catch (err) {
        console.error(err);
        return;
      }

      Object.keys(devmand)
        .filter(id => devmand[id])
        .map(id => {
          const status = devmand[id];
          if (!(id in devicemap)) {
            // If a device does not exist in devicemap,
            //   then we're in an inconsistent state, but display it anyway
            console.error(
              `Warning gateway ${gateway.id} is reporting state for non-existing device ${id}`,
            );
            // display it anyway because it's interesting
            devicemap[id] = {
              id,
              agentIds: [],
              config: null,
              status,
              statusAgentId: gateway.id,
            };
          } else {
            if (devicemap[id].status) {
              console.error(
                `Warning: device id ${id} managed by multiple devices`,
              );
            }
            devicemap[id].status = status;
            devicemap[id].statusAgentId = gateway.id;
          }
        });
    });
  }
  return devicemap;
}
