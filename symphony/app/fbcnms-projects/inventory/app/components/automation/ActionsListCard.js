/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Paper from '@material-ui/core/Paper';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  paper: {
    padding: theme.spacing(3, 2),
  },
}));

export default function ActionsListCard() {
  const classes = useStyles();
  return (
    <Paper className={classes.paper}>
      <div>My Custom Action</div>
      <div>Owner: jbraeg@fb.com</div>
      <div>Status: Enabled</div>
      <div>Last Updated: 3/7/2019</div>
    </Paper>
  );
}
