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

import Grid from '@material-ui/core/Grid';
import React from 'react';
import Typography from '@material-ui/core/Typography';

import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  labelName: {
    color: theme.palette.grey.A700,
    fontWeight: 500,
  },
});

type Props = WithStyles & {
  label: string,
  value?: ?string | ?number,
};

class FormField extends React.Component<Props> {
  render() {
    const {classes, label, value} = this.props;
    return (
      <div>
        <Grid container spacing={16}>
          <Grid item xs={4}>
            <Typography className={classes.labelName} variant="body2">
              {label}:
            </Typography>
          </Grid>
          <Grid item xs={8}>
            <Typography variant="body2" color="secondary">
              {value}
            </Typography>
          </Grid>
        </Grid>
      </div>
    );
  }
}

export default withStyles(styles)(FormField);
