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
 * @flow
 * @format
 */

import type {
  gateway_status,
  gateway_wifi_configs,
  wifi_gateway,
} from '@fbcnms/magma-api';

import nullthrows from '@fbcnms/util/nullthrows';
import type {
  MagmaConnectionFeature,
  MagmaFeatureCollection,
  MagmaGatewayFeature,
} from '@fbcnms/ui/insights/map/GeoJSON';

import {assign, flatMap, flatten, groupBy, partition} from 'lodash';

const MS_IN_MIN = 60 * 1000;
const MS_IN_HOUR = 60 * MS_IN_MIN;
const MS_IN_DAY = 24 * MS_IN_HOUR;

type WifiVersionAttributes = {
  hash: string,
  user: string,
  buildtime: string,
  fbpkg: string,
  cfg: string,
};

export type WifiGateway = {
  id: string,
  meshid: string,

  wifi_config: ?gateway_wifi_configs,

  hwid: string,
  info: string,

  readTime: number,
  checkinTime: ?number,
  lastCheckin: string,
  version: string,
  versionParsed: ?WifiVersionAttributes,

  up: ?boolean,
  upDanger: ?boolean,

  status: ?gateway_status,

  coordinates: [number, number],
  isGateway: ?boolean,
};

// parses long version str into key/value of type WifiVersionAttributes
const versionRegExp = new RegExp(
  [
    '.*\\s+Release\\s+', // Facebook Wi-Fi soma-image-1.0 (nbg6817) Release
    '([a-f0-9]+(?:\\+dirty)?)', // 96ee0a665ada+dirty
    '\\s+\\(', // <space>(
    '([\\w]+[@])[\\w]+\\s+', // alanw@devvm886
    '([^)]+)', // Mon Nov  5 18:03:48 PST 2018
    '\\)\\s+', // )<space>
    '(?:\\(fbpkg\\:([^)]*)\\))?', // (fbpkg:none)
    '\\s+', // <space>
    '(?:\\(cfg\\:([^)]*)\\))?', // (cfg:clown-alanwdev)
  ].join(''),
);

function _parseVersion(version: string): ?WifiVersionAttributes {
  const match = version.match(versionRegExp);
  if (match) {
    return {
      hash: match[1],
      user: match[2],
      buildtime: match[3],
      fbpkg: match[4],
      cfg: match[5],
    };
  }
  return null;
}

export function buildWifiGatewayFromPayloadV1(
  gateway: wifi_gateway,
  now?: number,
): WifiGateway {
  if (!gateway.device || !gateway.wifi) {
    throw Error(
      `Something very wrong with ${gateway.id}. Gateway without 'device' or 'wifi'`,
    );
  }

  const currentTime = now === undefined ? new Date().getTime() : now;

  let meshid = 'Not Configured';

  let info = 'Not Configured';

  let lastCheckin = 'Not Reported';
  let version = 'Not Reported';
  let versionParsed = null;
  let checkinTime = null;

  let up = null;
  let upDanger = null;

  let coordinates = [NaN, NaN];
  let isGateway = null;

  const {status, wifi} = gateway;

  if (status) {
    checkinTime = status.checkin_time;

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
      // danger zone is between 1-2 minutes
      up = elapsedTime < 2 * MS_IN_MIN;
      upDanger = up && elapsedTime > 1 * MS_IN_MIN;
    }

    if (status.meta) {
      if (status.meta.version) {
        version = status.meta.version;
        versionParsed = _parseVersion(version);
      }

      isGateway = false;
      // TODO: remove openr_inet_monitor clause when all
      //       devices are updated past D13547708
      if (status.meta.openr_inet_monitor) {
        // previously legacy, now back.
        if (status.meta.openr_inet_monitor === 'success') {
          isGateway = true;
        }
      } else if (status.meta.is_gateway === 'true') {
        isGateway = true;
      }
    }
  }

  let longitude = parseFloat(wifi?.longitude);
  let latitude = parseFloat(wifi?.latitude);

  // 180 is the antimeridian
  if (longitude !== NaN && Math.abs(longitude) > 180) longitude %= 180;
  // 90 is north/south pole
  if (latitude !== NaN && Math.abs(latitude) > 90) latitude %= 90;

  if (wifi) {
    meshid = wifi.mesh_id || meshid;
    info = wifi.info || info;
    coordinates = [longitude, latitude];
  }

  return {
    id: gateway.id,
    meshid,

    hwid: gateway.device?.hardware_id || 'Error: Missing hwid',
    info,

    wifi_config: wifi,

    readTime: currentTime,
    checkinTime,
    lastCheckin,
    version,
    versionParsed,

    up, // 2 minutes
    upDanger, // 1 minute

    status,

    coordinates,
    isGateway,
  };
}

