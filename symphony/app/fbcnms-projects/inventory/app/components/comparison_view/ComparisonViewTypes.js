/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterEntity} from './hooks/__generated__/filterBookmarksHookReportFiltersQuery.graphql';
import type {PropertyType} from '../../common/PropertyType';

export const EntityTypeMap = Object.freeze({
  equipment: 'equipment',
  location: 'location',
  location_by_types: 'location_by_types',
  link: 'link',
  port: 'port',
  properties: 'properties',
  work_order: 'work_order',
  cell_scan: 'cell_scan',
  wifi_scan: 'wifi_scan',
  service: 'service',
  alert: 'alert',
});

export type EntityType = $Values<typeof EntityTypeMap>;

export const OperatorMap = Object.freeze({
  is: 'is',
  contains: 'contains',
  is_one_of: 'is_one_of',
  is_not_one_of: 'is_not_one_of',
  date_greater_than: 'date_greater_than',
  date_less_than: 'date_less_than',
});

export type Operator = $Values<typeof OperatorMap>;

export type EntityLocationFilter = {
  id: string,
  name: string,
};

export type FilterConfig = {
  key: string,
  name: string,
  entityType: EntityType,
  label: string,
  component: Object,
  defaultOperator: Operator,
  extraData?: ?Object,
};

export type SavedSearchConfig = {
  id: string,
  label: string,
  key: string,
  entity: FilterEntity,
  filters: Array<FilterValue>,
};

export type EntityConfig = {
  type: EntityType,
  label: string,
  filters: Array<FilterConfig>,
};

export type FilterValue = {
  id: string,
  key: string,
  name: string,
  operator: Operator,
  stringValue?: ?string,
  idSet?: ?Array<string>,
  stringSet?: ?Array<string>,
  boolValue?: ?boolean,
  propertyValue?: ?PropertyType,
};

export type FilterProps = {
  config: FilterConfig,
  onInputBlurred: () => void,
  onNewInputBlurred: (newValue: FilterValue) => void,
  value: FilterValue,
  editMode: boolean,
  onValueChanged: (newValue: FilterValue) => void,
  onRemoveFilter: () => void,
  title?: string,
};

export type FiltersQuery = Array<FilterValue>;
