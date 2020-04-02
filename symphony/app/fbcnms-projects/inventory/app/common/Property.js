/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FragmentReference} from 'relay-runtime';
import type {PropertyFormField_property} from '../components/form/__generated__/PropertyFormField_property.graphql';
import type {PropertyType} from './PropertyType';

import DateTimeFormat from './DateTimeFormat.js';
import {toMutablePropertyType} from './PropertyType';

export type Property = {|
  id?: ?string,
  propertyType: PropertyType,

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
  equipmentValue?: ?{id: string, name: string},
  locationValue?: ?{id: string, name: string},
  serviceValue?: ?{id: string, name: string},
|};

export const sortPropertiesByIndex = (a: Property, b: Property) =>
  (a.propertyType.index ?? 0) - (b.propertyType.index ?? 0);

export const getNonInstancePropertyTypes = (
  properties: Array<Property>,
  propertyTypes: Array<PropertyType>,
): Array<PropertyType> => {
  properties = properties || [];
  const propIds = properties.map(x => x.propertyType.id);
  return propertyTypes.filter(type => !propIds.includes(type.id));
};

export const getPropertyValue = (property: Property | PropertyType) => {
  {
    const type = property.propertyType
      ? property.propertyType.type
      : property.type;
    switch (type) {
      case 'date':
      case 'email':
      case 'enum':
      case 'string':
        return property.stringValue;
      case 'datetime_local':
        return DateTimeFormat.dateTime(property.stringValue);
      case 'bool':
        return property.booleanValue != undefined
          ? property.booleanValue.toString()
          : '';
      case 'int':
        return property.intValue;
      case 'float':
        return property.floatValue;
      case 'range':
        return property.rangeFromValue !== null &&
          property.rangeToValue !== null
          ? (property.rangeFromValue ?? '') +
              ' - ' +
              (property.rangeToValue ?? '')
          : '';
      case 'gps_location':
        return property.latitudeValue !== null &&
          property.longitudeValue !== null
          ? (property.latitudeValue ?? '') +
              ', ' +
              (property.longitudeValue ?? '')
          : '';
      /**
       * Since this function accepts either property or property type,
       * we need to check which one we recieved.
       * In the case of PropertyType, there isn't an equipment/location value.
       */
      case 'equipment':
        return property.propertyType ? property.equipmentValue?.name : null;
      case 'location':
        return property.propertyType ? property.locationValue?.name : null;
      case 'service':
        return property.propertyType ? property.serviceValue?.name : null;
    }
  }
};

export const toPropertyInput = (properties: Array<Property>): Array<any> => {
  return properties
    .map(property => ({
      ...property,
      propertyTypeID: property.propertyType.id,
    }))
    .map(propInput => {
      const {propertyType: _, ...newPropInput} = propInput;
      return newPropInput;
    })
    .map(property => {
      if ((property.id && property.id.includes('@tmp')) || property.id == '0') {
        const {id: _, ...newProp} = property;
        return newProp;
      }
      return property;
    })
    .map(property => ({
      ...property,
      equipmentValue: undefined,
      equipmentIDValue: property.equipmentValue?.id ?? null,
      locationValue: undefined,
      locationIDValue: property.locationValue?.id ?? null,
      serviceValue: undefined,
      serviceIDValue: property.serviceValue?.id ?? null,
    }));
};

export const toMutableProperty = (
  immutableProperty: $ReadOnly<
    $Diff<PropertyFormField_property, {$refType: FragmentReference, ...}>,
  >,
): Property => ({
  id: immutableProperty.id,
  propertyType: toMutablePropertyType(immutableProperty.propertyType),
  booleanValue: immutableProperty.booleanValue,
  stringValue: immutableProperty.stringValue,
  intValue: immutableProperty.intValue,
  floatValue: immutableProperty.floatValue,
  latitudeValue: immutableProperty.latitudeValue,
  longitudeValue: immutableProperty.longitudeValue,
  rangeFromValue: immutableProperty.rangeFromValue,
  rangeToValue: immutableProperty.rangeToValue,
  equipmentValue:
    immutableProperty.equipmentValue != null
      ? {
          id: immutableProperty.equipmentValue.id,
          name: immutableProperty.equipmentValue.name,
        }
      : null,
  locationValue:
    immutableProperty.locationValue != null
      ? {
          id: immutableProperty.locationValue.id,
          name: immutableProperty.locationValue.name,
        }
      : null,
  serviceValue:
    immutableProperty.serviceValue != null
      ? {
          id: immutableProperty.serviceValue.id,
          name: immutableProperty.serviceValue.name,
        }
      : null,
});
