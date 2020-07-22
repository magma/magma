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
import type {ComponentType} from 'react';

import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '../../theme/design-system/Text';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: theme.spacing(1),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

type CardTitleRowProps = {
  icon?: ComponentType<SvgIconExports>,
  label: string,
};

export const CardTitleRow = (props: CardTitleRowProps) => {
  const classes = useStyles();
  const Icon = props.icon;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      {Icon ? <Icon className={classes.cardTitleIcon} /> : null}
      <Text variant="body1">{props.label}</Text>
    </Grid>
  );
};

type CardTitleFilterRowProps = {
  icon?: ComponentType<SvgIconExports>,
  label: string,
  filter: () => React$Node,
};

export const CardTitleFilterRow = (props: CardTitleFilterRowProps) => {
  const classes = useStyles();
  const Filters = props.filter;
  const Icon = props.icon;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      <Grid container xs>
        {Icon ? <Icon className={classes.cardTitleIcon} /> : null}
        <Text variant="body1">{props.label}</Text>
      </Grid>
      <Grid item>
        <Filters />
      </Grid>
    </Grid>
  );
};
