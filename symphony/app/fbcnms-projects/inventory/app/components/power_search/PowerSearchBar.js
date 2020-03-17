/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Bookmark} from './PowerSearchContext';
import type {
  EntityConfig,
  FilterConfig,
  FilterValue,
  FiltersQuery,
  SavedSearchConfig,
} from '../comparison_view/ComparisonViewTypes';
import type {FilterEntity} from '../../mutations/__generated__/AddReportFilterMutation.graphql';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import CSVFileExport from '../CSVFileExport';
import FilterBookmark from '../FilterBookmark';
import FiltersTypeahead from '../comparison_view/FiltersTypeahead';
import PowerSearchContext from './PowerSearchContext';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import update from 'immutability-helper';
import {
  configToFilterQuery,
  doesFilterHasValue,
} from '../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useRef, useState} from 'react';

const useStyles = makeStyles(theme => ({
  root: {
    position: 'relative',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    width: '100%',
    backgroundColor: theme.palette.common.white,
    padding: '8px 14px',
  },
  searchBarContainer: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    flexGrow: 1,
  },
  searchTypeahead: {
    display: 'flex',
    cursor: 'text',
  },
  headerContainer: {
    marginRight: theme.spacing(),
  },
  placeholder: {
    position: 'absolute',
    top: '0px',
    bottom: '0px',
    zIndex: 2,
    lineHeight: '36px',
    color: theme.palette.grey.A200,
    fontWeight: 'bold',
    pointerEvents: 'none',
    alignItems: 'center',
    display: 'flex',
  },
  filter: {
    marginRight: theme.spacing(),
  },
  typeahead: {
    flexGrow: 1,
  },
  footer: {
    padding: '8px 4px 8px 10px',
    color: theme.palette.grey.A200,
    fontWeight: 'bold',
    pointerEvents: 'none',
  },
}));

type Props = {
  filterValues?: FiltersQuery,
  placeholder?: string,
  filterConfigs: Array<FilterConfig>,
  savedSearches?: Array<SavedSearchConfig>,
  searchConfig: Array<EntityConfig>,
  header?: React.Node,
  footer?: ?string,
  exportPath?: ?string,
  className?: string,
  onFiltersChanged: (filters: Array<FilterValue>) => void,
  onFilterRemoved?: (filter: FilterValue) => void,
  onFilterBlurred?: (filter: FilterValue) => void,
  // used when a filter is selected from filter typeahead
  getSelectedFilter: (filterConfig: FilterConfig) => FilterValue,
  entity?: FilterEntity,
};

const PowerSearchBar = (props: Props) => {
  const filtersTypeaheadRef = useRef();

  const classes = useStyles();
  const {
    entity,
    placeholder,
    searchConfig,
    filterConfigs,
    savedSearches,
    onFiltersChanged,
    header,
    footer,
    exportPath,
  } = props;
  const [filterValues, setFilterValues] = useState(props.filterValues ?? []);

  useEffect(() => {
    setFilterValues(props.filterValues ?? []);
  }, [props.filterValues]);

  const [editingFilterIndex, setEditingFilterIndex] = useState((null: ?number));
  const [isInputFocused, setIsInputFocused] = useState(false);
  const [bookmark, setBookmark] = useState<?Bookmark>(null);

  const onFilterValueChanged = (index: number, filterValue: FilterValue) => {
    setBookmark(null);
    setFilterValues([
      ...filterValues.slice(0, index),
      filterValue,
      ...filterValues.slice(index + 1),
    ]);
  };

  const removeFilter = (index: number) => {
    props.onFilterRemoved && props.onFilterRemoved(filterValues[index]);
    const newFilterValues = update(filterValues, {
      $splice: [[index, 1]],
    });
    setBookmark(null);
    setFilterValues(newFilterValues);
    onFiltersChanged(newFilterValues);
    setEditingFilterIndex(null);
  };

  const onFilterBlurred = (index: number, filterValue: FilterValue) => {
    filtersTypeaheadRef.current && filtersTypeaheadRef.current.focus();

    if (!doesFilterHasValue(filterValue)) {
      removeFilter(index);
      return;
    }

    props.onFilterBlurred && props.onFilterBlurred(filterValue);
    setEditingFilterIndex(null);
    const newFilterValues = update(filterValues, {
      [index]: {$set: filterValue},
    });
    setBookmark(null);
    setFilterValues(newFilterValues);
    onFiltersChanged(newFilterValues);
  };

  const savedSearch = React.useContext(AppContext).isFeatureEnabled(
    'saved_searches',
  );

  return (
    <PowerSearchContext.Provider
      value={{
        bookmark,
        setBookmark,
      }}>
      <div className={classNames(classes.root, props.className)}>
        <div className={classes.headerContainer}>
          {header != null && header}
        </div>
        <div className={classes.searchBarContainer}>
          <div className={classes.searchTypeahead}>
            {filterValues.length > 0 || isInputFocused ? null : (
              <Text variant="body2" className={classes.placeholder}>
                {placeholder}
              </Text>
            )}
            {filterValues.map((filterValue, i) => {
              const filterConfig = filterConfigs.find(
                filter => filter.key === filterValue.key,
              );
              if (filterConfig == null || filterConfig.component == null) {
                return null;
              }
              const FilterComponent = filterConfig.component;
              return (
                <div className={classes.filter} key={filterValue.id}>
                  <FilterComponent
                    config={filterConfig}
                    editMode={editingFilterIndex === i}
                    value={filterValue}
                    onInputBlurred={() => {
                      return onFilterBlurred(i, filterValue);
                    }}
                    onNewInputBlurred={value => onFilterBlurred(i, value)}
                    onValueChanged={value => onFilterValueChanged(i, value)}
                    onRemoveFilter={() => removeFilter(i)}
                  />
                </div>
              );
            })}
          </div>
          <div className={classes.typeahead}>
            <FiltersTypeahead
              ref={filtersTypeaheadRef}
              options={filterConfigs}
              searchConfig={searchConfig}
              savedSearches={savedSearches ?? []}
              selectedFilters={filterValues.map(filter => filter.name)}
              onFilterSelected={({option, optionType}) => {
                if (optionType == 'SAVED_SEARCH') {
                  const searchConfig = savedSearches?.find(
                    filter => filter.key === option.key,
                  );
                  if (searchConfig == null) {
                    return null;
                  }
                  setBookmark({id: searchConfig.id, name: searchConfig.label});
                  onFiltersChanged(configToFilterQuery(searchConfig));
                } else {
                  const filterConfig = filterConfigs.find(
                    filter => filter.key === option.key,
                  );
                  if (filterConfig == null) {
                    return null;
                  }
                  setEditingFilterIndex(filterValues.length);
                  setFilterValues([
                    ...filterValues,
                    props.getSelectedFilter(filterConfig),
                  ]);
                }
                setIsInputFocused(false);
              }}
              onInputFocused={() => setIsInputFocused(true)}
              onInputBlurred={() => setIsInputFocused(false)}
            />
          </div>
          {footer != null && (
            <Text variant="body2" className={classes.footer}>
              {footer}
            </Text>
          )}
          {savedSearch && entity && (
            <FilterBookmark filters={filterValues} entity={entity} />
          )}
          {exportPath && (
            <CSVFileExport
              title="Export"
              exportPath={exportPath}
              filters={filterValues}
            />
          )}
        </div>
      </div>
    </PowerSearchContext.Provider>
  );
};

export default PowerSearchBar;
