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

import type {FilterSpecification} from 'mapbox-gl/src/style-spec/types';
import type {MagmaConnectionFeature} from '@fbcnms/ui/insights/map/GeoJSON';
import type {WifiGateway} from './WifiUtils';

import {CONNECTION_TYPE_L2_ESTAB, CONNECTION_TYPE_L3_OPENR} from './WifiUtils';

export type LineMapLayer = 'none' | 'l2' | 'l3' | 'defaultRoute';
type LineMapLayerConfig = [string, number, number];

type LineMapLayerConfigs = {
  [LineMapLayer]: LineMapLayerConfig,
};

const LINE_MAP_LAYER_CONFIG: LineMapLayerConfigs = {
  // [layer, line-color, line-width, line-opacity]
  none: ['black', 1.5, 0.2],
  l2: ['red', 2, 0.23],
  l3: ['blue', 3, 0.23],
  defaultRoute: ['green', 5, 0.3],
};

export function layerNameForFeature(
  feature: MagmaConnectionFeature,
): LineMapLayer {
  // qualify highestConnectionType by 'unidirectional'
  if (feature.properties.unidirectional) {
    return 'none';
  }
  if (
    feature.properties.info0to1_isDefaultRoute ||
    feature.properties.info1to0_isDefaultRoute
  ) {
    return 'defaultRoute';
  } else if (
    feature.properties.highestConnectionType === CONNECTION_TYPE_L3_OPENR
  ) {
    return 'l3';
  } else if (
    feature.properties.highestConnectionType === CONNECTION_TYPE_L2_ESTAB
  ) {
    return 'l2';
  } else {
    // CONNECTION_TYPE_L2_LISTEN
    return 'none';
  }
}

export function buildLayer(
  layerName: LineMapLayer,
  features: MagmaConnectionFeature[],
) {
  const [color, width, opacity] = LINE_MAP_LAYER_CONFIG[layerName];
  return {
    id: `lines-${layerName}`,
    type: 'line',
    source: {
      type: 'geojson',
      data: {
        type: 'FeatureCollection',
        features,
      },
    },
    layout: {
      'line-join': 'round',
      'line-cap': 'round',
    },
    paint: {
      'line-color': color,
      'line-width': width,
      'line-opacity': opacity,
    },
  };
}

export function searchDevices(
  stateDevices: WifiGateway[],
  search: ?string,
  clickedDevice: ?WifiGateway,
  meshID: ?string,
): [WifiGateway[], Map<string, FilterSpecification>] {
  let searchRegex;
  // if there's a problem with the regex, return equivalent of "no filter"
  try {
    searchRegex = new RegExp(search || '', 'i');
  } catch (err) {
    return [stateDevices, mapLayerFilters(meshID, clickedDevice, null)];
  }

  const macsMatchedSearch = !search ? null : [];

  // get mac addresses of devices whose 'info' field contains searchRegex
  const macsOfMatchedInfo = new Set(
    stateDevices
      .map((device: WifiGateway) => {
        if (meshID && device.meshid !== meshID) {
          return null;
        }

        if (device.info.toLowerCase().match(searchRegex)) {
          return device.hwid.slice(-12);
        }
        return null;
      })
      .filter(Boolean),
  );

  const devices = stateDevices.filter((device: WifiGateway) => {
    // always show clicked device
    if (clickedDevice) {
      return clickedDevice.id === device.id;
    }

    if (meshID && device.meshid !== meshID) {
      return false;
    }

    // macsMatchedSearch is null if macfilter is inactive
    if (macsMatchedSearch === null) {
      return true;
    }

    if (
      // check hwid against macs that previously matched 'info'
      macsOfMatchedInfo.has(device.hwid.slice(-12)) ||
      // check mac address filter
      device.hwid.match(searchRegex)
    ) {
      macsMatchedSearch.push(device.hwid.slice(-12));
      return true;
    }

    // don't match if no status
    if (!device.status) {
      return false;
    }

    return false;
  });
  return [devices, mapLayerFilters(meshID, clickedDevice, macsMatchedSearch)];
}

function mapLayerFilters(
  meshID: ?string,
  device: ?WifiGateway,
  macs: ?Array<string>,
): Map<string, FilterSpecification> {
  /*
    filter on meshID if not null (empty string will match everything)
    filter on device if not null (empty string will match everything)
    filter on macs if not null (empty [] will match nothing)
    */

  // filter *must* exist, default filter should "match all"
  let layerFilterMeshId = ['has', 'id'];
  if (meshID) {
    layerFilterMeshId = [
      'any',
      ['in', 'meshid0', meshID],
      ['in', 'meshid1', meshID],
    ];
  }

  // filter *must* exist, default filter should "match all"
  let layerFilterDevice = ['has', 'id'];
  if (device) {
    layerFilterDevice = [
      'any',
      ['in', 'deviceId0', device.id],
      ['in', 'deviceId1', device.id],
    ];
  }

  // filter *must* exist, default filter should "match all"
  let layerFilterMac = ['has', 'id'];
  if (macs) {
    // note that this also matches '[]', such that passing []
    //   will match nothing and exclude everything.
    let filters = ['any'];
    macs.forEach(
      mac =>
        (filters = filters.concat([
          ['in', 'deviceMac0', mac],
          ['in', 'deviceMac1', mac],
        ])),
    );
    layerFilterMac = filters;
  }

  // compose filter from MeshId and Device and Mac selections
  return new Map(
    Object.keys(LINE_MAP_LAYER_CONFIG).map(layerName => [
      `lines-${layerName}`,
      ['all', layerFilterMeshId, layerFilterDevice, layerFilterMac],
    ]),
  );
}
