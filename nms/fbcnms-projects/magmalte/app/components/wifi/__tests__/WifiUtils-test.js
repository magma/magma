/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {RAW_GATEWAY} from '../test/GatewayMock';
import {
  additionalPropsToArray,
  additionalPropsToObject,
  buildWifiGatewayFromPayload,
  calculateLinkLocalMac,
  getAdditionalProp,
  isMeshNetwork,
  setAdditionalProp,
  wifiGeoJson,
  wifiGeoJsonConnections,
} from '../WifiUtils';

const NOW = 1561157125000;

it('renders without crashing', () => {
  const gateway = buildWifiGatewayFromPayload(RAW_GATEWAY, NOW);
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
  const gateway = buildWifiGatewayFromPayload(RAW_GATEWAY, NOW);
  const connections = wifiGeoJsonConnections([gateway, gateway]);
  expect(connections).toEqual([]);
});

it('geojson collection', () => {
  const gateway = buildWifiGatewayFromPayload(RAW_GATEWAY, NOW);
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
  const value = [['a', '1'], ['b', '2']];
  const expected = [['a', '1'], ['b', '2'], ['c', '3']];
  setAdditionalProp(value, 'c', '3');
  expect(value).toEqual(expected);

  setAdditionalProp(value, 'c', null);
  expect(value).toEqual([['a', '1'], ['b', '2']]);

  setAdditionalProp(value, 'b', '5');
  expect(value).toEqual([['a', '1'], ['b', '5']]);
});

it('getAdditionalProp test', () => {
  expect(getAdditionalProp([['a', '1'], ['b', '2'], ['c', '3']], 'b')).toEqual(
    '2',
  );
  expect(getAdditionalProp(null, 'b')).toEqual(null);
});

it('additionalPropsToArray test', () => {
  const initial = {a: '1', b: '2'};
  const expected = [['a', '1'], ['b', '2']];
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
