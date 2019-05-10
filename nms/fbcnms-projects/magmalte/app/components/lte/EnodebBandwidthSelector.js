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
import React from 'react';
import Select from '@material-ui/core/Select';

import {EnodebBandwidthOption} from './EnodebUtils';

type Props = WithStyles & {
  value: $Values<typeof EnodebBandwidthOption>,
  onChange: (SyntheticEvent<>) => void,
  className: string,
};

type State = {
  open: boolean,
};

class EnodebBandwidthSelector extends React.Component<Props, State> {
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
    const {className} = this.props;
    const deviceBandwidthArr = [];
    for (const property in EnodebBandwidthOption) {
      if (EnodebBandwidthOption.hasOwnProperty(property)) {
        deviceBandwidthArr.push(EnodebBandwidthOption[property]);
      }
    }

    const menuItems = deviceBandwidthArr.map(bandwidthMhz => {
      return (
        <MenuItem key={bandwidthMhz} value={bandwidthMhz}>
          {bandwidthMhz}
        </MenuItem>
      );
    });

    return (
      <form autoComplete="off">
        <FormControl className={className}>
          <InputLabel htmlFor="demo-controlled-open-select">
            eNodeB DL/UL Bandwidth (MHz)
          </InputLabel>
          <Select
            open={this.state.open}
            onClose={this.handleClose}
            onOpen={this.handleOpen}
            value={this.props.value}
            onChange={this.handleChange}
            inputProps={{
              name: 'Bandwidth (MHz)',
              id: 'demo-controlled-open-select',
            }}>
            {menuItems}
          </Select>
        </FormControl>
      </form>
    );
  }
}

export default EnodebBandwidthSelector;
