/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PropertyType} from '../../../../common/PropertyType';
import type {PropertyTypeTableDispatcherActionType} from './PropertyTypeTableDispatcherActionType';
import type {PropertyTypesTableState} from './PropertyTypesTableState';

import {getInitialPropertyType} from './PropertyTypesTableState';
import {reorder} from '../../../draggable/DraggableUtils';
import {sortByIndex} from '../../../draggable/DraggableUtils';

export function getInitialState(
  propertyTypes: Array<PropertyType>,
): PropertyTypesTableState {
  return propertyTypes.length === 0
    ? [getInitialPropertyType(0)]
    : propertyTypes.slice().map(p => ({...p}));
}

function editPropertyType<T: PropertyTypesTableState>(
  state: T,
  updatedPropertyTypeId: string,
  updatingCallback: PropertyType => PropertyType,
): T {
  const propertyTypeIndex = state.findIndex(
    p => p.id === updatedPropertyTypeId,
  );
  return [
    ...state.slice(0, propertyTypeIndex),
    updatingCallback(state[propertyTypeIndex]),
    ...state.slice(propertyTypeIndex + 1),
  ];
}

export function reducer(
  state: PropertyTypesTableState,
  action: PropertyTypeTableDispatcherActionType,
): PropertyTypesTableState {
  switch (action.type) {
    case 'ADD_PROPERTY_TYPE':
      return [...state, getInitialPropertyType(state.length)];
    case 'REMOVE_PROPERTY_TYPE':
      return editPropertyType(state, action.id, pt => ({
        ...pt,
        isDeleted: true,
      }));
    case 'UPDATE_PROPERTY_TYPE_NAME':
      return editPropertyType(state, action.id, pt => ({
        ...pt,
        name: action.name,
      }));
    case 'UPDATE_PROPERTY_TYPE_KIND':
      return editPropertyType(state, action.id, pt => ({
        ...getInitialPropertyType(pt.index ?? 0),
        id: action.id,
        type: action.kind,
        name: pt.name,
        nodeType: action.nodeType,
      }));
    case 'UPDATE_PROPERTY_TYPE':
      return editPropertyType(state, action.value.id, pt => ({
        ...pt,
        ...action.value,
      }));
    case 'CHANGE_PROPERTY_TYPE_INDEX':
      const sortedNotDeletedState = state
        .filter(pt => !pt.isDeleted)
        .sort(sortByIndex);
      return [
        ...reorder<PropertyType>(
          sortedNotDeletedState,
          action.sourceIndex,
          action.destinationIndex,
        ).map((p, index) => {
          return {
            ...p,
            index,
          };
        }),
        ...state.filter(pt => pt.isDeleted),
      ];
    default:
      return state;
  }
}
