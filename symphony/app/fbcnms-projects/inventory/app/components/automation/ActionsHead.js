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

import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  header: {
    backgroundColor: theme.palette.common.white,
    borderBottom: '1px solid ' + symphony.palette.separator,
    display: 'flex',
    padding: '0px 16px',
    width: '100%',
    height: '60px',
    alignItems: 'center',
  },
}));

export default function ActionsHead({children}: {children: React.Node}) {
  const classes = useStyles();
  return <div className={classes.header}>{children}</div>;
}
