/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {RAW_AGENT, RAW_DEVICES} from '../test/DevicesMock';
import {
  buildDevicesAgentFromPayload,
  mergeAgentsDevices,
} from '../DevicesUtils';

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

it('test merging devices and agents calls', () => {
  const agent = buildDevicesAgentFromPayload(
    RAW_AGENT,
    RAW_AGENT.status?.checkin_time,
  );

  if (agent.devmand_config?.managed_devices) {
    const managed_devices = agent.devmand_config?.managed_devices || [];
    managed_devices.push('nonexistdevice_but_ok_to_ignore');
  }

  const devicesAgents = mergeAgentsDevices([agent], Object.keys(RAW_DEVICES));

  expect(devicesAgents).toEqual({
    ens_switch_1: {
      id: 'ens_switch_1',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
    localhost_snmpd: {
      id: 'localhost_snmpd',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
    mikrotik: {
      id: 'mikrotik',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
    ping_fb_dns_from_lab: {
      id: 'ping_fb_dns_from_lab',
      agentIds: ['fbbosfbcdockerengine'],
      config: null,
      statusAgentId: 'fbbosfbcdockerengine',
      status: {
        'fbc-symphony-device:system': {
          'geo-location': {
            'reference-frame': {
              'astronomical-body': 'earth',
              'geodetic-system': {'geodetic-datum': 'wgs-84'},
            },
            latitude: 0,
            longitude: 0,
            height: 0,
          },
          latencies: {
            latency: [{type: 'ping', src: 'agent', dst: 'device', rtt: 11797}],
          },
          status: 'UP',
        },
      },
    },
    ping_fb_dns_ken_laptop: {
      id: 'ping_fb_dns_ken_laptop',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
    ping_google_ipv6: {
      id: 'ping_google_ipv6',
      agentIds: ['fbbosfbcdockerengine'],
      config: null,
      statusAgentId: 'fbbosfbcdockerengine',
      status: {
        'fbc-symphony-device:system': {
          status: 'UP',
          latencies: {
            latency: [{rtt: 12296, dst: 'device', src: 'agent', type: 'ping'}],
          },
          'geo-location': {
            height: 0,
            longitude: 0,
            latitude: 0,
            'reference-frame': {
              'geodetic-system': {'geodetic-datum': 'wgs-84'},
              'astronomical-body': 'earth',
            },
          },
        },
      },
    },
    ping_google_ipv6_ken_laptop: {
      id: 'ping_google_ipv6_ken_laptop',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
    ubnt: {
      id: 'ubnt',
      agentIds: [],
      config: null,
      statusAgentId: null,
      status: null,
    },
  });
});
