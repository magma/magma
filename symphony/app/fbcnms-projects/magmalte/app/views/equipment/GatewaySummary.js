/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';

import React from 'react';
import type {lte_gateway} from '@fbcnms/magma-api';

export default function GatewaySummary({gw_info}: {gw_info: lte_gateway}) {
  if (!gw_info) {
    return null;
  }

  const version = gw_info?.status?.platform_info?.packages?.[0]?.version;
  return (
    <>
      <Card>
        <CardHeader
          title={gw_info.description}
          titleTypographyProps={{variant: 'body2'}}
        />
      </Card>
      <Card>
        <CardHeader
          title="Gateway ID"
          subheader={gw_info.id}
          titleTypographyProps={{variant: 'caption'}}
          subheaderTypographyProps={{variant: 'body2'}}
        />
      </Card>

      <Card>
        <CardHeader
          title="Hardware UUID"
          subheader={gw_info.device.hardware_id}
          titleTypographyProps={{variant: 'caption'}}
          subheaderTypographyProps={{variant: 'body2'}}
        />
      </Card>

      <Card>
        <CardHeader
          title="Version"
          subheader={version ?? 'null'}
          titleTypographyProps={{variant: 'caption'}}
          subheaderTypographyProps={{variant: 'body2'}}
        />
      </Card>
    </>
  );
}
