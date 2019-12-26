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
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapIcon from '@material-ui/icons/Map';
import PowerSearchBar from '../power_search/PowerSearchBar';
import React, {useMemo, useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import WorkOrderCard from './WorkOrderCard';
import WorkOrderComparisonViewQueryRenderer from './WorkOrderComparisonViewQueryRenderer';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import useLocationTypes from '../comparison_view/hooks/locationTypesHook';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {InventoryAPIUrls} from '../../common/InventoryAPI';
import {WorkOrderSearchConfig} from './WorkOrderSearchConfig';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {getInitialFilterValue} from '../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardRoot: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    paddingLeft: '0px',
    paddingRight: '0px',
  },
  cardContent: {
    paddingLeft: '0px',
    paddingRight: '0px',
    paddingTop: '0px',
    flexGrow: 1,
    width: '100%',
    padding: '0px',
    backgroundColor: symphony.palette.background,
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    height: '100%',
  },
  searchResults: {
    display: 'flex',
    flexDirection: 'column',
    flexGrow: 1,
    backgroundColor: symphony.palette.background,
  },
  bar: {
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
  searchBar: {
    flexGrow: 1,
  },
  buttonContent: {
    paddingTop: '4px',
  },
  addWorkOrderButton: {
    alignSelf: 'flex-end',
  },
  comparisionViewTable: {
    margin: '0px 32px',
  },
  titleContainer: {
    margin: '32px',
    display: 'flex',
  },
  title: {
    flexGrow: 1,
    display: 'block',
  },
}));

const WorkOrderComparisonView = () => {
  const [filters, setFilters] = useState([]);
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [workOrderKey, setWorkOrderKey] = useState(1);
  const [resultsDisplayMode, setResultsDisplayMode] = useState('table');
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
    .concat(locationTypesFilterConfigs);

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
      <div className={classes.cardRoot}>
        <div className={classes.root}>
          <div className={classes.bar}>
            <div className={classes.searchBar}>
              <PowerSearchBar
                placeholder="Filter work orders"
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
            <MapButtonGroup
              initiallySelectedButton={resultsDisplayMode === 'table' ? 0 : 1}
              onIconClicked={id => {
                setResultsDisplayMode(id === 'table' ? 'table' : 'map');
              }}
              buttons={[
                {
                  item: <ListAltIcon className={classes.buttonContent} />,
                  id: 'table',
                },
                {
                  item: <MapIcon className={classes.buttonContent} />,
                  id: 'map',
                },
              ]}
            />
          </div>
          <div className={classes.searchResults}>
            <div className={classes.titleContainer}>
              <Text className={classes.title} variant="h6">
                Work Orders
              </Text>
              <Button
                className={classes.addWorkOrderButton}
                onClick={showDialog}>
                Add Work Order
              </Button>
            </div>
            <WorkOrderComparisonViewQueryRenderer
              className={classNames({
                [classes.comparisionViewTable]: resultsDisplayMode === 'table',
              })}
              limit={50}
              filters={filters}
              onWorkOrderSelected={selectedWorkOrderCardId =>
                navigateToWorkOrder(selectedWorkOrderCardId)
              }
              workOrderKey={workOrderKey}
              resultsDisplayMode={resultsDisplayMode}
            />
          </div>
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
