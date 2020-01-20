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
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';

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
  const [filters, setFilters] = useState(([]: FiltersQuery));
  const [count, setCount] = useState((0: number));
  const [subject, setSubject] = useState(('equipment': EntityType));

  const onSubjectChange = subject => {
    ServerLogger.info(LogEvents.COMPARISON_VIEW_SUBJECT_CHANGED, {
      subject: subject,
    });
    setFilters([]);
    setSubject(subject);
  };

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
                onFiltersChanged={setFilters}
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
