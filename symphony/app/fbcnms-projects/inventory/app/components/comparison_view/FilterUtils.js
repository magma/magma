/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  EntityLocationFilter,
  FilterConfig,
  FilterValue,
  Operator,
} from './ComparisonViewTypes';
import type {PropertyType} from '../../common/PropertyType';
import type {locationTypesHookLocationTypesQueryResponse} from './hooks/__generated__/locationTypesHookLocationTypesQuery.graphql.js';
import type {propertiesHookPossiblePropertiesQueryResponse} from './hooks/__generated__/propertiesHookPossiblePropertiesQuery.graphql.js';

import PowerSearchLocationFilter from '../comparison_view/PowerSearchLocationFilter';
import PowerSearchPropertyFilter from './PowerSearchPropertyFilter';
import shortid from 'shortid';
import {getPropertyValue} from '../../common/Property';
import {groupBy} from 'lodash';

export const PROPERTY_FILTER_NAME = 'property';

export const dateValues = [
  {
    value: 'date_greater_than',
    label: 'Greater than',
  },
  {
    value: 'date_less_than',
    label: 'Less than',
  },
];

export function getOperatorLabel(operator: Operator): string {
  switch (operator) {
    case 'is':
    case 'is_one_of':
      return 'is';
    case 'contains':
      return 'contains';
    case 'date_greater_than':
      return 'greater than';
    case 'date_less_than':
      return 'less than';
  }

  throw new Error(`Operator ${operator} doesn't have a label`);
}

export function getSelectedFilter(
  filterConfig: FilterConfig,
  possibleProperties: Array<PropertyType>,
): FilterValue {
  return getInitialFilterValue(
    filterConfig.key,
    filterConfig.name,
    filterConfig.defaultOperator,
    filterConfig.name === PROPERTY_FILTER_NAME
      ? possibleProperties.find(
          propDef =>
            filterConfig.key === `property_${propDef.name}_${propDef.type}`,
        )
      : null,
  );
}

export function getInitialFilterValue(
  key: string,
  name: string,
  operator: Operator,
  propertyType?: ?PropertyType,
): FilterValue {
  return {
    id: shortid.generate(),
    key,
    name,
    operator,
    stringValue: null,
    idSet: null,
    boolValue: null,
    propertyValue: propertyType
      ? {
          id: propertyType.id,
          type: propertyType.type,
          name: propertyType.name,
          index: propertyType.index,
          stringValue: propertyType.stringValue,
        }
      : null,
  };
}

export const buildLocationTypeFilterConfigs = (
  locationTypes: Array<EntityLocationFilter>,
): Array<FilterConfig> => {
  return locationTypes.map(type => ({
    key: `location_${type.id}`,
    name: 'location_inst',
    entityType: 'location_by_types',
    label: type.name,
    component: PowerSearchLocationFilter,
    defaultOperator: 'is_one_of',
    extraData: {
      locationTypeId: type.id,
    },
  }));
};

export function doesFilterHasValue(filterValue: FilterValue): boolean {
  const propValue = filterValue.propertyValue;
  return (
    !!filterValue.stringValue ||
    filterValue.boolValue != null ||
    (!!filterValue.idSet && filterValue.idSet.length > 0) ||
    (!!propValue && !!getPropertyValue(propValue))
  );
}

export function getLocationTypes(
  data: ?locationTypesHookLocationTypesQueryResponse,
): Array<EntityLocationFilter> {
  if (data == null || data.locationTypes == null) {
    return [];
  }

  return (data.locationTypes.edges ?? [])
    .filter(edge => edge != null && edge.node != null)
    .map(edge => ({
      id: edge?.node?.id || '',
      name: edge?.node?.name || '',
    }));
}

export const buildPropertyFilterConfigs = (
  definitions: ?Array<PropertyType>,
): Array<FilterConfig> => {
  if (definitions == null) {
    return [];
  }

  return definitions
    .filter(
      d =>
        d.type !== 'equipment' && d.type !== 'location' && d.type !== 'service',
    )
    .map(definition => ({
      key: `property_${definition.name}_${definition.type}`,
      name: 'property',
      entityType: 'properties',
      label: definition.name,
      component: PowerSearchPropertyFilter,
      defaultOperator: definition.type === 'date' ? 'date_less_than' : 'is', // Take from property type
    }));
};

export function getPossibleProperties(
  data: ?propertiesHookPossiblePropertiesQueryResponse,
): Array<PropertyType> {
  if (data == null || data.possibleProperties == null) {
    return [];
  }
  const propertiesGroup: {[string]: Array<PropertyType>} = groupBy(
    data.possibleProperties
      .filter(prop => prop.type !== 'gps_location' && prop.type !== 'range')
      .map((prop, index) => ({
        id: prop.name + prop.type,
        type: prop.type,
        name: prop.name,
        index: index,
        stringValue: prop.stringValue,
      })),
    prop => prop.name + prop.type,
  );
  const supportedProperties: Array<PropertyType> = [];
  for (const k in propertiesGroup) {
    supportedProperties.push(propertiesGroup[k][0]);
  }
  return supportedProperties;
}
