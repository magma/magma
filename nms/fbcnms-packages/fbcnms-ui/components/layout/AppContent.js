/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import {makeStyles} from '@material-ui/styles';

type Props = {
  children: any,
};

const useStyles = makeStyles({
  content: {
    flexGrow: 1,
    height: '100vh',
    overflow: 'auto',
  },
});

export default function AppContent(props: Props) {
  const classes = useStyles();
  return <main className={classes.content}>{props.children}</main>;
}
