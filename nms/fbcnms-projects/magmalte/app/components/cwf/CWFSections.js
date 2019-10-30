/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {SectionsConfigs} from '@fbcnms/magmalte/app/components/layout/Section';

import CWFGateways from './CWFGateways';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import React from 'react';

export function getCWFSections(): SectionsConfigs {
  const sections = [
    {
      path: 'gateways',
      label: 'Gateways',
      icon: <CellWifiIcon />,
      component: CWFGateways,
    },
  ];

  return [
    'gateways', // landing path
    sections,
  ];
}