export const macColonfy = (mac: string): string =>
  mac.replace(/(.{2})(.{2})(.{2})(.{2})(.{2})(.{2})/, '$1:$2:$3:$4:$5:$6');

function getOpenrNeighbors(device: WifiGateway): string[] {
  return (device.status?.meta?.openr_neighbors || '').split(',');
}

function getL2Peers(device: WifiGateway): {[string]: string[]} {
  // TODO: consider filter for inactive time?
  const meta = device.status?.meta || {};
  const mesh0_stations = meta.mesh0_stations || '';
  return groupBy(
    mesh0_stations.split(','),
    mac => meta[`mesh0_${mac}_mesh_plink`] || 'UNKNOWN',
  );
}

// given a device object, return list of default route ip v4 and v6 addresses
function getDefaultRouteIpList(device: WifiGateway): string[] {
  // get v4 default routes
  // meta.default_route in form of: interface1,ip1;interface2,ip2
  const defaultRouteIfIpList =
    (device.status &&
      device.status.meta &&
      device.status.meta.default_route &&
      device.status.meta.default_route.split(';')) ||
    [];

  return defaultRouteIfIpList
    .map(ifIpStr => {
      const ifIp = ifIpStr.split(',');
      // ensure that interface is mesh0
      if (ifIp.length === 2 && ifIp[0] === 'mesh0') {
        return ifIp[1];
      }
    })
    .filter(Boolean);
}

export function calculateLinkLocalMac(ipv6: string): ?string {
  // Get local-link mac address (no colons) given a v6 IP
  // See https://tools.ietf.org/html/rfc4291 appendix A
  // expecting fe80::xxxx:xxff:fexx:xxxx
  // where "x" portion is based on mac address

  // ** this is not intended to be a perfect calculation - this is intended only
  // to work in the wifi specific use case.

  const res = ipv6.toLowerCase().split('::');
  if (res.length != 2 || res[0] !== 'fe80') {
    return null;
  }

  // expecting xxxx:xxff:fexx:xxxx
  let ipMac = res[1].replace(/:/g, '');
  if (ipMac.length !== 16) {
    return null;
  }
  ipMac = ipMac.substring(0, 6) + ipMac.substring(10);

  // flip second-to-last bit in first octet
  const firstOctet = parseInt(ipMac.substring(0, 2), 16);
  const flippedOctet = firstOctet ^ parseInt('00000010', 2);
  ipMac = flippedOctet.toString(16) + ipMac.substring(2);

  return ipMac;
}

function _gatherOpenrData(
  prefix: string,
  device: WifiGateway,
  macOther: string,
): {[string]: string | Boolean} {
  // macOther must be mac address _without_ colons
  if (!device.status || !device.status.meta) {
    return {};
  }
  const {meta} = device.status;

  const data = {};
  data[`${prefix}_OpenrMetric`] = meta[`openr_${macOther}_metric`];
  data[`${prefix}_OpenrIp`] = meta[`openr_${macOther}_ip`];
  data[`${prefix}_OpenrIpv6`] = meta[`openr_${macOther}_ipv6`];

  if (
    getDefaultRouteIpList(device).filter(
      ip =>
        /* if v4 IP matches */
        ip === meta[`openr_${macOther}_ip`] ||
        /* if v6 IP matches */
        ip === meta[`openr_${macOther}_ipv6`] ||
        /* if v6 local-link address matches */
        calculateLinkLocalMac(ip) === macOther,
    ).length > 0
  ) {
    data[`${prefix}_isDefaultRoute`] = true;
  }

  return data;
}

function _gatherL2Data(
  prefix: string,
  device: WifiGateway,
  macOther: string,
): {[string]: string} {
  // macOther must be mac address _with_ colons
  if (!device.status || !device.status.meta) {
    return {};
  }
  const {meta} = device.status;

  // key contains underscore because geojson properties
  //   does not support named objects as values
  const data = {};
  data[`${prefix}_L2ExpectedThroughput`] =
    meta[`mesh0_${macOther}_expected_throughput`];
  data[`${prefix}_L2Exptime`] = meta[`mesh0_${macOther}_exptime`];
  data[`${prefix}_L2InactiveTime`] = meta[`mesh0_${macOther}_inactive_time`];
  data[`${prefix}_L2MeshPlink`] = meta[`mesh0_${macOther}_mesh_plink`];
  data[`${prefix}_L2Metric`] = meta[`mesh0_${macOther}_metric`];
  data[`${prefix}_L2RxBitrate`] = meta[`mesh0_${macOther}_rx_bitrate`];
  data[`${prefix}_L2Signal`] = meta[`mesh0_${macOther}_signal`];
  return data;
}

