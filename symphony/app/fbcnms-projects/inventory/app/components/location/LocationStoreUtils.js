/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import nullthrows from '@fbcnms/util/nullthrows';
import {ConnectionHandler} from 'relay-runtime';

import type {RecordProxy, RecordSourceSelectorProxy} from 'relay-runtime';

export const removeLocationFromStore = (
  store: RecordSourceSelectorProxy,
  locationId: string,
  parentLocationId: ?string,
) => {
  if (parentLocationId != null) {
    const parentProxy = nullthrows(store.get(parentLocationId));
    const currNodes = parentProxy.getLinkedRecords('children') || [];
    const withoutCurrentLocation = currNodes.filter(
      child => child?.getDataID() !== locationId,
    );
    parentProxy.setLinkedRecords(withoutCurrentLocation, 'children');
    parentProxy.setValue(
      Number(parentProxy.getValue('numChildren')) - 1,
      'numChildren',
    );
  } else {
    const rootQuery = store.getRoot();
    const locations = nullthrows(
      ConnectionHandler.getConnection(rootQuery, 'LocationsTree_locations', {
        onlyTopLevel: true,
      }),
    );
    ConnectionHandler.deleteNode(locations, locationId);
  }
  store.delete(locationId);
};

export const addLocationToStore = (
  store: RecordSourceSelectorProxy,
  newNode: RecordProxy,
  parentLocationId: ?string,
) => {
  if (parentLocationId != null) {
    const parentProxy = nullthrows(store.get(parentLocationId));
    if (parentProxy != null) {
      const currNodes = parentProxy.getLinkedRecords('children') ?? [];
      const parentLoaded =
        currNodes !== null &&
        (currNodes.length === 0 || !!currNodes.find(node => node != undefined));
      if (parentLoaded) {
        parentProxy.setLinkedRecords([...currNodes, newNode], 'children');
        parentProxy.setValue(
          Number(parentProxy.getValue('numChildren')) + 1,
          'numChildren',
        );
      }
    }
  } else {
    const rootQuery = store.getRoot();
    const locations = nullthrows(
      ConnectionHandler.getConnection(rootQuery, 'LocationsTree_locations', {
        onlyTopLevel: true,
      }),
    );
    if (locations != null) {
      const edge = ConnectionHandler.createEdge(
        store,
        locations,
        newNode,
        'LocationsEdge',
      );
      ConnectionHandler.insertEdgeAfter(locations, edge);
    }
  }
};
