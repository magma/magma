/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Property} from './Property';
import type {PropertyKind} from '../components/form/__generated__/PropertyTypeFormField_propertyType.graphql';

export type PropertyType = {|
  id: string,
  type: PropertyKind,
  name: string,
  index: number,
  // one or more of the following potential value fields will have actual data,
  // depending on the property type selected for this property.
  // e.g. for 'email' the stringValue field will be populated
  booleanValue?: ?boolean,
  stringValue?: ?string,
  intValue?: ?number,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};

export const getPropertyDefaultValue = (propertyType: PropertyType) => {
  {
    switch (propertyType.type) {
      case 'date':
      case 'email':
      case 'enum':
      case 'string':
        return propertyType.stringValue;
      case 'bool':
        return propertyType.booleanValue != undefined
          ? propertyType.booleanValue.toString()
          : '';
      case 'int':
        return propertyType.intValue;
      case 'float':
        return propertyType.floatValue;
      case 'range':
        return propertyType.rangeFromValue !== null &&
          propertyType.rangeToValue !== null
          ? (propertyType.rangeFromValue ?? '') +
              ' - ' +
              (propertyType.rangeToValue ?? '')
          : '';
      case 'gps_location':
        return propertyType.latitudeValue !== null &&
          propertyType.longitudeValue !== null
          ? (propertyType.latitudeValue ?? '') +
              ', ' +
              (propertyType.longitudeValue ?? '')
          : '';
      case 'equipment':
      case 'location':
      case 'service':
        return '';
    }
  }
};

export const getInitialPropertyFromType = (
  propType: PropertyType,
): Property => {
  return {
    id: 'prop@tmp' + propType.id,
    propertyType: propType,
    booleanValue: propType.booleanValue,
    stringValue: propType.type !== 'enum' ? propType.stringValue : '',
    intValue: propType.intValue,
    floatValue: propType.floatValue,
    latitudeValue: propType.latitudeValue,
    longitudeValue: propType.longitudeValue,
    rangeFromValue: propType.rangeFromValue,
    rangeToValue: propType.rangeToValue,
    isInstanceProperty: propType.isInstanceProperty,
    equipmentValue: null,
  };
};
