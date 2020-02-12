/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

type Props = {
  titleLabel: string,
  value: number | string,
  valueOptionsByKey: {+[string]: number | string},
  onChange: (SyntheticInputEvent<>) => void,
  className: string,
};

type State = {
  open: boolean,
};

class EnodebPropertySelector extends React.Component<Props, State> {
  state = {
    open: false,
  };

  handleChange = (event: SyntheticInputEvent<>) => {
    this.props.onChange(event);
  };

  handleClose = () => {
    this.setState({open: false});
  };

  handleOpen = () => {
    this.setState({open: true});
  };

  render() {
    const {className, valueOptionsByKey} = this.props;
    const valueOptionsArr = [];
    for (const property in valueOptionsByKey) {
      if (valueOptionsByKey.hasOwnProperty(property)) {
        valueOptionsArr.push(valueOptionsByKey[property]);
      }
    }

    const menuItems = valueOptionsArr.map(valueOption => {
      return (
        <MenuItem key={valueOption} value={valueOption}>
          {valueOption}
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
              name: this.props.titleLabel,
              id: 'demo-controlled-open-select',
            }}>
            {menuItems}
          </Select>
        </FormControl>
      </form>
    );
  }
}

export default EnodebPropertySelector;
