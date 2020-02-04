/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {EntityType, FiltersQuery} from './ComparisonViewTypes';

import Divider from '@material-ui/core/Divider';
import InventoryComparisonViewRouter from './InventoryComparisonViewRouter';
import InventoryErrorBoundary from '../../common/InventoryErrorBoundary';
import PowerSearchBarRouter from './PowerSearchBarRouter';
import PowerSearchFilterSubjectDropDown from '../power_search/PowerSearchFilterSubjectDropDown';
import React, {useState} from 'react';
import {EntityTypeMap} from './ComparisonViewTypes';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {useRouter} from '@fbcnms/ui/hooks';

import {makeStyles} from '@material-ui/styles';

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
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
}));

const QUERY_LIMIT = 100;

const InventoryComparisonView = () => {
  const classes = useStyles();
  const [count, setCount] = useState((0: number));

  const onSubjectChange = subject => {
    ServerLogger.info(LogEvents.COMPARISON_VIEW_SUBJECT_CHANGED, {
      subject: subject,
    });
    // add subject to url
    let path = history.location.pathname;
    if (
      path.endsWith(InventoryAPIUrls.search) ||
      path.endsWith(InventoryAPIUrls.search + '/')
    ) {
      if (subject) {
        path = path.replace(/\/$/, '');
        history.push(path + '/' + subject);
      }
    }
    history.replace({pathname: subject});
  };

  const onFiltersChange = filters => {
    ServerLogger.info(LogEvents.COMPARISON_VIEW_FILTERS_CHANGED, {
      filters: filters,
    });
    // add filters to URL
    const filtersStr =
      filters.length > 0 ? `filters=${JSON.stringify(filters)}` : '';
    if (getSubject()) {
      history.replace({
        search: filtersStr,
      });
    }
  };

  const getFiltersFromURL = (): FiltersQuery => {
    if (getSubject()) {
      const urlParams = new URLSearchParams(history.location.search);
      const filtersStr = urlParams.get('filters') ?? '[]';
      return JSON.parse(filtersStr);
    }
    return [];
  };

  const getSubjectFromURL = (): EntityType => {
    const subj = getSubject();
    if (subj) {
      return subj;
    }
    onSubjectChange('equipment');
    return 'equipment';
  };

  const getSubject = (): ?EntityType => {
    const path = history.location.pathname;
    const subj = path.split(InventoryAPIUrls.search + '/')[1];
    if (Object.keys(EntityTypeMap).includes(subj)) {
      return subj;
    }
    return null;
  };

  const {history} = useRouter();
  const filters = getFiltersFromURL();
  const subject = getSubjectFromURL();
  return (
    <InventoryErrorBoundary>
      <div className={classes.root}>
        <div className={classes.searchResults}>
          <div className={classes.root}>
            <div className={classes.searchBar}>
              <PowerSearchFilterSubjectDropDown
                subject={subject}
                onSubjectChange={onSubjectChange}
              />
              <Divider orientation="vertical" />
              <PowerSearchBarRouter
                filters={filters}
                subject={subject}
                onFiltersChanged={onFiltersChange}
                footer={
                  count != null
                    ? count > QUERY_LIMIT
                      ? `1 to ${QUERY_LIMIT} of ${count}`
                      : `1 to ${count}`
                    : null
                }
              />
            </div>
            <InventoryComparisonViewRouter
              filters={filters}
              limit={QUERY_LIMIT}
              onQueryReturn={x => setCount(x)}
              subject={subject}
            />
          </div>
        </div>
      </div>
    </InventoryErrorBoundary>
  );
};

export default InventoryComparisonView;
