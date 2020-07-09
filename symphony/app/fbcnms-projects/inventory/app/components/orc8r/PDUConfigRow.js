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
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  container: {
    margin: '5px 0',
    width: '100%',
  },
  inputRow: {
    margin: 0,
  },
}));

export type DeviceConfig = {
  id: string,
  name: string,
  ipAddress: string,
};

type Props = {
  config: DeviceConfig,
  onChange: DeviceConfig => void,
};

export default function (props: Props) {
  const classes = useStyles();
  const onChange = (field: 'name' | 'ipAddress') => event => {
    const newConfig = {
      ...props.config,
      // $FlowFixMe Set state for each field
      [field]: event.target.value,
    };
    props.onChange(newConfig);
  };

  return (
    <div className={classes.container}>
      <TextField
        required
        className={classes.inputRow}
        label="Name"
        margin="normal"
        onChange={onChange('name')}
        value={props.config.name}
      />
      <TextField
        required
        className={classes.inputRow}
        label="IP Address"
        margin="normal"
        onChange={onChange('ipAddress')}
        value={props.config.ipAddress}
      />
    </div>
  );
}
