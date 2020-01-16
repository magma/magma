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
  FilterConfig,
  FilterValue,
  FiltersQuery,
} from './ComparisonViewTypes';

import AppContext from '@fbcnms/ui/context/AppContext';
import PowerSearchBar from '../power_search/PowerSearchBar';
import React, {useContext} from 'react';
import useLocationTypes from './hooks/locationTypesHook';
import usePropertyFilters from './hooks/propertiesHook';
import {LinkCriteriaConfig} from './LinkSearchConfig';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {buildPropertyFilterConfigs, getSelectedFilter} from './FilterUtils';

type Props = {
  filters: FiltersQuery,
  onFiltersChanged: FiltersQuery => void,
  footer?: ?string,
};

const LinksPowerSearchBar = (props: Props) => {
  const {onFiltersChanged, filters, footer} = props;
  const {isFeatureEnabled} = useContext(AppContext);
  const linkStatusEnabled = isFeatureEnabled('planned_equipment');
  const locationTypesFilterConfigs = useLocationTypes();

  const possibleProperties = usePropertyFilters('link');
  const linkPropertiesFilterConfigs = buildPropertyFilterConfigs(
    possibleProperties,
  );

  const filterConfigs = LinkCriteriaConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .filter(conf => linkStatusEnabled || conf.key != 'link_future_status')
    .concat(linkPropertiesFilterConfigs ?? [])
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
      searchConfig={LinkCriteriaConfig}
      filterConfigs={filterConfigs}
      footer={footer}
      exportPath={'/links'}
    />
  );
};

const handleFilterRemoved = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.LINK_COMPARISON_VIEW_FILTER_REMOVED, {
    filterName: filter.name,
  });
};

const handleFilterBlurred = (filter: FilterValue) => {
  ServerLogger.info(LogEvents.LINK_COMPARISON_VIEW_FILTER_SET, {
    value: filter,
  });
};

export default LinksPowerSearchBar;
