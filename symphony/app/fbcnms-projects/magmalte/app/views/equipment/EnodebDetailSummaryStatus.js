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
import type {KPIRows} from '../../components/KPIGrid';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import KPIGrid from '../../components/KPIGrid';
import React from 'react';

import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: props => (props.hasStatus ? 'calc(100% - 16px)' : '100%'),
  },
  kpiValue: {
    color: colors.primary.brightGray,
  },
  kpiBox: {
    width: '100%',
  },
}));

export function EnodebSummary({enbInfo}: {enbInfo: EnodebInfo}) {
  const classes = useStyles();
  return (
    <Card elevation={0}>
      <CardHeader
        className={classes.kpiBox}
        data-testid="eNodeB Serial Number"
        title="eNodeB Serial Number"
        subheader={enbInfo.enb.serial}
        titleTypographyProps={{
          variant: 'body3',
          className: classes.kpiLabel,
          title: 'eNodeB Serial Number',
        }}
        subheaderTypographyProps={{
          variant: 'body1',
          className: classes.kpiValue,
          title: enbInfo.enb.serial,
        }}
      />
    </Card>
  );
}

export function EnodebStatus({enbInfo}: {enbInfo: EnodebInfo}) {
  const classes = useStyles();
  const isEnbHealthy = isEnodebHealthy(enbInfo);

  const kpiData: KPIRows[] = [
    [
      {
        category: 'Health',
        value: isEnbHealthy ? 'Good' : 'Bad',
        statusCircle: true,
        status: isEnbHealthy,
      },
      {
        category: 'Transmit Enabled',
        value: enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: enbInfo.enb.config.transmit_enabled,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: enbInfo.enb_state.reporting_gateway_id ?? '',
        statusCircle: true,
        status: enbInfo.enb_state.enodeb_connected,
      },
      {
        category: 'Mme Connected',
        value: enbInfo.enb_state.mme_connected ? 'Connected' : 'Disconnected',
        statusCircle: false,
        status: enbInfo.enb_state.mme_connected,
      },
    ],
  ];
  return <KPIGrid data={kpiData} />;
}
