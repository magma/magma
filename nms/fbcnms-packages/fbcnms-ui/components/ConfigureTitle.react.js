/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';

import React from 'react';
import {Typography} from '@material-ui/core';
import {gray13} from '../theme/colors';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  title: string,
  subtitle: ?string,
  className?: string,
} & WithStyles<typeof styles>;

const styles = theme => ({
  title: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
    color: theme.palette.blueGrayDark,
  },
  subtitle: {
    fontSize: '14px',
    lineHeight: '24px',
    fontWeight: 500,
    color: gray13,
  },
});

const ConfigureTitle = (props: Props) => {
  const {title, subtitle, classes, className} = props;
  return (
    <div className={className}>
      <Typography className={classes.title}>{title}</Typography>
      <Typography className={classes.subtitle}>{subtitle}</Typography>
    </div>
  );
};

export default withStyles(styles)(ConfigureTitle);
