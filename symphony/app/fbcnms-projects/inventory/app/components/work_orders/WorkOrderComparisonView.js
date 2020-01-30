/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterConfig} from '../comparison_view/ComparisonViewTypes';

import AddWorkOrderCard from './AddWorkOrderCard';
import AddWorkOrderDialog from './AddWorkOrderDialog';
import ErrorBoundary from '@fbcnms/ui/components/ErrorBoundary/ErrorBoundary';
import InventoryViewHeader, {DisplayOptions} from '../InventoryViewHeader';
import PowerSearchBar from '../power_search/PowerSearchBar';
import React, {useMemo, useState} from 'react';
import WorkOrderCard from './WorkOrderCard';
import WorkOrderComparisonViewQueryRenderer from './WorkOrderComparisonViewQueryRenderer';
import classNames from 'classnames';
import useLocationTypes from '../comparison_view/hooks/locationTypesHook';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {WorkOrderSearchConfig} from './WorkOrderSearchConfig';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {getInitialFilterValue} from '../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
  powerSearchBarWrapper: {
    paddingRight: '8px',
  },
  powerSearchBar: {
    borderRadius: '8px',
  },
  searchResults: {
    flexGrow: 1,
    paddingTop: '8px',
  },
});

const WorkOrderComparisonView = () => {
  const [filters, setFilters] = useState([]);
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [workOrderKey, setWorkOrderKey] = useState(1);
  const [resultsDisplayMode, setResultsDisplayMode] = useState(
    DisplayOptions.table,
  );
  const {match, history, location} = useRouter();
  const classes = useStyles();

  const selectedWorkOrderTypeId = useMemo(
    () => extractEntityIdFromUrl('workorderType', location.search),
    [location.search],
  );

  const selectedWorkOrderCardId = useMemo(
    () => extractEntityIdFromUrl('workorder', location.search),
    [location],
  );

  const locationTypesFilterConfigs = useLocationTypes();

  const filterConfigs = WorkOrderSearchConfig.map(ent => ent.filters)
    .reduce((allFilters, currentFilter) => allFilters.concat(currentFilter), [])
    .concat(locationTypesFilterConfigs ?? []);

  function navigateToAddWorkOrder(selectedWorkOrderTypeId: ?string) {
    history.push(
      match.url +
        (selectedWorkOrderTypeId
          ? `?workorderType=${selectedWorkOrderTypeId}`
          : ''),
    );
  }

  function navigateToWorkOrder(selectedWorkOrderCardId: ?string) {
    history.push(InventoryAPIUrls.workorder(selectedWorkOrderCardId));
  }

  const showDialog = () => {
    setDialogOpen(true);
    setDialogKey(dialogKey + 1);
    setWorkOrderKey(workOrderKey + 1);
  };

  const hideDialog = () => setDialogOpen(false);

  if (selectedWorkOrderTypeId != null) {
    return (
      <ErrorBoundary>
        <AddWorkOrderCard workOrderTypeId={selectedWorkOrderTypeId} />
      </ErrorBoundary>
    );
  }

  if (selectedWorkOrderCardId != null) {
    return (
      <ErrorBoundary>
        <WorkOrderCard
          workOrderId={selectedWorkOrderCardId}
          onWorkOrderExecuted={() => {}}
          onWorkOrderRemoved={() => navigateToWorkOrder(null)}
        />
      </ErrorBoundary>
    );
  }

  return (
    <ErrorBoundary>
      <div className={classes.root}>
        <InventoryViewHeader
          title="Work Orders"
          onViewToggleClicked={setResultsDisplayMode}
          searchBar={
            <div className={classes.powerSearchBarWrapper}>
              <PowerSearchBar
                placeholder="Filter work orders"
                className={classes.powerSearchBar}
                filterConfigs={filterConfigs}
                searchConfig={WorkOrderSearchConfig}
                getSelectedFilter={(filterConfig: FilterConfig) =>
                  getInitialFilterValue(
                    filterConfig.key,
                    filterConfig.name,
                    filterConfig.defaultOperator,
                    null,
                  )
                }
                onFiltersChanged={filters => setFilters(filters)}
              />
            </div>
          }
          actionButtons={[
            {
              title: 'Add Work Order',
              action: showDialog,
            },
          ]}
        />
        <div className={classes.searchResults}>
          <WorkOrderComparisonViewQueryRenderer
            className={classNames({
              [classes.comparisionViewTable]:
                resultsDisplayMode === DisplayOptions.table,
            })}
            limit={50}
            filters={filters}
            onWorkOrderSelected={selectedWorkOrderCardId =>
              navigateToWorkOrder(selectedWorkOrderCardId)
            }
            workOrderKey={workOrderKey}
            displayMode={
              resultsDisplayMode === DisplayOptions.map
                ? DisplayOptions.map
                : DisplayOptions.table
            }
          />
        </div>
      </div>
      <AddWorkOrderDialog
        key={`new_work_order_${dialogKey}`}
        open={dialogOpen}
        onClose={hideDialog}
        onWorkOrderTypeSelected={typeId => {
          navigateToAddWorkOrder(typeId);
          setDialogOpen(false);
        }}
      />
    </ErrorBoundary>
  );
};

export default WorkOrderComparisonView;
