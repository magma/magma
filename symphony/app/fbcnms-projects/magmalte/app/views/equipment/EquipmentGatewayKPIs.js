/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import KPITray from '../../components/KPITray';
import React from 'react';
import type {KPIData} from '../../components/KPITray';

export default function EquipmentGatewayKPIs() {
  const kpiData: KPIData[] = [
    {category: 'KPI1', value: 0},
    {category: 'KPI2', value: 0},
    {category: 'KPI3', value: 0},
    {category: 'KPI4', value: 0},
    {category: 'KPI5', value: 0},
    {category: 'KPI6', value: 0},
  ];
  return <KPITray data={kpiData} />;
}
