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
  input: {
    width: '500px',
    paddingRight: '10px',
  },
  icon: {
    width: '30px',
    height: '30px',
    verticalAlign: 'bottom',
  },
};

type Props = WithStyles<typeof styles> & {
  itemList: Array<string>,
  onChange: (Array<string>) => void,
};

type State = {};

class ListFields extends React.Component<Props, State> {
  render() {
    return this.props.itemList.map((item, index) => (
      <div className={this.props.classes.container} key={index}>
        <TextField
          label="Item"
          margin="none"
          value={item}
          onChange={({target}) => this.onChange(index, target.value)}
          className={this.props.classes.input}
        />
        {this.props.itemList.length !== 1 && (
          <IconButton
            onClick={() => this.removeField(index)}
            className={this.props.classes.icon}>
            <RemoveCircleOutline />
          </IconButton>
        )}
        {index === this.props.itemList.length - 1 && (
          <IconButton
            onClick={this.addField}
            className={this.props.classes.icon}>
            <AddCircleOutline />
          </IconButton>
        )}
      </div>
    ));
  }

  onChange = (index, value) => {
    const itemList = this.props.itemList.slice(0);
    itemList[index] = value;
    this.props.onChange(itemList);
  };

  removeField = index => {
    const itemList = this.props.itemList.slice(0);
    itemList.splice(index, 1);
    this.props.onChange(itemList);
  };

  addField = () => {
    const itemList = this.props.itemList.slice(0);
    itemList.push('');
    this.props.onChange(itemList);
  };
}

export default withStyles(styles)(ListFields);
