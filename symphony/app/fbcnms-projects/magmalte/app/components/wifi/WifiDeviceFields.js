/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {gateway_status, gateway_wifi_configs} from '@fbcnms/magma-api';

import Checkbox from '@material-ui/core/Checkbox';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {withStyles} from '@material-ui/core/styles';

import {getAdditionalProp, setAdditionalProp} from './WifiUtils';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
};

type Props = WithStyles<typeof styles> & {
  macAddress: string,
  status: ?gateway_status,
  configs: gateway_wifi_configs,
  additionalProps: Array<[string, string]>,
  handleMACAddressChange?: string => void,
  configChangeHandler: (string, string | number) => void,
  additionalPropsChangeHandler: (Array<[string, string]>) => void,
};

class WifiDeviceFields extends React.Component<Props> {
  render() {
    const reboot_if_bootid = getAdditionalProp(
      this.props.additionalProps,
      'reboot_if_bootid',
    );
    return (
      <>
        <FormGroup row>
          {this.props.handleMACAddressChange && (
            <TextField
              required
              className={this.props.classes.input}
              label="MAC Address"
              margin="normal"
              onChange={this.handleMACAddressChange}
              value={this.props.macAddress}
            />
          )}
          <TextField
            required
            className={this.props.classes.input}
            label="Info"
            margin="normal"
            onChange={this.handleInfoChange}
            value={this.props.configs.info}
          />
          <TextField
            className={this.props.classes.input}
            label="Latitude"
            margin="normal"
            onChange={this.handleLatitudeChange}
            value={this.props.configs.latitude}
          />
          <TextField
            className={this.props.classes.input}
            label="Longitude"
            margin="normal"
            onChange={this.handleLongitudeChange}
            value={this.props.configs.longitude}
          />
          <TextField
            className={this.props.classes.input}
            label="Client Channel"
            margin="normal"
            onChange={this.handleClientChannelChange}
            value={this.props.configs.client_channel}
          />
          <FormControlLabel
            control={
              <Checkbox
                checked={this.props.configs.is_production}
                onChange={this.handleIsProductionChange}
                color="primary"
              />
            }
            label="Is Production"
          />

          <FormControlLabel
            control={
              <Checkbox
                disabled={this.props.status === null}
                checked={
                  this.props.status !== null &&
                  reboot_if_bootid !== null &&
                  reboot_if_bootid === this.props.status?.meta?.boot_id
                }
                onChange={this.handleRequestReboot}
                color="primary"
              />
            }
            label="Reboot requested"
          />
        </FormGroup>
        <KeyValueFields
          keyValuePairs={this.props.additionalProps || [['', '']]}
          onChange={this.props.additionalPropsChangeHandler}
        />
      </>
    );
  }

  handleMACAddressChange = ({target}) =>
    nullthrows(this.props.handleMACAddressChange)(target.value);
  handleInfoChange = ({target}) =>
    this.props.configChangeHandler('info', target.value);
  handleLatitudeChange = ({target}) =>
    this.props.configChangeHandler('latitude', target.value);
  handleLongitudeChange = ({target}) =>
    this.props.configChangeHandler('longitude', target.value);
  handleClientChannelChange = ({target}) =>
    this.props.configChangeHandler('client_channel', target.value);
  handleIsProductionChange = ({target}) =>
    this.props.configChangeHandler('is_production', target.checked);

  handleRequestReboot = ({target}) => {
    const keyValuePairs = this.props.additionalProps.slice(0);
    if (target.checked && this.props.status && this.props.status.meta) {
      // add the reboot directive
      setAdditionalProp(
        keyValuePairs,
        'reboot_if_bootid',
        this.props.status.meta.boot_id,
      );
    } else {
      // remove the reboot directive
      setAdditionalProp(keyValuePairs, 'reboot_if_bootid', undefined);
      // if there are no key/values, then add a dummy line for UI purposes
      if (keyValuePairs.length === 0) {
        keyValuePairs.push(['', '']);
      }
    }
    this.props.additionalPropsChangeHandler(keyValuePairs);
  };
}

export default withStyles(styles)(WifiDeviceFields);
