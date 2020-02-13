/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import {gray7} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

type Props = {
  children: React.Node,
};

const useStyles = makeStyles(() => ({
  content: {
    flexGrow: 1,
    height: '100vh',
    overflow: 'auto',
    overflowX: 'hidden',
    backgroundColor: gray7,
  },
}));

export default function AppContent(props: Props) {
  const classes = useStyles();
  return <main className={classes.content}>{props.children}</main>;
}
