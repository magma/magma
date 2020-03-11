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

type Props = {
  isGrey: boolean,
  isActive: boolean,
};

const useStyles = makeStyles(() => ({
  status: {
    width: '10px',
    height: '10px',
    borderRadius: '50%',
    display: 'inline-block',
    textAlign: 'center',
    color: 'white',
    fontSize: '10px',
    fontWeight: 'bold',
    marginRight: '5px',
  },
}));

export default function DeviceStatusCircle(props: Props) {
  const classes = useStyles();
  if (props.isGrey) {
    return (
      <span className={classes.status} style={{backgroundColor: '#bec3c8'}} />
    );
  } else if (props.isActive) {
    return (
      <span className={classes.status} style={{backgroundColor: '#05a503'}} />
    );
  } else {
    return (
      <span className={classes.status} style={{border: '2px solid #fa3a3f'}} />
    );
  }
}
