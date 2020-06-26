/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

// TODO: Remove unnecessary styles and resolve overflow

import type {lte_gateway} from '@fbcnms/magma-api';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React from 'react';

import {colors} from '../../theme/default';
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

export default function GatewaySummary({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const version = gwInfo.status?.platform_info?.packages?.[0]?.version;
  return (
    <Card elevation={0}>
      <Grid
        container
        zeroMinWidth
        xs={12}
        direction="column"
        className={classes.kpiBlock}>
        <CardHeader
          className={classes.kpiBox}
          title={gwInfo.description}
          titleTypographyProps={{
            variant: 'body2',
            className: classes.kpiValue,
            title: gwInfo.description,
          }}
        />
        <CardHeader
          title="Gateway ID"
          subheader={gwInfo.id}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Gateway ID',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: gwInfo.id,
          }}
        />

        <CardHeader
          title="Hardware UUID"
          subheader={gwInfo.device.hardware_id}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Hardware UUID',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: gwInfo.device.hardware_id,
          }}
        />

        <CardHeader
          title="Version"
          subheader={version ?? 'null'}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Version',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: version ?? 'null',
          }}
        />
      </Grid>
    </Card>
  );
}
