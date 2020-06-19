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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import React from 'react';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';

export function EnodebSummary({enbInfo}: {enbInfo: EnodebInfo}) {
  return (
    <Card variant={'outlined'}>
      <CardHeader
        data-testid="eNodeB Serial Number"
        title="eNodeB Serial Number"
        subheader={enbInfo.enb.serial}
        titleTypographyProps={{variant: 'caption'}}
        subheaderTypographyProps={{variant: 'body1'}}
      />
    </Card>
  );
}

// Status Indicator displays a small text with an DeviceStatusCircle icon
// disabled indicates if the status color is to be grayed out
// up/down indicates if we have to display status to be in green or in red
function StatusIndicator(disabled: boolean, up: boolean, val: string) {
  return (
    <Grid container>
      <Grid item>
        <DeviceStatusCircle isGrey={disabled} isActive={up} isFilled={true} />
      </Grid>
      <Grid item>{val}</Grid>
    </Grid>
  );
}

export function EnodebStatus({enbInfo}: {enbInfo: EnodebInfo}) {
  const isEnbHealthy = isEnodebHealthy(enbInfo);
  return (
    <Grid container>
      <Grid item xs={6}>
        <Card variant={'outlined'}>
          <CardHeader
            titleTypographyProps={{variant: 'caption'}}
            subheaderTypographyProps={{variant: 'body1'}}
            data-testid="Health"
            title="Health"
            subheader={StatusIndicator(
              false,
              isEnbHealthy,
              isEnbHealthy ? 'Good' : 'Bad',
            )}
          />
        </Card>
      </Grid>
      <Grid item xs={6}>
        <Card variant={'outlined'}>
          <CardHeader
            titleTypographyProps={{variant: 'caption'}}
            subheaderTypographyProps={{variant: 'body1'}}
            data-testid="Transmit Enabled"
            title="Transmit Enabled"
            subheader={StatusIndicator(
              !enbInfo.enb.config.transmit_enabled,
              enbInfo.enb.config.transmit_enabled,
              enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled',
            )}
          />
        </Card>
      </Grid>

      <Grid item xs={6}>
        <Card variant={'outlined'}>
          <CardHeader
            titleTypographyProps={{variant: 'caption'}}
            subheaderTypographyProps={{variant: 'body1'}}
            data-testid="Gateway ID"
            title="Gateway ID"
            subheader={StatusIndicator(
              false,
              enbInfo.enb_state.enodeb_connected,
              enbInfo.enb_state.reporting_gateway_id ?? '',
            )}
          />
        </Card>
      </Grid>

      <Grid item xs={6}>
        <Card variant={'outlined'}>
          <CardHeader
            titleTypographyProps={{variant: 'caption'}}
            subheaderTypographyProps={{variant: 'body1'}}
            data-testid="Mme Connected"
            title="Mme Connected"
            subheader={StatusIndicator(
              false,
              enbInfo.enb_state.mme_connected,
              enbInfo.enb_state.mme_connected ? 'Connected' : 'Disconnected',
            )}
          />
        </Card>
      </Grid>
    </Grid>
  );
}
