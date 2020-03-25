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
import Typography from '@material-ui/core/Typography';
import WifiTethering from '@material-ui/icons/WifiTethering';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  noAccess: {
    color: theme.palette.gray13,
    top: '50%',
    width: '520px',
    position: 'relative',
    margin: 'auto',
    textAlign: 'center',
  },
  icon: {
    width: '60px',
    height: '60px',
  },
}));

export default function({children}: {children: React.Node}) {
  const classes = useStyles();
  return (
    <Typography variant="h6" className={classes.noAccess}>
      <div>
        <WifiTethering className={classes.icon} />
      </div>
      {children}
    </Typography>
  );
}
