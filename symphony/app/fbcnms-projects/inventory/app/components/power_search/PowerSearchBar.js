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
  EntityConfig,
  FilterConfig,
  FilterValue,
  FiltersQuery,
} from '../comparison_view/ComparisonViewTypes';

import * as React from 'react';
import CSVFileExport from '../CSVFileExport';
import FiltersTypeahead from '../comparison_view/FiltersTypeahead';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import {useRef, useState} from 'react';

import update from 'immutability-helper';
import {doesFilterHasValue} from '../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

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
  exportButton: {
    paddingLeft: '16px',
    paddingRight: '16px',
    marginLeft: 'auto',
  },
  exportButtonContainer: {
    display: 'flex',
  },
}));

type Props = {
  filterValues?: FiltersQuery,
  placeholder?: string,
  filterConfigs: Array<FilterConfig>,
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
};

const PowerSearchBar = (props: Props) => {
  const filtersTypeaheadRef = useRef();

  const classes = useStyles();
  const {
    placeholder,
    searchConfig,
    filterConfigs,
    onFiltersChanged,
    header,
    footer,
    exportPath,
  } = props;
  const [filterValues, setFilterValues] = useState(props.filterValues ?? []);

  const [editingFilterIndex, setEditingFilterIndex] = useState((null: ?number));
  const [isInputFocused, setIsInputFocused] = useState(false);

  const onFilterValueChanged = (index: number, filterValue: FilterValue) => {
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
    setFilterValues(newFilterValues);
    onFiltersChanged(newFilterValues);
  };
  return (
    <div className={classNames(classes.root, props.className)}>
      <div className={classes.headerContainer}>{header != null && header}</div>
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
            if (filterConfig == null) {
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
            selectedFilters={filterValues.map(filter => filter.name)}
            onFilterSelected={filterOption => {
              const filterConfig = filterConfigs.find(
                filter => filter.key === filterOption.key,
              );
              if (filterConfig == null) {
                return null;
              }
              setIsInputFocused(false);
              setEditingFilterIndex(filterValues.length);
              setFilterValues([
                ...filterValues,
                props.getSelectedFilter(filterConfig),
              ]);
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
        {exportPath && (
          <CSVFileExport
            title="Export"
            exportPath={exportPath}
            filters={filterValues}
          />
        )}
      </div>
    </div>
  );
};

export default PowerSearchBar;
