/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContext from '@fbcnms/ui/context/AppContext';
import EquipmentComparisonViewQueryRenderer from './EquipmentComparisonViewQueryRenderer';
import InventoryErrorBoundary from '../../common/InventoryErrorBoundary';
import React, {useContext} from 'react';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {makeStyles} from '@material-ui/styles';

const QUERY_LIMIT = 50;

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
  },
}));

const EquipmentComparisonView = () => {
  const {history} = useRouter();
  const classes = useStyles();
  const equipmentExportEnabled = useContext(AppContext).isFeatureEnabled(
    'equipment_export',
  );
  return (
    <InventoryErrorBoundary>
      <div className={classes.root}>
        <div className={classes.searchResults}>
          <EquipmentComparisonViewQueryRenderer
            limit={QUERY_LIMIT}
            showExport={equipmentExportEnabled}
            onEquipmentSelected={equipment => {
              ServerLogger.info(
                LogEvents.EQUIPMENT_COMPARISON_VIEW_EQUIPMENT_CLICKED,
              );
              history.replace(`inventory?equipment=${equipment.id}`);
            }}
          />
        </div>
      </div>
    </InventoryErrorBoundary>
  );
};

export default EquipmentComparisonView;
