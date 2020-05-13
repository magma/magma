/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PropertyType} from '../../../../common/PropertyType';
import type {PropertyTypeTableDispatcherActionType} from './PropertyTypeTableDispatcherActionType';

import {generateTempId} from '../../../../common/EntUtils';
import {getInitialState, reducer} from './PropertyTypesTableReducer';
import {useReducer} from 'react';

export type PropertyTypesTableState = Array<PropertyType>;

export const getInitialPropertyType = (index: number): PropertyType => ({
  id: generateTempId(),
  name: '',
  index: index,
  type: 'string',
  nodeType: null,
  booleanValue: false,
  stringValue: null,
  intValue: null,
  floatValue: null,
  latitudeValue: null,
  longitudeValue: null,
  rangeFromValue: null,
  rangeToValue: null,
  isEditable: true,
  isInstanceProperty: true,
});

export const usePropertyTypesReducer = (
  initialPropertyTypes: Array<PropertyType>,
) => {
  return useReducer<
    PropertyTypesTableState,
    PropertyTypeTableDispatcherActionType,
    Array<PropertyType>,
  >(reducer, initialPropertyTypes, getInitialState);
};
