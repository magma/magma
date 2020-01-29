/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import classNames from 'classnames';
import update from 'immutability-helper';
import {isJSON} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  className: string,
  inputClassName?: ?string,
  label: ?string,
  disabled?: boolean,
  property: Property | PropertyType,
  onChange: (Property | PropertyType) => void,
  margin: 'none' | 'dense' | 'normal',
  autoFocus?: boolean,
} & WithStyles<typeof styles>;

const styles = {
  input: {
    width: '300px',
    display: 'flex',
  },
  container: {
    display: 'flex',
    width: '250px',
  },
};

class EnumPropertySelectValueInput extends React.Component<Props> {
  render() {
    const {
      classes,
      onChange,
      margin,
      property,
      className,
      inputClassName,
      autoFocus,
      disabled = false,
    } = this.props;
    const propertyType = !!property.propertyType
      ? property.propertyType
      : property;
    const jsonStr = propertyType.stringValue || '';
    const options = isJSON(jsonStr) ? JSON.parse(jsonStr) : [];
    const optionsArr = Array.isArray(options) ? options : [];
    return (
      <TextField
        select
        id="property-value"
        variant="outlined"
        autoFocus={autoFocus}
        disabled={disabled}
        className={classNames(classes.input, className)}
        inputProps={{className: inputClassName}}
        label={this.props.label}
        margin={margin}
        value={property && property.stringValue ? property.stringValue : ''}
        onChange={event => {
          if (property != null) {
            onChange(
              update(property, {
                stringValue: {
                  $set: event.target.value,
                },
              }),
            );
          } else {
            onChange(
              update(propertyType, {
                stringValue: {
                  $set: event.target.value,
                },
              }),
            );
          }
        }}>
        {optionsArr.map(stringVal => (
          <MenuItem key={stringVal} value={stringVal}>
            {stringVal}
          </MenuItem>
        ))}
      </TextField>
    );
  }
}

export default withStyles(styles)(EnumPropertySelectValueInput);
