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
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {makeStyles} from '@material-ui/styles';
import {gray7} from '@fbcnms/ui/theme/colors';
import type {ComponentType} from 'react';

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
    fill: '#545F77',
    marginRight: theme.spacing(1),
  },
  kpiBlock: {
    '& + &': {
      boxShadow: `-2px 0 0 ${gray7}`,
    },
  },
  kpiLabel: {
    color: '#545F77',
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: '#323845',
  },
  test: {
    width: '100%',

    '& div': {
      width: '100%',
    },
  },
}));

export type KPIData = {category: string, value: number | string, unit?: string};
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
          <Text variant="h6" className={classes.kpiTitle}>
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
          className={classes.test}
          subheader={kpi.value + (kpi.unit ?? '')}
          titleTypographyProps={{
            variant: 'body2',
            className: classes.kpiLabel,
            title: kpi.category,
          }}
          subheaderTypographyProps={{
            variant: 'h5',
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