function _geoJsonDevice(
  device: WifiGateway,
  useLargeIcon?: boolean,
): Array<MagmaGatewayFeature> {
  /*
    return list of features for this device, currently a single Point

    expect each device.coordinates to not include NaN
  */

  return [
    {
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: device.coordinates,
      },
      properties: {
        iconSize: [60, 60],
        useLargeIcon,
        id: device.id,

        name: device.info,
        category: 'Device',

        meshid: device.meshid,
        version: device.version,
        status: device.status,

        device: device,
      },
    },
  ];
}

export const CONNECTION_TYPE_L2_ESTAB = '802.11s:ESTAB';
export const CONNECTION_TYPE_L2_LISTEN = '802.11s:LISTEN';
export const CONNECTION_TYPE_L3_OPENR = 'openr';

function _geoJsonConnection(
  deviceA: WifiGateway,
  deviceB: WifiGateway,
): ?MagmaConnectionFeature {
  /*
    LineString
    Connections between device A <--> device B may be any of:
      L3 (openr)
      L2 (802.11s)
      default route
  */

  if (deviceA.id == deviceB.id) {
    return null;
  }

  if (deviceA.coordinates.includes(NaN)) {
    return null;
  }

  if (deviceB.coordinates.includes(NaN)) {
    return null;
  }

  // require status.meta to enumerate neighbors/peers and device is up
  if (
    !deviceA.status ||
    !deviceB.status ||
    !deviceA.status.meta ||
    !deviceB.status.meta
  ) {
    return null;
  }

  let device0 = deviceA;
  let device1 = deviceB;

  // ensure consistent ordering
  if (device0.id > device1.id) {
    device0 = deviceB;
    device1 = deviceA;
  }

  const mac0 = device0.hwid.slice(-12);
  const macColons0 = macColonfy(mac0);
  const mac1 = device1.hwid.slice(-12);
  const macColons1 = macColonfy(mac1);

  // check if device0 and device1 has an L3 (openr) connection
  const openrNeighbors0 = device0.up ? getOpenrNeighbors(device0) : [];
  const openrNeighbors1 = device1.up ? getOpenrNeighbors(device1) : [];

  // check if device0 and device1 has an L2 (802.11s) connection
  const l2Peers0 = device0.up ? getL2Peers(device0) : {};
  const l2Peers1 = device1.up ? getL2Peers(device1) : {};

  let highestConnectionType;
  if (
    // detect if openr (L3 connectivity)
    openrNeighbors0.includes(mac1) &&
    openrNeighbors1.includes(mac0)
  ) {
    highestConnectionType = CONNECTION_TYPE_L3_OPENR;
  } else if (
    // detect if 802.11s:ESTAB (L2 connectivity)
    (l2Peers0['ESTAB'] || []).includes(macColons1) &&
    (l2Peers1['ESTAB'] || []).includes(macColons0)
  ) {
    highestConnectionType = CONNECTION_TYPE_L2_ESTAB;
  } else if (
    // detect if 802.11s:LISTEN (visible but can't communicate)
    flatMap(l2Peers0).includes(macColons1) &&
    flatMap(l2Peers1).includes(macColons0)
  ) {
    highestConnectionType = CONNECTION_TYPE_L2_LISTEN;
  }

  let highestUniConnectionType;
  if (
    // detect if openr (L3 connectivity) on EITHER side
    openrNeighbors0.includes(mac1) ||
    openrNeighbors1.includes(mac0)
  ) {
    highestUniConnectionType = CONNECTION_TYPE_L3_OPENR;
  } else if (
    // detect if 802.11s:ESTAB (L2 connectivity) on EITHER side
    (l2Peers0['ESTAB'] || []).includes(macColons1) ||
    (l2Peers1['ESTAB'] || []).includes(macColons0)
  ) {
    highestUniConnectionType = CONNECTION_TYPE_L2_ESTAB;
  } else if (
    // detect if 802.11s:LISTEN (visible but can't communicate) on EITHER side
    flatMap(l2Peers0).includes(macColons1) ||
    flatMap(l2Peers1).includes(macColons0)
  ) {
    highestUniConnectionType = CONNECTION_TYPE_L2_LISTEN;
  }

  // bail if no connection types
  let unidirectional = false; // begin by assuming bidirectional connection
  if (!highestConnectionType && !highestUniConnectionType) {
    return null;
  } else if (
    // if a bidirectional connection does not exist, consider unidirectional
    !highestConnectionType
  ) {
    unidirectional = true;
    highestConnectionType = highestUniConnectionType;
  }

  // construct geojson properties for this connection
  const properties = {
    title: `${device0.info} <--> ${device1.info}: ${nullthrows(
      highestConnectionType,
    )}${unidirectional ? ' (unidirectional)' : ''}`,

    id: `${device0.id}-${device1.id}`,
    name: `${device0.info} ${device1.info}`,

    deviceId0: device0.id,
    deviceId1: device1.id,

    deviceInfo0: device0.info,
    deviceInfo1: device1.info,

    deviceMac0: mac0,
    deviceMac1: mac1,

    // filtering properties to be used with map.setFilter()
    meshid0: device0.meshid,
    meshid1: device1.meshid,

    highestConnectionType, // qualified by unidirectional
    unidirectional,
  };

  assign(
    properties,
    // L3 openr attribute compilation
    device0.up ? _gatherOpenrData('info0to1', device0, mac1) : {},
    device1.up ? _gatherOpenrData('info1to0', device1, mac0) : {},
    // L2 attribute compilation
    device0.up ? _gatherL2Data('info0to1', device0, macColons1) : {},
    device1.up ? _gatherL2Data('info1to0', device1, macColons0) : {},
  );

  return {
    type: 'Feature',
    geometry: {
      type: 'LineString',
      coordinates: [device0.coordinates, device1.coordinates],
    },
    properties,
  };
}

