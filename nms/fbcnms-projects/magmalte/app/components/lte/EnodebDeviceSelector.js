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

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import PropTypes from 'prop-types';
import React from 'react';
import Select from '@material-ui/core/Select';

import {EnodebDeviceClass} from './EnodebUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  formControl: {
    margin: theme.spacing(),
    minWidth: 120,
  },
});

type Props = WithStyles & {
  value: $Values<typeof EnodebDeviceClass>,
  onChange: (SyntheticEvent<>) => void,
  className: string,
};

type State = {
  open: boolean,
};

class EnodebDeviceSelector extends React.Component<Props, State> {
  state = {
    open: false,
  };

  handleChange = (event: SyntheticEvent<>) => {
    this.props.onChange(event);
    this.setState({
      // $FlowFixMe: event target will have name and value
      [event.target.name]: event.target.value,
    });
  };

  handleClose = () => {
    this.setState({open: false});
  };

  handleOpen = () => {
    this.setState({open: true});
  };

  render() {
    const deviceLabelArr = [];
    for (const property in EnodebDeviceClass) {
      if (EnodebDeviceClass.hasOwnProperty(property)) {
        deviceLabelArr.push(EnodebDeviceClass[property]);
      }
    }

    const menuItems = deviceLabelArr.map(deviceClass => {
      return (
        <MenuItem key={deviceClass} value={deviceClass}>
          {deviceClass}
        </MenuItem>
      );
    });

    return (
      <form autoComplete="off">
        <FormControl className={this.props.classes.input}>
          <InputLabel htmlFor="demo-controlled-open-select">
            eNodeB Device Class
          </InputLabel>
          <Select
            open={this.state.open}
            onClose={this.handleClose}
            onOpen={this.handleOpen}
            value={this.props.value}
            onChange={this.handleChange}
            inputProps={{
              name: 'Device Class',
              id: 'demo-controlled-open-select',
            }}>
            {menuItems}
          </Select>
        </FormControl>
      </form>
    );
  }
}

EnodebDeviceSelector.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(EnodebDeviceSelector);
