/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import BarChartIcon from '@material-ui/icons/BarChart';
import React from 'react';
import RhinoMetrics from './RhinoMetrics';

export function getRhinoSections(): SectionsConfigs {
  return [
    'metrics', // landing path
    [
      {
        path: 'metrics',
        label: 'Metrics',
        icon: <BarChartIcon />,
        component: RhinoMetrics,
      },
    ],
  ];
}
