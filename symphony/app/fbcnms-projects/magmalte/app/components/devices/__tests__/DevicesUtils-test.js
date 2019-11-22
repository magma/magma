/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {RAW_AGENT} from '../test/DevicesMock';
import {buildDevicesAgentFromPayload} from '../DevicesUtils';

it('disect agent', () => {
  const agent = buildDevicesAgentFromPayload(
    RAW_AGENT,
    RAW_AGENT.status?.checkin_time,
  );
  expect(agent).toEqual({
    checkinTime: RAW_AGENT.status?.checkin_time,
    readTime: RAW_AGENT.status?.checkin_time,
    devmand_config: {
      managed_devices: ['ping_fb_dns_from_lab', 'ping_google_ipv6'],
    },
    hardware_id: 'faceb00c-face-b00c-face-000c2940b2bf',
    id: 'fbbosfbcdockerengine',
    lastCheckin: '0.00 min ago',
    status: RAW_AGENT.status,
    up: true,
    version: 'Not Reported',
    rawAgent: RAW_AGENT,
  });
});
