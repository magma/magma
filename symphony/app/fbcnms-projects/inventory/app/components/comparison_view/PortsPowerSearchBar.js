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
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {PortCriteriaConfig} from './PortSearchConfig';
import {buildPropertyFilterConfigs, getSelectedFilter} from './FilterUtils';

type Props = {
  filters: FiltersQuery,
  onFiltersChanged: FiltersQuery => void,
  footer?: ?string,
};

const PortsPowerSearchBar = (props: Props) => {
  const {onFiltersChanged, filters, footer} = props;
  const locationTypesFilterConfigs = useLocationTypes();

  const possibleProperties = usePropertyFilters('port');
  const portPropertiesFilterConfigs = buildPropertyFilterConfigs(
    possibleProperties,
  );

  const filterConfigs = PortCriteriaConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .concat(portPropertiesFilterConfigs ?? [])
    .concat(locationTypesFilterConfigs ?? []);

  return (
    <PowerSearchBar
      filters={filters}
      filterValues={filters}
      onFiltersChanged={onFiltersChanged}
      onFilterRemoved={handleFilterRemoved}
      onFilterBlurred={handleFilterBlurred}
      getSelectedFilter={(filterConfig: FilterConfig) =>
        getSelectedFilter(filterConfig, possibleProperties ?? [])
      }
      placeholder="Filter..."
      searchConfig={PortCriteriaConfig}
      filterConfigs={filterConfigs}
      footer={footer}
      exportPath={'/ports'}
    />
  );
};

const handleFilterRemoved = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.PORT_COMPARISON_VIEW_FILTER_REMOVED, {
    filterName: filter.name,
  });
};

const handleFilterBlurred = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.PORT_COMPARISON_VIEW_FILTER_SET, {
    value: filter,
  });
};

export default PortsPowerSearchBar;
