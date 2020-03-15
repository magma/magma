/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {
  FilterConfig,
  FilterValue,
  FiltersQuery,
} from './ComparisonViewTypes';

import PowerSearchBar from '../power_search/PowerSearchBar';
import React from 'react';
import useLocationTypes from './hooks/locationTypesHook';
import usePropertyFilters from './hooks/propertiesHook';
import {EquipmentCriteriaConfig} from './EquipmentSearchConfig';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {buildPropertyFilterConfigs, getSelectedFilter} from './FilterUtils';

type Props = {
  filters: FiltersQuery,
  onFiltersChanged: FiltersQuery => void,
  footer?: ?string,
};

const EquipmentPowerSearchBar = (props: Props) => {
  const {onFiltersChanged, filters, footer} = props;

  const possibleProperties = usePropertyFilters('equipment');
  const equipmentPropertiesFilterConfigs = buildPropertyFilterConfigs(
    possibleProperties,
  );

  const locationTypesFilterConfigs = useLocationTypes();

  const filterConfigs = EquipmentCriteriaConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .concat(equipmentPropertiesFilterConfigs ?? [])
    .concat(locationTypesFilterConfigs ?? []);

  return (
    <PowerSearchBar
      filters={filters}
      filterValues={filters}
      exportPath={'/equipment'}
      onFiltersChanged={onFiltersChanged}
      onFilterRemoved={handleFilterRemoved}
      onFilterBlurred={handleFilterBlurred}
      getSelectedFilter={(filterConfig: FilterConfig) =>
        getSelectedFilter(filterConfig, possibleProperties ?? [])
      }
      placeholder="Filter..."
      searchConfig={EquipmentCriteriaConfig}
      filterConfigs={filterConfigs}
      footer={footer}
      entity="EQUIPMENT"
    />
  );
};

const handleFilterRemoved = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.EQUIPMENT_COMPARISON_VIEW_FILTER_REMOVED, {
    filterName: filter.name,
  });
};

const handleFilterBlurred = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.EQUIPMENT_COMPARISON_VIEW_FILTER_SET, {
    value: filter,
  });
};

export default EquipmentPowerSearchBar;
