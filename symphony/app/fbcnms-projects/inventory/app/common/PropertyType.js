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
import type {PropertyFormField_property} from '../components/form/__generated__/PropertyFormField_property.graphql';
import type {PropertyKind} from '../components/form/__generated__/PropertyTypeFormField_propertyType.graphql';

export type PropertyType = {|
  id: string,
  type: PropertyKind,
  nodeType?: ?string,
  name: string,
  index?: ?number,
  category?: ?string,
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
      case 'node':
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
    nodeValue: null,
  };
};

export const toMutablePropertyType = (
  immutablePropertyType: $ReadOnly<
    $ElementType<PropertyFormField_property, 'propertyType'>,
  >,
): PropertyType => ({
  id: immutablePropertyType.id,
  type: immutablePropertyType.type,
  nodeType: immutablePropertyType.nodeType,
  name: immutablePropertyType.name,
  index: immutablePropertyType.index,
  category: immutablePropertyType.category,
  booleanValue: immutablePropertyType.booleanValue,
  stringValue: immutablePropertyType.stringValue,
  intValue: immutablePropertyType.intValue,
  floatValue: immutablePropertyType.floatValue,
  latitudeValue: immutablePropertyType.latitudeValue,
  longitudeValue: immutablePropertyType.longitudeValue,
  rangeFromValue: immutablePropertyType.rangeFromValue,
  rangeToValue: immutablePropertyType.rangeToValue,
  isEditable: immutablePropertyType.isEditable,
  isInstanceProperty: immutablePropertyType.isInstanceProperty,
  isMandatory: immutablePropertyType.isMandatory,
  isDeleted: immutablePropertyType.isDeleted,
});
