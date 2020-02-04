/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {ButtonProps} from '@fbcnms/ui/components/design-system/Button';
import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import classNames from 'classnames';
import update from 'immutability-helper';
import {isJSON} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

type Props = {|
  className: string,
  property: Property | PropertyType,
  onChange: (Property | PropertyType) => void,
  ...ButtonProps,
|} & WithStyles<typeof styles>;

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
      property,
      className,
      ...restButtonProps
    } = this.props;
    const propertyType = !!property.propertyType
      ? property.propertyType
      : property;
    const jsonStr = propertyType.stringValue || '';
    const options = isJSON(jsonStr) ? JSON.parse(jsonStr) : [];
    const optionsArr = Array.isArray(options) ? options : [];
    return (
      <Select
        className={classNames(classes.input, className)}
        options={optionsArr.map(stringVal => ({
          value: stringVal,
          label: stringVal,
        }))}
        selectedValue={
          property && property.stringValue ? property.stringValue : ''
        }
        {...restButtonProps}
        onChange={value => {
          if (property != null) {
            onChange(
              update(property, {
                stringValue: {
                  $set: value,
                },
              }),
            );
          } else {
            onChange(
              update(propertyType, {
                stringValue: {
                  $set: value,
                },
              }),
            );
          }
        }}
      />
    );
  }
}

export default withStyles(styles)(EnumPropertySelectValueInput);
