/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  noCommentsEmptyState: {
    textAlign: 'center',
    paddingTop: '16px',
  },
});

const CommentsLogEmptyState = () => {
  const classes = useStyles();
  return <div className={classes.noCommentsEmptyState}>No comments yet</div>;
};

export default CommentsLogEmptyState;
