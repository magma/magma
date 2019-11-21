/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Theme, WithStyles} from '@material-ui/core';

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import InputAffix from '@fbcnms/ui/components/design-system/Input/InputAffix';
import React from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import classNames from 'classnames';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  value: any,
  className?: string,
  required: boolean,
  disabled: boolean,
  onLatitudeChange: (event: SyntheticInputEvent<any>) => void,
  onLongitudeChange: (event: SyntheticInputEvent<any>) => void,
  margin: 'none' | 'dense' | 'normal',
  fullWidth?: boolean,
  autoFocus?: boolean,
} & WithStyles<typeof styles>;

const styles = (theme: Theme) => ({
  container: {
    display: 'flex',
    width: '280px',
  },
  fullWidth: {
    width: '100%',
  },
  input: {
    marginLeft: '0px',
    marginRight: theme.spacing(),
    width: '100%',
  },
  lngField: {
    marginLeft: '16px',
  },
  field: {
    flexGrow: 1,
  },
});

class GPSPropertyValueInput extends React.Component<Props> {
  constructor(props) {
    super(props);
  }

  static defaultProps = {
    required: false,
    disabled: false,
  };

  render() {
    const {
      className,
      classes,
      disabled,
      margin,
      required,
      autoFocus,
      fullWidth,
    } = this.props;
    const {latitude, longitude} = this.props.value || {
      latitude: '',
      longitude: '',
    };
    const latError = Number(latitude) < -90 || Number(latitude) > 90;
    const longError = Number(longitude) < -180 || Number(longitude) > 180;
    return (
      <div
        className={classNames(classes.container, className, {
          [classes.fullWidth]: fullWidth,
        })}>
        <FormField
          errorText={'Latitude should be between -90 and 90'}
          hasError={latError}
          className={classes.field}>
          <TextInput
            autoFocus={autoFocus}
            prefix={<InputAffix>Lat.</InputAffix>}
            required={required}
            disabled={disabled}
            id="longitude-value"
            variant="outlined"
            className={classes.input}
            margin={margin}
            value={latitude}
            type="number"
            onChange={this.props.onLatitudeChange}
          />
        </FormField>
        <FormField
          className={classNames(classes.lngField, classes.field)}
          errorText={'Longitude should be between -180 and 180'}
          hasError={longError}>
          <TextInput
            required={required}
            prefix={<InputAffix>Long.</InputAffix>}
            disabled={disabled}
            id="latitude-value"
            variant="outlined"
            className={classes.input}
            margin={margin}
            type="number"
            value={longitude}
            error={longError}
            onChange={this.props.onLongitudeChange}
          />
        </FormField>
      </div>
    );
  }
}

export default withStyles(styles)(GPSPropertyValueInput);
