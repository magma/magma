/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import LoadingFiller from './LoadingFiller';
import React from 'react';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  backdrop: {
    alignItems: 'center',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    bottom: 0,
    display: 'flex',
    justifyContent: 'center',
    left: 0,
    position: 'fixed',
    right: 0,
    top: 0,
    zIndex: '13000',
  },
}));

export default function LoadingFillerBackdrop() {
  const classes = useStyles();
  return (
    <div className={classes.backdrop}>
      <LoadingFiller />
    </div>
  );
}
