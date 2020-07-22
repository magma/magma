/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {KPIRows} from '../../components/KPIGrid';
import type {lte_gateway} from '@fbcnms/magma-api';

import Card from '@material-ui/core/Card';
import KPIGrid from '../../components/KPIGrid';
import React from 'react';

export default function GatewaySummary({gwInfo}: {gwInfo: lte_gateway}) {
  const version = gwInfo.status?.platform_info?.packages?.[0]?.version;

  const kpiData: KPIRows[] = [
    [
      {
        value: gwInfo.description,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: gwInfo.id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Hardware UUID',
        value: gwInfo.device.hardware_id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Version',
        value: version ?? 'null',
        statusCircle: false,
      },
    ],
  ];

  return (
    <Card elevation={0}>
      <KPIGrid data={kpiData} />
    </Card>
  );
}
