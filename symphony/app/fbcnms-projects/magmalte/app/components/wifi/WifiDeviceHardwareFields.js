/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {gateway_device} from '@fbcnms/magma-api';

import React from 'react';
import TextField from '@material-ui/core/TextField';

import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
};

type Props = WithStyles<typeof styles> & {
  record: gateway_device,
};

class WifiDeviceHardwareFields extends React.Component<Props> {
  render() {
    return (
      <>
        <TextField
          label="HW ID"
          className={this.props.classes.input}
          value={this.props.record.hardware_id}
          disabled={true}
        />
      </>
    );
  }
}

export default withStyles(styles)(WifiDeviceHardwareFields);
