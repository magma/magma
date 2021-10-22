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

import {RAW_GATEWAY} from '../test/GatewayMock';
import {
  additionalPropsToArray,
  additionalPropsToObject,
  buildWifiGatewayFromPayloadV1,
  calculateLinkLocalMac,
  getAdditionalProp,
  isMeshNetwork,
  setAdditionalProp,
  wifiGeoJson,
  wifiGeoJsonConnections,
} from '../WifiUtils';

const NOW = 1561157125000;

it('renders without crashing', () => {
  const gateway = buildWifiGatewayFromPayloadV1(RAW_GATEWAY, NOW);
  expect(gateway).toEqual({
    id: 'shared_d_id_5ce28cf1aeb6',
    meshid: 'shared_d',
    hwid: 'faceb00c-face-b00c-face-5ce28cf1aeb6',
    info: 'binney lab, top shelf back wall',
    readTime: 1561157125000,
    checkinTime: 1561156214384,
    lastCheckin: '15.18 min ago',
    version:
      'Facebook Wi-Fi soma-image-1.0 (nbg6817) Release 52461bc4d24d+dirty (yerv@devvm354 Wed Jun 19 18:30:23 UTC 2019) (fbpkg:none) (cfg:none)',
    versionParsed: {
      buildtime: 'Wed Jun 19 18:30:23 UTC 2019',
      cfg: 'none',
      fbpkg: 'none',
      hash: '52461bc4d24d+dirty',
      user: 'yerv@',
    },
    wifi_config: {
      additional_props: {
        example_prop: 'lol',
        expected_is_gateway: '1',
      },
      client_channel: '1',
      info: 'binney lab, top shelf back wall',
      is_production: true,
      latitude: 83,
      longitude: -70,
      mesh_id: 'shared_d',
      mesh_rssi_threshold: -80,
    },
    up: false,
    upDanger: false,
    status: RAW_GATEWAY.status,
    coordinates: [-70, 83],
    isGateway: true,
  });
});

it('geojson connections', () => {
  const now = RAW_GATEWAY.status?.checkin_time || NOW;
  const gatewayA = buildWifiGatewayFromPayloadV1(
    {
      ...RAW_GATEWAY,
      name: '5ce28cf1ae00',
      id: 'devicea_id_5ce28cf1ae00',
      device: {
        hardware_id: 'faceb00c-face-b00c-face-5ce28cf1ae00',
        key: {key_type: 'ECHO'},
      },
      status: {
        ...RAW_GATEWAY.status,
        meta: {
          ...(RAW_GATEWAY.status?.meta || {}),
          openr_neighbors: '5ce28cf1ae01',
        },
      },
    },
    now,
  );
  const gatewayB = buildWifiGatewayFromPayloadV1(
    {
      ...RAW_GATEWAY,
      name: '5ce28cf1ae01',
      id: 'deviceb_id_5ce28cf1ae01',
      device: {
        hardware_id: 'faceb00c-face-b00c-face-5ce28cf1ae01',
        key: {key_type: 'ECHO'},
      },
      status: {
        ...RAW_GATEWAY.status,
        meta: {
          ...(RAW_GATEWAY.status?.meta || {}),
          openr_neighbors: '5ce28cf1ae00',
        },
      },
    },
    now,
  );
  const connections = wifiGeoJsonConnections([gatewayB, gatewayA]);
  expect(connections).toEqual([
    {
      geometry: {
        coordinates: [
          [-70, 83],
          [-70, 83],
        ],
        type: 'LineString',
      },
      properties: {
        deviceId0: 'devicea_id_5ce28cf1ae00',
        deviceId1: 'deviceb_id_5ce28cf1ae01',
        deviceInfo0: 'binney lab, top shelf back wall',
        deviceInfo1: 'binney lab, top shelf back wall',
        deviceMac0: '5ce28cf1ae00',
        deviceMac1: '5ce28cf1ae01',
        highestConnectionType: 'openr',
        id: 'devicea_id_5ce28cf1ae00-deviceb_id_5ce28cf1ae01',
        info0to1_L2ExpectedThroughput: undefined,
        info0to1_L2Exptime: undefined,
        info0to1_L2InactiveTime: undefined,
        info0to1_L2MeshPlink: undefined,
        info0to1_L2Metric: undefined,
        info0to1_L2RxBitrate: undefined,
        info0to1_L2Signal: undefined,
        info0to1_OpenrIp: undefined,
        info0to1_OpenrIpv6: undefined,
        info0to1_OpenrMetric: undefined,
        info1to0_L2ExpectedThroughput: undefined,
        info1to0_L2Exptime: undefined,
        info1to0_L2InactiveTime: undefined,
        info1to0_L2MeshPlink: undefined,
        info1to0_L2Metric: undefined,
        info1to0_L2RxBitrate: undefined,
        info1to0_L2Signal: undefined,
        info1to0_OpenrIp: undefined,
        info1to0_OpenrIpv6: undefined,
        info1to0_OpenrMetric: undefined,
        meshid0: 'shared_d',
        meshid1: 'shared_d',
        name: 'binney lab, top shelf back wall binney lab, top shelf back wall',
        title:
          'binney lab, top shelf back wall <--> binney lab, top shelf back wall: openr',
        unidirectional: false,
      },
      type: 'Feature',
    },
  ]);
});

