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
import {augmentDevicesMap, buildDevicesAgentFromPayload} from '../DevicesUtils';

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

  const devicesAgents = augmentDevicesMap(RAW_DEVICES);

  expect(devicesAgents).toEqual({
    ens_switch_1: {
      id: 'ens_switch_1',
      managingAgentId: '',
      config: {
        channels: {
          cambium_channel: {
            client_id: 'randomid',
            client_ip: '10.0.0.1',
            client_mac: '58:C1:7A:90:36:50',
            client_secret: 'randomsecret',
          },
          frinx_channel: {
            authorization: 'Basic auth',
            device_type: 'ios',
            device_version: '15.2',
            frinx_port: 8181,
            host: 'frinx',
            password: 'frinx',
            port: 23,
            transport_type: 'telnet',
            username: 'username',
          },
          other_channel: {},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_config: '{}',
        device_type: ['snmp'],
        host: '2620:10d:c089:1:822a:a8ff:fe1c:d3c1',
        platform: 'snmp',
      },
      status: null,
    },
    localhost_snmpd: {
      id: 'localhost_snmpd',
      managingAgentId: '',
      config: {
        channels: {
          cambium_channel: {
            client_id: 'randomid',
            client_ip: '10.0.0.1',
            client_mac: '58:C1:7A:90:36:50',
            client_secret: 'randomsecret',
          },
          frinx_channel: {
            authorization: 'Basic auth',
            device_type: 'ios',
            device_version: '15.2',
            frinx_port: 8181,
            host: 'frinx',
            password: 'frinx',
            port: 23,
            transport_type: 'telnet',
            username: 'username',
          },
          other_channel: {},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_config: '{}',
        device_type: ['snmp'],
        host: '127.0.0.1',
        platform: 'snmp',
      },
      status: null,
    },
    mikrotik: {
      id: 'mikrotik',
      managingAgentId: '',
      config: {
        channels: {
          other_channel: {channel_props: {password: '', username: 'admin'}},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_type: [],
        host: '192.168.90.1',
        platform: 'mikrotik',
      },
      status: null,
    },
    ping_fb_dns_from_lab: {
      id: 'ping_fb_dns_from_lab',
      managingAgentId: 'fbbosfbcdockerengine',
      config: {
        channels: {},
        device_config: '{}',
        device_type: [],
        host: '192.168.96.18',
        platform: 'ping',
      },
      status: null,
    },
    ping_fb_dns_ken_laptop: {
      id: 'ping_fb_dns_ken_laptop',
      managingAgentId: '',
      config: {
        channels: {
          cambium_channel: {
            client_id: 'randomid',
            client_ip: '10.0.0.1',
            client_mac: '58:C1:7A:90:36:50',
            client_secret: 'randomsecret',
          },
          frinx_channel: {
            authorization: 'Basic auth',
            device_type: 'ios',
            device_version: '15.2',
            frinx_port: 8181,
            host: 'frinx',
            password: 'frinx',
            port: 23,
            transport_type: 'telnet',
            username: 'username',
          },
          other_channel: {},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_config: '{}',
        device_type: ['snmp'],
        host: '192.168.96.18',
        platform: 'ping',
      },
      status: null,
    },
    ping_google_ipv6: {
      id: 'ping_google_ipv6',
      managingAgentId: 'fbbosfbcdockerengine',
      config: {
        channels: {},
        device_config: '{}',
        device_type: [],
        host: '2607:f8b0:4004:814::200e',
        platform: 'ping',
      },
      status: null,
    },
    ping_google_ipv6_ken_laptop: {
      id: 'ping_google_ipv6_ken_laptop',
      managingAgentId: '',
      config: {
        channels: {
          cambium_channel: {
            client_id: 'randomid',
            client_ip: '10.0.0.1',
            client_mac: '58:C1:7A:90:36:50',
            client_secret: 'randomsecret',
          },
          frinx_channel: {
            authorization: 'Basic auth',
            device_type: 'ios',
            device_version: '15.2',
            frinx_port: 8181,
            host: 'frinx',
            password: 'frinx',
            port: 23,
            transport_type: 'telnet',
            username: 'username',
          },
          other_channel: {},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_config: '{}',
        device_type: ['snmp'],
        host: '2607:f8b0:4004:803::200e',
        platform: 'ping',
      },
      status: null,
    },
    ubnt: {
      id: 'ubnt',
      managingAgentId: '',
      config: {
        channels: {
          cambium_channel: {
            client_id: 'randomid',
            client_ip: '10.0.0.1',
            client_mac: '58:C1:7A:90:36:50',
            client_secret: 'randomsecret',
          },
          frinx_channel: {
            authorization: 'Basic auth',
            device_type: 'ios',
            device_version: '15.2',
            frinx_port: 8181,
            host: 'frinx',
            password: 'frinx',
            port: 23,
            transport_type: 'telnet',
            username: 'username',
          },
          other_channel: {},
          snmp_channel: {community: 'public', version: 'v1'},
        },
        device_config: '{}',
        device_type: ['snmp'],
        host: '192.168.88.253',
        platform: 'Ubnt',
      },
      status: null,
    },
  });
});
