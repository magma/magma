/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {GeoJSONGeometry} from '@mapbox/geojson-types';

export type ProjectGeoJSONFeature = {
  type: 'Feature',
  geometry: ?GeoJSONGeometry,
  properties: ?{
    name: string,
    id: string,
    numberOfWorkOrders: number,
    location: ProjectLocation,
  },
  id?: number | string,
};

export type ProjectLocation = {
  id: string,
  name: string,
  latitude: number,
  longitude: number,
};

export type ProjectMapMarkerData = {
  id: string,
  name: string,
  location: ProjectLocation,
  numberOfWorkOrders: number,
};

export type ProjectGeoJSONFeatureCollection = {
  type: 'FeatureCollection',
  features: Array<ProjectGeoJSONFeature>,
};

export const projectToGeoJson = (
  projects: Array<ProjectMapMarkerData>,
): ProjectGeoJSONFeatureCollection => {
  return {
    type: 'FeatureCollection',
    features: projects.map<ProjectGeoJSONFeature>(project =>
      projectToGeoFeature(project),
    ),
  };
};

export const projectToGeoFeature = (
  project: ProjectMapMarkerData,
): ProjectGeoJSONFeature => {
  return {
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates: [project.location.longitude, project.location.latitude],
    },
    properties: {
      id: project.id,
      name: project.name,
      numberOfWorkOrders: project.numberOfWorkOrders,
      location: project.location,
    },
  };
};
