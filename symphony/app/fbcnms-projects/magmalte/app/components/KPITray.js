/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ComponentType} from 'react';

import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '../theme/design-system/Text';

import {colors} from '../theme/default';
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
    '& + &': {
      boxShadow: `-2px 0 0 ${colors.primary.concrete}`,
    },
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
  },
  kpiBox: {
    width: '100%',
    '& div': {
      width: '100%',
    },
  },
}));

export type KPIData = {
  category: string,
  value: number | string,
  unit?: string,
  icon?: ComponentType<SvgIconExports>,
};
type Props = {
  icon?: ComponentType<SvgIconExports>,
  description?: string,
  data: KPIData[],
};

export default function KPITray(props: Props) {
  const classes = useStyles();
  const kpiTray = [];

  if (props.icon) {
    const KpiIcon = props.icon;
    kpiTray.push(
      <Grid
        item
        alignItems="center"
        className={classes.kpiHeaderBlock}
        key="kpiTitle">
        <CardContent className={classes.kpiHeaderContent}>
          <KpiIcon className={classes.kpiHeaderIcon} />
          <Text variant="body1" className={classes.kpiTitle}>
            {props.description}
          </Text>
        </CardContent>
      </Grid>,
    );
  }

  kpiTray.push(
    props.data.map((kpi, i) => (
      <Grid
        container
        xs
        zeroMinWidth
        key={'data-' + i}
        className={classes.kpiBlock}
        alignItems="center">
        <CardHeader
          title={kpi.category}
          className={classes.kpiBox}
          subheader={kpi.value + (kpi.unit ?? '')}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: kpi.category,
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
          }}
          data-testid={kpi.category}
        />
      </Grid>
    )),
  );

  return (
    <Grid container zeroMinWidth>
      {kpiTray}
    </Grid>
  );
}
