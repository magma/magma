/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import SearchIcon from '@material-ui/icons/Search';
import Text from '@fbcnms/ui/components/design-system/Text';
import WorkOrdersMap from './WorkOrdersMap';
import WorkOrdersView from './WorkOrdersView';
import classNames from 'classnames';
import {DisplayOptions} from '../InventoryViewHeader';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

import type {DisplayOptionTypes} from '../InventoryViewHeader';

const useStyles = makeStyles(theme => ({
  root: {
    height: '100%',
    width: '100%',
    flexGrow: 1,
  },
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
  tableViewContainer: {
    paddingRight: '24px',
    paddingLeft: '24px',
  },
}));

type Props = {
  className?: string,
  onWorkOrderSelected: (workOrderId: string) => void,
  limit?: number,
  filters: Array<any>,
  workOrderKey: number,
  displayMode?: DisplayOptionTypes,
};

const workOrderSearchQuery = graphql`
  query WorkOrderComparisonViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [WorkOrderFilterInput!]!
  ) {
    workOrderSearch(limit: $limit, filters: $filters) {
      ...WorkOrdersView_workOrder
      ...WorkOrdersMap_workOrders
    }
  }
`;

const WorkOrderComparisonViewQueryRenderer = (props: Props) => {
  const classes = useStyles();
  const {
    filters,
    limit,
    onWorkOrderSelected,
    workOrderKey,
    displayMode,
    className,
  } = props;

  return (
    <InventoryQueryRenderer
      query={workOrderSearchQuery}
      variables={{
        limit: limit,
        filters: filters.map(f => ({
          filterType: f.name.toUpperCase(),
          operator: f.operator.toUpperCase(),
          stringValue: f.stringValue,
          propertyValue: f.propertyValue,
          idSet: f.idSet,
        })),
        workOrderKey: workOrderKey,
      }}
      render={props => {
        const {workOrderSearch} = props;
        if (!workOrderSearch || workOrderSearch.length === 0) {
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
          <div className={classNames(classes.root, className)}>
            {displayMode === DisplayOptions.map ? (
              <WorkOrdersMap workOrders={workOrderSearch} />
            ) : (
              <WorkOrdersView
                workOrder={workOrderSearch}
                onWorkOrderSelected={onWorkOrderSelected}
                className={classes.tableViewContainer}
              />
            )}
          </div>
        );
      }}
    />
  );
};

export default WorkOrderComparisonViewQueryRenderer;
