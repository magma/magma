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
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import Grid from '@material-ui/core/Grid';
import React from 'react';

import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  kpiHeaderBlock: {
    display: 'flex',
    alignItems: 'center',
    padding: 0,
  },
  kpiHeaderContent: {
    display: 'flex',
    alignItems: 'center',
  },
  kpiHeaderIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  kpiBlock: {
    boxShadow: `0 0 0 1px ${colors.primary.concrete}`,
  },
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

// Status Indicator displays a small text with an DeviceStatusCircle icon
// disabled indicates if the status color is to be grayed out
// up/down indicates if we have to display status to be in green or in red
function StatusIndicator(disabled: boolean, up: boolean, val: string) {
  const props = {hasStatus: true};
  const classes = useStyles(props);
  return (
    <Grid container zeroMinWidth alignItems="center" xs={12}>
      <Grid item>
        <DeviceStatusCircle isGrey={disabled} isActive={up} isFilled={true} />
      </Grid>
      <Grid item className={classes.kpiLabel}>
        {val}
      </Grid>
    </Grid>
  );
}

export function EnodebStatus({enbInfo}: {enbInfo: EnodebInfo}) {
  const classes = useStyles();
  const isEnbHealthy = isEnodebHealthy(enbInfo);
  return (
    <Card elevation={0}>
      <Grid container>
        <Grid
          container
          xs={6}
          zeroMinWidth
          className={classes.kpiBlock}
          alignItems="center">
          <CardHeader
            className={classes.kpiBox}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: 'Health',
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: isEnbHealthy ? 'Good' : 'Bad',
            }}
            data-testid="Health"
            title="Health"
            subheader={StatusIndicator(
              false,
              isEnbHealthy,
              isEnbHealthy ? 'Good' : 'Bad',
            )}
          />
        </Grid>
        <Grid
          container
          xs={6}
          zeroMinWidth
          className={classes.kpiBlock}
          alignItems="center">
          <CardHeader
            className={classes.kpiBox}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: 'Transmit Enabled',
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: enbInfo.enb.config.transmit_enabled
                ? 'Enabled'
                : 'Disabled',
            }}
            data-testid="Transmit Enabled"
            title="Transmit Enabled"
            subheader={StatusIndicator(
              !enbInfo.enb.config.transmit_enabled,
              enbInfo.enb.config.transmit_enabled,
              enbInfo.enb.config.transmit_enabled ? 'Enabled' : 'Disabled',
            )}
          />
        </Grid>
        <Grid
          container
          xs={6}
          zeroMinWidth
          className={classes.kpiBlock}
          alignItems="center">
          <CardHeader
            className={classes.kpiBox}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: 'Gateway ID',
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: enbInfo.enb_state.reporting_gateway_id ?? '',
            }}
            data-testid="Gateway ID"
            title="Gateway ID"
            subheader={StatusIndicator(
              false,
              enbInfo.enb_state.enodeb_connected,
              enbInfo.enb_state.reporting_gateway_id ?? '',
            )}
          />
        </Grid>
        <Grid
          container
          xs={6}
          zeroMinWidth
          className={classes.kpiBlock}
          alignItems="center">
          <CardHeader
            className={classes.kpiBox}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: 'Mme Connected',
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: enbInfo.enb_state.mme_connected
                ? 'Connected'
                : 'Disconnected',
            }}
            data-testid="Mme Connected"
            title="Mme Connected"
            subheader={StatusIndicator(
              false,
              enbInfo.enb_state.mme_connected,
              enbInfo.enb_state.mme_connected ? 'Connected' : 'Disconnected',
            )}
          />
        </Grid>
      </Grid>
    </Card>
  );
}
