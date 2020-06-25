/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {KPIRows} from '../../components/KPIGrid';
import type {lte_gateway} from '@fbcnms/magma-api';

import KPIGrid from '../../components/KPIGrid';
import React from 'react';

import isGatewayHealthy from '../../components/GatewayUtils';
export default function GatewayDetailStatus({gwInfo}: {gwInfo: lte_gateway}) {
  let checkInTime = new Date(0);
  if (
    gwInfo.status &&
    gwInfo.status.checkin_time !== undefined &&
    gwInfo.status.checkin_time !== null
  ) {
    checkInTime = new Date(gwInfo.status.checkin_time);
  }

  const logAggregation =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes('td-agent-bit');

  const eventAggregation =
    !!gwInfo.magmad.dynamic_services &&
    gwInfo.magmad.dynamic_services.includes('eventd');

  const kpiData: KPIRows[] = [
    [
      {
        category: 'Health',
        value: isGatewayHealthy(gwInfo) ? 'Good' : 'Bad',
        statusCircle: true,
        status: isGatewayHealthy(gwInfo) ? 'Up' : 'Down',
      },
      {
        category: 'Last Check in',
        value: checkInTime.toLocaleString(),
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Event Aggregation',
        value: eventAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: eventAggregation ? 'Up' : 'Disabled',
      },
      {
        category: 'Log Aggregation',
        value: logAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: logAggregation ? 'Up' : 'Disabled',
      },
      {
        category: 'CPU Usage',
        value: '0',
        unit: '%',
        statusCircle: false,
      },
    ],
  ];
  return <KPIGrid data={kpiData} />;
}
