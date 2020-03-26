/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {gray14, green30, orange} from '@fbcnms/ui/theme/colors';
import type {BasicLocation} from '../../common/Location';
import type {
  GeoJSONFeature,
  GeoJSONFeatureCollection,
} from '@mapbox/geojson-types';
import type {GeoJSONSource} from './MapView';
import type {Location} from '../../common/Location.js';
import type {ShortUser} from '../../common/EntUtils';
import type {WorkOrder} from '../../common/WorkOrder';
import type {
  WorkOrderPriority,
  WorkOrderStatus,
} from '../../mutations/__generated__/EditWorkOrderMutation.graphql';

export type CoordsWithProps = {
  latitude: number,
  longitude: number,
  properties: Object,
};

export type WorkOrderLocation = BasicLocation & {
  randomizedLatitude: number,
};

export type WorkOrderWithLocation = {
  workOrder: WorkOrder,
  location: WorkOrderLocation,
};

export const locationsToGeoJSONSource = (
  key: string,
  locations: Array<Location>,
  properties: Object,
): GeoJSONSource => {
  return {
    key: key,
    data: {
      type: 'FeatureCollection',
      features: locations.map<GeoJSONFeature>(location =>
        locationToGeoFeature(location, properties),
      ),
    },
  };
};

export const workOrderToGeoJSONSource = (
  key: string,
  workOrders: Array<WorkOrderWithLocation>,
  properties: Object,
): GeoJSONSource => {
  return {
    key: key,
    data: {
      type: 'FeatureCollection',
      features: workOrders.map<GeoJSONFeature>(workOrder =>
        workOrderToGeoFeature(workOrder, properties),
      ),
    },
  };
};

export type WorkOrderProperties = {
  id: string,
  name: string,
  description: string,
  status: WorkOrderStatus,
  priority: WorkOrderPriority,
  owner: ShortUser,
  assignee: ?ShortUser,
  installDate: string,
  location: WorkOrderLocation,
  iconStatus: string,
  iconTech: string,
  text: string,
  textColor: string,
};

export const workOrderToGeoFeature = (
  workOrder: WorkOrderWithLocation,
  properties: WorkOrderProperties,
): GeoJSONFeature => {
  return {
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates: [
        workOrder.location.longitude,
        workOrder.location.randomizedLatitude,
      ],
    },
    properties: {
      id: workOrder.workOrder.id,
      name: workOrder.workOrder.name,
      description: workOrder.workOrder.description,
      status: workOrder.workOrder.status,
      priority: workOrder.workOrder.priority,
      owner: workOrder.workOrder.owner,
      assignee: workOrder.workOrder.assignedTo,
      installDate: workOrder.workOrder.installDate,
      location: workOrder.workOrder.location,
      iconStatus: getWorkOrderStatusIcon(workOrder.workOrder.status),
      iconTech: workOrder.workOrder.assignedTo
        ? 'icon_pin'
        : 'unassignedActive',
      text: workOrder.workOrder.assignedTo
        ? workOrder.workOrder.assignedTo.email.slice(0, 2)
        : '',
      textColor: getWorkOrderIconTextColor(workOrder.workOrder.status),
      ...properties,
    },
  };
};

const getWorkOrderStatusIcon = (status: string) => {
  if (status === 'DONE') {
    return 'doneActive';
  } else if (status == 'PENDING') {
    return 'pendingActive';
  }
  return 'plannedActive';
};

const getWorkOrderIconTextColor = (status: string) => {
  if (status === 'DONE') {
    return green30;
  } else if (status == 'PENDING') {
    return orange;
  }
  return gray14;
};

export const locationToGeoFeature = (
  location: {
    id: string,
    name: string,
    latitude: number,
    longitude: number,
  },
  properties: Object,
): GeoJSONFeature => {
  return {
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates: [location.longitude, location.latitude],
    },
    properties: {
      id: location.id,
      title: location.name,
      description: '',
      ...properties,
    },
  };
};

export const locationToGeoJson = (location: {
  id: string,
  name: string,
  latitude: number,
  longitude: number,
}): GeoJSONFeatureCollection => {
  return {
    type: 'FeatureCollection',
    features: [locationToGeoFeature(location)],
  };
};

export const coordsToGeoJson = (
  lat: number,
  long: number,
): GeoJSONFeatureCollection => {
  return {
    type: 'FeatureCollection',
    features: [
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [long, lat],
        },
        properties: {},
      },
    ],
  };
};

export const coordsToGeoJSONSource = (
  key: string,
  coordsWithProps: Array<CoordsWithProps>,
): GeoJSONSource => {
  const features: Array<GeoJSONFeature> = coordsWithProps.map(coords => ({
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates: [coords.longitude, coords.latitude],
    },
    properties: coords.properties,
  }));
  return {
    key: key,
    data: {
      type: 'FeatureCollection',
      features: features,
    },
  };
};

export const polygonToGeoJSONSource = (
  key: string,
  coords: Array<{latitude: number, longitude: number}>,
): GeoJSONSource => {
  const coordinates = [coords.map(coord => [coord.longitude, coord.latitude])];
  return {
    key: key,
    data: {
      type: 'FeatureCollection',
      features: [
        {
          type: 'Feature',
          geometry: {
            type: 'Polygon',
            coordinates: coordinates,
          },
          properties: null,
        },
      ],
    },
  };
};
