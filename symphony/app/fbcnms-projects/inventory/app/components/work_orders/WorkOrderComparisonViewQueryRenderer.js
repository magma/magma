/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import ComparisonViewNoResults from '../comparison_view/ComparisonViewNoResults';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import WorkOrdersMap from './WorkOrdersMap';
import WorkOrdersView from './WorkOrdersView';
import classNames from 'classnames';
import {DisplayOptions} from '../InventoryViewContainer';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

import type {DisplayOptionTypes} from '../InventoryViewContainer';

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
}));

type Props = {
  className?: string,
  onWorkOrderSelected: (workOrderId: string) => void,
  limit?: number,
  filters: Array<any>,
  workOrderKey: number,
  displayMode?: DisplayOptionTypes,
  onQueryReturn?: (resultCount: number) => void,
};

const workOrderSearchQuery = graphql`
  query WorkOrderComparisonViewQueryRendererSearchQuery(
    $limit: Int
    $filters: [WorkOrderFilterInput!]!
  ) {
    workOrderSearch(limit: $limit, filters: $filters) {
      count
      workOrders {
        ...WorkOrdersView_workOrder
        ...WorkOrdersMap_workOrders
      }
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
    onQueryReturn,
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
          stringSet: f.stringSet,
        })),
        workOrderKey: workOrderKey,
      }}
      render={props => {
        const {count, workOrders} = props.workOrderSearch;
        onQueryReturn && onQueryReturn(count);
        if (count === 0) {
          return <ComparisonViewNoResults />;
        }
        return (
          <div className={classNames(classes.root, className)}>
            {displayMode === DisplayOptions.map ? (
              <WorkOrdersMap workOrders={workOrders} />
            ) : (
              <WorkOrdersView
                workOrder={workOrders}
                onWorkOrderSelected={onWorkOrderSelected}
              />
            )}
          </div>
        );
      }}
    />
  );
};

export default WorkOrderComparisonViewQueryRenderer;
