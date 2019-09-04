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

import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import RemoveCircleOutline from '@material-ui/icons/RemoveCircleOutline';
import TextField from '@material-ui/core/TextField';

import {withStyles} from '@material-ui/core/styles';

const styles = {
  container: {
    display: 'block',
    margin: '5px 0',
    whiteSpace: 'nowrap',
    width: '100%',
  },
  inputKey: {
    width: '245px',
    paddingRight: '10px',
  },
  inputValue: {
    width: '240px',
  },
  icon: {
    width: '30px',
    height: '30px',
    verticalAlign: 'bottom',
  },
};

type Props = WithStyles<typeof styles> & {
  keyValuePairs: Array<[string, string]>,
  onChange: (Array<[string, string]>) => void,
};

type State = {};

class KeyValueFields extends React.Component<Props, State> {
  render() {
    return this.props.keyValuePairs.map((pair, index) => (
      <div className={this.props.classes.container} key={index}>
        <TextField
          label="Key"
          margin="none"
          value={pair[0]}
          onChange={({target}) => this.onChange(index, 0, target.value)}
          className={this.props.classes.inputKey}
        />
        <TextField
          label="Value"
          margin="none"
          value={pair[1]}
          onChange={({target}) => this.onChange(index, 1, target.value)}
          className={this.props.classes.inputValue}
        />
        {this.props.keyValuePairs.length !== 1 && (
          <IconButton
            onClick={() => this.removeField(index)}
            className={this.props.classes.icon}>
            <RemoveCircleOutline />
          </IconButton>
        )}
        {index === this.props.keyValuePairs.length - 1 && (
          <IconButton
            onClick={this.addField}
            className={this.props.classes.icon}>
            <AddCircleOutline />
          </IconButton>
        )}
      </div>
    ));
  }

  onChange = (index, subIndex, value) => {
    const keyValuePairs = this.props.keyValuePairs.slice(0);
    keyValuePairs[index] = [keyValuePairs[index][0], keyValuePairs[index][1]];
    keyValuePairs[index][subIndex] = value;
    this.props.onChange(keyValuePairs);
  };

  removeField = index => {
    const keyValuePairs = this.props.keyValuePairs.slice(0);
    keyValuePairs.splice(index, 1);
    this.props.onChange(keyValuePairs);
  };

  addField = () => {
    const keyValuePairs = this.props.keyValuePairs.slice(0);
    keyValuePairs.push(['', '']);
    this.props.onChange(keyValuePairs);
  };
}

export default withStyles(styles)(KeyValueFields);