it('geojson collection', () => {
  const gateway = buildWifiGatewayFromPayloadV1(RAW_GATEWAY, NOW);
  const geojson = wifiGeoJson([gateway], new Set([gateway.id]));
  const expected = [
    {
      type: 'FeatureCollection',
      features: [
        {
          type: 'Feature',
          geometry: {
            coordinates: [-70, 83],
            type: 'Point',
          },
          properties: {
            id: gateway.id,
            iconSize: [60, 60],
            meshid: gateway.meshid,
            name: gateway.info,
            category: 'Device',
            device: gateway,
            status: gateway.status,
            useLargeIcon: true,
            version: gateway.version,
          },
        },
      ],
    },
    [],
  ];
  expect(geojson).toEqual(expected);
});

it('setAdditionalProp test', () => {
  const value = [
    ['a', '1'],
    ['b', '2'],
  ];
  const expected = [
    ['a', '1'],
    ['b', '2'],
    ['c', '3'],
  ];
  setAdditionalProp(value, 'c', '3');
  expect(value).toEqual(expected);

  setAdditionalProp(value, 'c', null);
  expect(value).toEqual([
    ['a', '1'],
    ['b', '2'],
  ]);

  setAdditionalProp(value, 'b', '5');
  expect(value).toEqual([
    ['a', '1'],
    ['b', '5'],
  ]);
});

it('getAdditionalProp test', () => {
  expect(
    getAdditionalProp(
      [
        ['a', '1'],
        ['b', '2'],
        ['c', '3'],
      ],
      'b',
    ),
  ).toEqual('2');
  expect(getAdditionalProp(null, 'b')).toEqual(null);
});

it('additionalPropsToArray test', () => {
  const initial = {a: '1', b: '2'};
  const expected = [
    ['a', '1'],
    ['b', '2'],
  ];
  const result = additionalPropsToArray(initial);
  expect(result).toEqual(expected);

  expect(additionalPropsToObject(expected)).toEqual(initial);

  expect(additionalPropsToObject(null)).toEqual(null);
  expect(additionalPropsToArray(null)).toEqual(null);
});

it('mesh network validation', () => {
  expect(isMeshNetwork('mesh_test')).toEqual(true);
  expect(isMeshNetwork('not_mesh')).toEqual(false);
});

it('calculateLinkLocalMac test', () => {
  expect(calculateLinkLocalMac('fe80::5683:3aff:feb0:2762')).toEqual(
    '54833ab02762',
  );
  expect(calculateLinkLocalMac('fe80::5683:3aff:feb0:286a')).toEqual(
    '54833ab0286a',
  );
});
