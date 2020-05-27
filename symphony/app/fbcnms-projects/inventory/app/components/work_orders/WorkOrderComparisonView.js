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
import Button from '@fbcnms/ui/components/design-system/Button';
import ErrorBoundary from '@fbcnms/ui/components/ErrorBoundary/ErrorBoundary';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import InventorySuspense from '../../common/InventorySuspense';
import InventoryView, {DisplayOptions} from '../InventoryViewContainer';
import PowerSearchBar from '../power_search/PowerSearchBar';
import React, {useMemo, useState} from 'react';
import WorkOrderCard from './WorkOrderCard';
import WorkOrderComparisonViewQueryRenderer from './WorkOrderComparisonViewQueryRenderer';
import fbt from 'fbt';
import useFilterBookmarks from '../comparison_view/hooks/filterBookmarksHook';
import useLocationTypes from '../comparison_view/hooks/locationTypesHook';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {WorkOrderSearchConfig} from './WorkOrderSearchConfig';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {getInitialFilterValue} from '../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
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
}));

const QUERY_LIMIT = 100;

const WorkOrderComparisonView = () => {
  const [filters, setFilters] = useState([]);
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [workOrderKey, setWorkOrderKey] = useState(1);
  const [resultsDisplayMode, setResultsDisplayMode] = useState(
    DisplayOptions.table,
  );
  const [count, setCount] = useState((0: number));
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
  const filterBookmarksFilterConfig = useFilterBookmarks('WORK_ORDER');

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
        <InventorySuspense>
          <AddWorkOrderCard workOrderTypeId={selectedWorkOrderTypeId} />
        </InventorySuspense>
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

  const header = {
    title: 'Work Orders',
    searchBar: (
      <div className={classes.powerSearchBarWrapper}>
        <PowerSearchBar
          placeholder="Filter work orders"
          className={classes.powerSearchBar}
          filterConfigs={filterConfigs}
          searchConfig={WorkOrderSearchConfig}
          filterValues={filters}
          savedSearches={filterBookmarksFilterConfig}
          getSelectedFilter={(filterConfig: FilterConfig) =>
            getInitialFilterValue(
              filterConfig.key,
              filterConfig.name,
              filterConfig.defaultOperator,
              null,
            )
          }
          onFiltersChanged={filters => setFilters(filters)}
          exportPath={'/work_orders'}
          entity={'WORK_ORDER'}
          footer={
            count !== 0
              ? count > QUERY_LIMIT
                ? fbt(
                    '1 to ' +
                      fbt.param('size of page', QUERY_LIMIT) +
                      ' of ' +
                      fbt.param('total number possible rows', count),
                    'header to indicate partial results',
                  )
                : fbt(
                    '1 to ' + fbt.param('number of results in page', count),
                    'header to indicate number of results',
                  )
              : null
          }
        />
      </div>
    ),
    actionButtons: [
      <FormActionWithPermissions
        permissions={{
          entity: 'workorder',
          action: 'create',
        }}>
        <Button onClick={showDialog}>
          <fbt desc="">Create Work Order</fbt>
        </Button>
      </FormActionWithPermissions>,
    ],
  };

  return (
    <ErrorBoundary>
      <InventoryView
        permissions={{entity: 'workorder'}}
        header={header}
        onViewToggleClicked={setResultsDisplayMode}>
        <WorkOrderComparisonViewQueryRenderer
          limit={QUERY_LIMIT}
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
          onQueryReturn={c => setCount(c)}
        />
        <AddWorkOrderDialog
          key={`new_work_order_${dialogKey}`}
          open={dialogOpen}
          onClose={hideDialog}
          onWorkOrderTypeSelected={typeId => {
            navigateToAddWorkOrder(typeId);
            setDialogOpen(false);
          }}
        />
      </InventoryView>
    </ErrorBoundary>
  );
};

export default WorkOrderComparisonView;
