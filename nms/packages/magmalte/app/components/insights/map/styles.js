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
 * @flow local
 * @format
 */

import mapboxgl from 'mapbox-gl';

import 'mapbox-gl/dist/mapbox-gl.css';

mapboxgl.accessToken = window.CONFIG.MAPBOX_ACCESS_TOKEN;

const OSM_STYLE = {
  version: 8,
  sources: {
    'osm-raster': {
      type: 'raster',
      tiles: [
        '//a.tile.openstreetmap.org/{z}/{x}/{y}.png',
        '//b.tile.openstreetmap.org/{z}/{x}/{y}.png',
      ],
      tileSize: 256,
    },
  },
  layers: [
    {
      id: 'osm-raster',
      type: 'raster',
      source: 'osm-raster',
      minzoom: 0,
      maxzoom: 22,
    },
  ],
};

const FB_STREETS_MAP_STYLE = 'mapbox://styles/fbmaps/cjnurl0x351tg2srp2jm1f6l1';
const FB_SATELLITE_MAP_STYLE =
  'mapbox://styles/fbmaps/cjzwgcnnv0rpi1cny70j6hvj6';

const MAPBOX_OUTDOORS_STYLE = 'mapbox://styles/mapbox/outdoors-v10';

export type MapType = 'satellite' | 'streets';

export function getDefaultMapStyle() {
  return mapboxgl.accessToken ? MAPBOX_OUTDOORS_STYLE : OSM_STYLE;
}

export function getMapStyleForType(mapType: MapType) {
  if (mapboxgl.accessToken) {
    return mapType == 'satellite'
      ? FB_SATELLITE_MAP_STYLE
      : FB_STREETS_MAP_STYLE;
  } else {
    return OSM_STYLE;
  }
}
