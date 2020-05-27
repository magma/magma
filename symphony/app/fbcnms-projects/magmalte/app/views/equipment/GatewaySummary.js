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

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {useRouter} from '@fbcnms/ui/hooks';
import type {lte_gateway} from '@fbcnms/magma-api';

export default function GatewaySummary() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);

  const {response, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGatewaysByGatewayId,
    {
      networkId: networkId,
      gatewayId: gatewayId,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }
  if (response) {
    const gateway: lte_gateway = response;
    let version = 'null';
    if (
      gateway.status &&
      gateway.status.platform_info &&
      gateway.status.platform_info.packages &&
      gateway.status.platform_info.packages.length > 0
    ) {
      version = gateway.status.platform_info.packages[0].version;
    }

    if (gateway) {
      return (
        <>
          <Card>
            <CardHeader
              title={gateway.description}
              titleTypographyProps={{variant: 'body2'}}
            />
          </Card>
          <Card>
            <CardHeader
              title="Gateway ID"
              subheader={gateway.id}
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body2'}}
            />
          </Card>

          <Card>
            <CardHeader
              title="Hardware UUID"
              subheader={gateway.device.hardware_id}
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body2'}}
            />
          </Card>

          <Card>
            <CardHeader
              title="Version"
              subheader={version}
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body2'}}
            />
          </Card>
        </>
      );
    }
  }
  return null;
}
