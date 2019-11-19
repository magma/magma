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
import type {FiringAlarm, Labels} from '@fbcnms/alarms/components/AlarmAPIType';

export default function Alarms() {
  const {isFeatureEnabled} = React.useContext(AppContext);
  const experimentalAlertsEnabled = isFeatureEnabled('alerts_experimental');
  return (
    <FBCAlarms
      apiUtil={MagmaAlarmsApiUtil}
      makeTabLink={({match, keyName}) =>
        `/nms/${match.params.networkId || ''}/alerts/${keyName}`
      }
      experimentalTabsEnabled={experimentalAlertsEnabled}
      filterLabels={filterSymphonyLabels}
    />
  );
}

/**
 * Filters out hidden system labels from the firing alerts table
 */
function filterSymphonyLabels(labels: Labels, _alarm: FiringAlarm) {
  const labelsToFilter = ['monitor', 'networkID'];
  const filtered = {...labels};
  for (const label of labelsToFilter) {
    delete filtered[label];
  }
  return filtered;
}
