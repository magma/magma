/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {magmad_gateway_configs} from '@fbcnms/magma-api';

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
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
  configs: magmad_gateway_configs,
  configChangeHandler: (string, any) => void,
};

class MagmaDeviceFields extends React.Component<Props> {
  render() {
    return (
      <>
        <FormControl className={this.props.classes.input}>
          <InputLabel htmlFor="autoupgradeEnabled">
            Autoupgrade Enabled
          </InputLabel>
          <Select
            inputProps={{id: 'autoupgradeEnabled'}}
            value={this.props.configs.autoupgrade_enabled ? 1 : 0}
            onChange={this.autoupgradeEnabledChanged}>
            <MenuItem value={1}>Enabled</MenuItem>
            <MenuItem value={0}>Disabled</MenuItem>
          </Select>
        </FormControl>
        <TextField
          label="Autoupgrade Poll Interval (seconds)"
          className={this.props.classes.input}
          value={this.props.configs.autoupgrade_poll_interval}
          onChange={this.autoupgradePollIntervalChanged}
          placeholder="E.g. 300"
        />
        <TextField
          label="Checkin Interval (seconds)"
          className={this.props.classes.input}
          value={this.props.configs.checkin_interval}
          onChange={this.checkinIntervalChanged}
          placeholder="E.g. 60"
        />
        <TextField
          label="Checkin Timeout (seconds)"
          className={this.props.classes.input}
          value={this.props.configs.checkin_timeout}
          onChange={this.checkinTimeoutChanged}
          placeholder="E.g. 5"
        />
      </>
    );
  }

  autoupgradeEnabledChanged = ({target}) =>
    this.props.configChangeHandler('autoupgrade_enabled', !!target.value);

  autoupgradePollIntervalChanged = ({target}) =>
    this.props.configChangeHandler('autoupgrade_poll_interval', target.value);

  checkinIntervalChanged = ({target}) =>
    this.props.configChangeHandler('checkin_interval', target.value);

  checkinTimeoutChanged = ({target}) =>
    this.props.configChangeHandler('checkin_timeout', target.value);
}

export default withStyles(styles)(MagmaDeviceFields);
