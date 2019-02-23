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
import Typography from '@material-ui/core/Typography';

import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    display: 'flex',
    marginBottom: '5px',
  },
  heading: {
    flexBasis: '33.33%',
    marginRight: '15px',
    textAlign: 'left',
  },
  secondaryHeading: {
    color: theme.palette.text.secondary,
    flexBasis: '66.66%',
  },
});

type Props = WithStyles & {
  label: string,
  children?: any,
};

class FormField extends React.Component<Props> {
  render() {
    const {classes, label} = this.props;
    return (
      <div className={classes.root}>
        <Typography className={classes.heading}>{label}</Typography>
        <Typography className={classes.secondaryHeading} component="div">
          {this.props.children}
        </Typography>
      </div>
    );
  }
}

export default withStyles(styles)(FormField);
