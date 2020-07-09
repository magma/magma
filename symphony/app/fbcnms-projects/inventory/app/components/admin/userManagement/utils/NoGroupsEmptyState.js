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
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  noGroups: {
    width: '124px',
    position: 'absolute',
    top: '33%',
    right: '50%',
    transform: 'translate(50%, -50%)',
    textAlign: 'center',
  },
}));

export default function NoGroupsEmptyState() {
  const classes = useStyles();

  return (
    <img
      className={classes.noGroups}
      src="/inventory/static/images/noGroups.png"
    />
  );
}
