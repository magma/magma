/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {gateway_device} from '@fbcnms/magma-api';

import React from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  record: gateway_device,
};

export default function WifiDeviceHardwareFields(props: Props) {
  const classes = useStyles();
  return (
    <>
      <TextField
        label="HW ID"
        className={classes.input}
        value={props.record.hardware_id}
        disabled={true}
      />
    </>
  );
}
