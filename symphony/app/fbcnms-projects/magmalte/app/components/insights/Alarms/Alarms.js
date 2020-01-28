/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import FBCAlarms from '@fbcnms/alarms/components/Alarms';
import {MagmaAlarmsApiUtil} from './AlarmApi';
import type {Labels} from '@fbcnms/alarms/components/AlarmAPIType';

export default function Alarms() {
  const {isFeatureEnabled} = React.useContext(AppContext);
  const disabledTabs = React.useMemo(
    () =>
      [
        isFeatureEnabled('alert_receivers') ? null : 'receivers',
        isFeatureEnabled('alert_routes') ? null : 'routes',
        isFeatureEnabled('alert_suppressions') ? null : 'suppressions',
      ].filter(Boolean),
    [isFeatureEnabled],
  );
  return (
    <FBCAlarms
      apiUtil={MagmaAlarmsApiUtil}
      makeTabLink={({match, keyName}) =>
        `/nms/${match.params.networkId || ''}/alerts/${keyName}`
      }
      disabledTabs={disabledTabs}
      thresholdEditorEnabled={true}
      filterLabels={filterSymphonyLabels}
    />
  );
}

/**
 * Filters out hidden system labels from the firing alerts table
 */
function filterSymphonyLabels(labels: Labels) {
  const labelsToFilter = ['monitor', 'networkID'];
  const filtered = {...labels};
  for (const label of labelsToFilter) {
    delete filtered[label];
  }
  return filtered;
}