export function wifiGeoJson(
  devices: Array<WifiGateway>,
  matchedDeviceIds?: Set<string>,
): [MagmaFeatureCollection, Array<WifiGateway>] {
  const [invalidDevices, validDevices] = partition(devices, device =>
    device.coordinates.includes(NaN),
  );

  return [
    {
      type: 'FeatureCollection',
      features: flatten(
        validDevices.map(d =>
          _geoJsonDevice(d, matchedDeviceIds && matchedDeviceIds.has(d.id)),
        ),
      ),
    },
    invalidDevices,
  ];
}

export function wifiGeoJsonConnections(
  devices: Array<WifiGateway>,
): Array<MagmaConnectionFeature> {
  const connectionFeaturesArray = devices.map(
    deviceA =>
      devices
        .map(deviceB =>
          // using a comparison is a cheap way to ensure only gathers
          //   connections in one direction, use null as placeholder
          deviceA.id >= deviceB.id
            ? null
            : _geoJsonConnection(deviceA, deviceB),
        )
        .filter(Boolean), // removes null from list
  );

  return flatten(connectionFeaturesArray);
}

export function additionalPropsToArray(
  props: ?{[string]: string},
): ?Array<[string, string]> {
  if (!props) {
    return null;
  }

  return Object.keys(props)
    .sort()
    .map(key => [key, nullthrows(props)[key]]);
}

export function additionalPropsToObject(
  props: ?Array<[string, string]>,
): ?{[string]: string} {
  if (!props) {
    return null;
  }

  const results = {};
  props.filter(p => p[0] && p[1]).map(pair => (results[pair[0]] = pair[1]));
  return results;
}

export function getAdditionalProp(
  props?: ?Array<[string, string]>,
  prop: string,
) {
  if (props === null || props === undefined) {
    return null;
  }
  for (let i = 0; i < props.length; i++) {
    if (props[i][0] === prop) {
      return props[i][1];
    }
  }
  return null;
}

export function setAdditionalProp(
  props: Array<[string, string]>,
  prop: string,
  value: ?string,
) {
  // if value is undefined, remove prop from props
  // if prop already exists, then update value in props
  // if prop does not exist, then add to props
  // if a prop in props is blank, then use that entry to fill in prop/value

  if (value === undefined || value === null) {
    // remove prop, if it exists
    for (let i = 0; i < props.length; i++) {
      if (props[i][0] === prop) {
        props.splice(i, 1);
        break;
      }
    }
    return;
  }

  for (let i = 0; i < props.length; i++) {
    if (props[i][0] === '' || props[i][0] === prop) {
      props[i] = [prop, value];
      return;
    }
  }
  props.push([prop, value]);
  return;
}

export function isMeshNetwork(networkID: string): boolean {
  return networkID.startsWith('mesh');
}

export const DEFAULT_HW_ID_PREFIX = 'faceb00c-face-b00c-face-';
export const DEFAULT_WIFI_GATEWAY_CONFIGS = {
  autoupgrade_enabled: false,
  autoupgrade_poll_interval: 300,
  checkin_interval: 15,
  checkin_timeout: 12,
  tier: 'default',
};
