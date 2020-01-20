/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EntityType} from './comparison_view/ComparisonViewTypes';

import AppContext from '@fbcnms/ui/context/AppContext';
import InventoryEntitiesTypeahead from '../components/InventoryEntitiesTypeahead';
import React, {useContext} from 'react';
import WorkOrdersPopover from './work_orders/WorkOrdersPopover';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    backgroundColor: theme.palette.common.white,
    borderBottom: '1px solid rgba(0, 0, 0, 0.1)',
    display: 'flex',
    padding: '0px 16px',
    width: '100%',
    height: '60px',
    alignItems: 'center',
  },
  locationSearch: {
    width: '270px',
  },
  spacer: {
    flexGrow: 1,
  },
}));

type Props = {
  onWorkOrderSelected: (workOrderId: ?string) => void,
  onSearchEntitySelected: (entityId: string, entityType: EntityType) => void,
  onNavigateToWorkOrder: (workOrderId: ?string) => void,
};

const InventoryTopBar = (props: Props) => {
  const {
    onWorkOrderSelected,
    onSearchEntitySelected,
    onNavigateToWorkOrder,
  } = props;
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);
  const woPopoverEnabled = isFeatureEnabled('planned_equipment');
  return (
    <div className={classes.root}>
      <div className={classes.locationSearch}>
        <InventoryEntitiesTypeahead onEntitySelected={onSearchEntitySelected} />
      </div>
      <div className={classes.spacer} />
      <div className={classes.workOrders}>
        {woPopoverEnabled && (
          <WorkOrdersPopover
            onSelect={onWorkOrderSelected}
            onNavigateToWorkOrder={onNavigateToWorkOrder}
          />
        )}
      </div>
    </div>
  );
};

export default InventoryTopBar;
