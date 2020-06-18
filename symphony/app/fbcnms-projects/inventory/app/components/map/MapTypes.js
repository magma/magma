/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MapFeature} from './GeoJSON';

import mapboxgl from 'mapbox-gl';

export type MapMarkerProps = {
  map: mapboxgl.Map,
  feature: MapFeature,
  onClick?: (string | number) => void,
};
