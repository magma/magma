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
  gateway_status,
  magmad_gateway_configs,
  symphony_agent,
  symphony_device,
  symphony_device_config,
} from '@fbcnms/magma-api';

const MS_IN_MIN = 60 * 1000;
const MS_IN_HOUR = 60 * MS_IN_MIN;
const MS_IN_DAY = 24 * MS_IN_HOUR;

export type DevicesAgent = {
  id: string,

  devmand_config: ?{
    managed_devices: string[],
  },

  hardware_id: string,

  readTime: number,
  checkinTime: ?number,
  lastCheckin: string,

  up: ?boolean,

  status: ?gateway_status,

  // TODO: deprecate this
  rawAgent: symphony_agent,
};

export type DeviceStatus = {
  [key: string]: string | number | DeviceStatus | Array<DeviceStatus>,
};

export type FullDevice = {
  id: string,
  config: ?symphony_device_config,
  managingAgentId: ?string,
  // status meta from agents are unstructured
  // eslint-disable-next-line flowtype/no-weak-types
  status: any,
};

export function buildDevicesAgentFromPayload(
  agent: symphony_agent,
  now?: number,
): DevicesAgent {
  const currentTime = now === undefined ? new Date().getTime() : now;

  let lastCheckin = 'Not Reported';
  let version = 'Not Reported';
  let checkinTime = null;
  let up = null;

  const {status} = agent;

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
    id: agent.id,

    hardware_id: agent?.device?.hardware_id || 'Error: Missing hardware_id',

    devmand_config: {managed_devices: agent.managed_devices},

    readTime: currentTime,
    checkinTime,
    lastCheckin,
    version,

    up, // 2 minutes

    status,

    rawAgent: agent,
  };
}

export const DEFAULT_MAGMAD_CONFIGS: magmad_gateway_configs = {
  autoupgrade_enabled: false,
  autoupgrade_poll_interval: 300,
  checkin_interval: 15,
  checkin_timeout: 12,
};

export function augmentDevicesMap(devices: {
  [string]: symphony_device,
}): {[key: string]: FullDevice} {
  const devicemap = {};

  Object.keys(devices).map(id => {
    const managingAgentId = devices[id]?.managing_agent;

    /*
    Currently, raw_state contains a JSON string using a YANG model. In the
    future this string should be returned expanded and typed; this logic here
    does the expansion on the client side while we figure out how to
    autogenerate swagger and javascript types based on a YANG model.

    The 'catch' clause is in case raw_state is published incorrectly, so the UI
    can detect the error and display something sane.
    */
    let status = null;
    if (devices[id].state?.raw_state) {
      try {
        status = JSON.parse(devices[id].state?.raw_state || '{}');
      } catch (e) {
        status = devices[id].state;
      }
    }

    devicemap[id] = {
      id,
      managingAgentId,
      config: devices[id].config,
      status,
    };
  });
  return devicemap;
}
