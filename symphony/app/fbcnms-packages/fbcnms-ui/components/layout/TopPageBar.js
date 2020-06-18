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
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    backgroundColor: theme.palette.common.white,
    height: '60px',
    borderBottom: '1px solid rgba(0, 0, 0, 0.1)',
    display: 'flex',
    padding: '0px 16px',
    width: '100%',
    alignItems: 'center',
  },
}));

type Props = {
  children: ?any,
};

export default function TopPageBar(props: Props) {
  const classes = useStyles();
  return <div className={classes.root}>{props.children}</div>;
}
