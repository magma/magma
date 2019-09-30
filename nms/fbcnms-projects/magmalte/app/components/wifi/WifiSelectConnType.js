/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LineMapLayer} from './WifiMapLayers';
import type {WithStyles} from '@material-ui/core';

import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  formControl: {
    margin: theme.spacing(),
    minWidth: 120,
    width: 'calc(100% - 15px)',
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
  },
});

type Props = WithStyles<typeof styles> & {
  onChange: (connType: LineMapLayer | '') => void,
  selectedConnType: LineMapLayer | '',
};

class WifiSelectConnType extends React.Component<Props> {
  handleChange = event => {
    // convert type
    const t: LineMapLayer | '' = (event.target.value: any);
    this.props.onChange(t);
  };

  render() {
    return (
      <FormControl className={this.props.classes.formControl}>
        <InputLabel htmlFor="conntype-helper">Connection Filter</InputLabel>
        <Select
          value={this.props.selectedConnType}
          onChange={this.handleChange}
          input={<Input name="connType" id="conntype-helper" />}>
          <MenuItem value={''}>All</MenuItem>
          <MenuItem value={'defaultRoute'}>Default Routes</MenuItem>
          <MenuItem value={'l3'}>L3 only</MenuItem>
          <MenuItem value={'l2'}>L2 only</MenuItem>
          <MenuItem value={'none'}>Visible (low signal)</MenuItem>
        </Select>
        <FormHelperText>Filter by Connection Type</FormHelperText>
      </FormControl>
    );
  }
}

export default withStyles(styles)(WifiSelectConnType);
