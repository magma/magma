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
  onChange: (meshId: string) => void,
  meshes: Array<string>,
  selectedMeshID: string,
  disallowEmpty?: boolean,
  helperText?: string,
};

class WifiSelectMesh extends React.Component<Props> {
  handleChange = event => this.props.onChange(event.target.value);

  render() {
    if (!this.props.meshes) {
      return null;
    }

    const {classes} = this.props;

    this.props.meshes.sort((a, b) =>
      a.toLowerCase() > b.toLowerCase() ? 1 : -1,
    );

    const meshItems = this.props.meshes.map(meshId => (
      <MenuItem value={meshId} key={meshId}>
        {meshId}
      </MenuItem>
    ));

    return (
      <FormControl className={classes.formControl}>
        <InputLabel htmlFor="meshid-helper">Mesh ID</InputLabel>
        <Select
          value={this.props.selectedMeshID}
          onChange={this.handleChange}
          input={<Input name="meshId" id="meshid-helper" />}>
          {!this.props.disallowEmpty && (
            <MenuItem value="">
              <em>All</em>
            </MenuItem>
          )}
          {meshItems}
        </Select>
        {this.props.helperText && (
          <FormHelperText>{this.props.helperText}</FormHelperText>
        )}
      </FormControl>
    );
  }
}

export default withStyles(styles)(WifiSelectMesh);
