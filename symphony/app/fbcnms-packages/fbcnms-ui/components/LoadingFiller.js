/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CircularProgress from '@material-ui/core/CircularProgress';
import React from 'react';

import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  loadingContainer: {
    minHeight: 500,
    paddingTop: 200,
    textAlign: 'center',
  },
});

const LoadingFiller = ({classes}) => (
  <div className={classes.loadingContainer}>
    <CircularProgress size={50} />
  </div>
);

export default withStyles(styles)(LoadingFiller);
