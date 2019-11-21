/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Equipment} from '../../common/Equipment';
import type {EquipmentComparisonViewQueryRendererSearchQueryResponse} from './__generated__/EquipmentComparisonViewQueryRendererSearchQuery.graphql.js';
import type {
  FilterConfig,
  FilterValue,
  FiltersQuery,
} from './ComparisonViewTypes';

import InventoryQueryRenderer from '../InventoryQueryRenderer';
import PowerSearchBar from '../power_search/PowerSearchBar';
import PowerSearchEquipmentResultsTable from './PowerSearchEquipmentResultsTable';
import React, {useState} from 'react';
import SearchIcon from '@material-ui/icons/Search';
import Text from '@fbcnms/ui/components/design-system/Text';
import useLocationTypes from './hooks/locationTypesHook';
import usePropertyFilters from './hooks/propertiesHook';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {EquipmentCriteriaConfig} from './EquipmentSearchConfig';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {
  buildPropertyFilterConfigs,
  getInitialFilterValue,
  getPossibleProperties,
} from './FilterUtils';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const PROPERTY_FILTER_NAME = 'inst_property';

const useStyles = makeStyles(theme => ({
  noResultsRoot: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: '100px',
  },
  noResultsLabel: {
    color: theme.palette.grey[600],
  },
  searchIcon: {
    color: theme.palette.grey[600],
    marginBottom: '6px',
    fontSize: '36px',
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
  },
  searchBar: {
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
}));

type Props = {
  onEquipmentSelected?: (equipment: Equipment) => void,
  limit?: number,
  showExport?: boolean,
};

const equipmentSearchQuery = graphql`
  query EquipmentComparisonViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [EquipmentFilterInput!]!
  ) {
    equipmentSearch(limit: $limit, filters: $filters) {
      equipment {
        ...PowerSearchEquipmentResultsTable_equipment
      }
      count
    }
  }
`;

const EquipmentComparisonViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {limit, onEquipmentSelected, showExport} = props;
  const [filters, setFilters] = useState(([]: FiltersQuery));
  const [count, setCount] = useState((0: number));
  const {history} = useRouter();

  const equipmentDataResponse = usePropertyFilters('equipment');

  const possibleProperties = getPossibleProperties(
    equipmentDataResponse.response,
  );
  const equipmentPropertiesFilterConfigs = buildPropertyFilterConfigs(
    possibleProperties,
  );

  const locationTypesFilterConfigs = useLocationTypes();

  const filterConfigs = EquipmentCriteriaConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .concat(equipmentPropertiesFilterConfigs)
    .concat(locationTypesFilterConfigs);

  return (
    <div className={classes.root}>
      <div className={classes.searchBar}>
        <PowerSearchBar
          filterValues={filters}
          exportPath={showExport ? '/equipment' : null}
          onFiltersChanged={setFilters}
          onFilterRemoved={handleFilterRemoved}
          onFilterBlurred={handleFilterBlurred}
          getSelectedFilter={(filterConfig: FilterConfig) =>
            getInitialFilterValue(
              filterConfig.key,
              filterConfig.name,
              filterConfig.defaultOperator,
              filterConfig.name === PROPERTY_FILTER_NAME
                ? possibleProperties.find(
                    propDef => propDef.name === filterConfig.label,
                  )
                : null,
            )
          }
          placeholder="Filter equipment"
          searchConfig={EquipmentCriteriaConfig}
          filterConfigs={filterConfigs}
          footer={
            count != null
              ? limit != null && count > limit
                ? `1 to ${limit} of ${count}`
                : `1 to ${count}`
              : null
          }
        />
      </div>
      <InventoryQueryRenderer
        query={equipmentSearchQuery}
        variables={{
          limit: limit,
          filters: filters.map(f => ({
            filterType: f.name.toUpperCase(),
            operator: f.operator.toUpperCase(),
            stringValue: f.stringValue,
            propertyValue: f.propertyValue,
            idSet: f.idSet,
          })),
        }}
        render={(
          props: EquipmentComparisonViewQueryRendererSearchQueryResponse,
        ) => {
          const {count, equipment} = props.equipmentSearch;
          setCount(count);
          if (!equipment || equipment.length === 0) {
            return (
              <div className={classes.noResultsRoot}>
                <SearchIcon className={classes.searchIcon} />
                <Text variant="h6" className={classes.noResultsLabel}>
                  No results found
                </Text>
              </div>
            );
          }
          return (
            <div className={classes.searchResults}>
              <PowerSearchEquipmentResultsTable
                equipment={equipment}
                onEquipmentSelected={onEquipmentSelected}
                onWorkOrderSelected={workOrderId =>
                  history.replace(`inventory?workorder=${workOrderId}`)
                }
              />
            </div>
          );
        }}
      />
    </div>
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

export default EquipmentComparisonViewQueryRenderer;
