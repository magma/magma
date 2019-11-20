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
import {graphql} from 'relay-runtime';
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
  workOrdersTable: {
    margin: '24px',
  },
}));

type Props = {
  onWorkOrderSelected: (workOrderId: string) => void,
  limit?: number,
  filters: Array<any>,
  workOrderKey: number,
  resultsDisplayMode: ?'map' | 'table',
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
    resultsDisplayMode,
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
          <>
            {resultsDisplayMode === 'map' ? (
              <WorkOrdersMap workOrders={workOrderSearch} />
            ) : (
              <WorkOrdersView
                className={classes.workOrdersTable}
                workOrder={workOrderSearch}
                onWorkOrderSelected={onWorkOrderSelected}
              />
            )}
          </>
        );
      }}
    />
  );
};

export default WorkOrderComparisonViewQueryRenderer;
