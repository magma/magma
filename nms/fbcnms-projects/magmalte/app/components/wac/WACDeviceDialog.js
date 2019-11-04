/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WifiGateway} from '../wifi/WifiUtils';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import nullthrows from '@fbcnms/util/nullthrows';
import {
  DEFAULT_HW_ID_PREFIX,
  DEFAULT_WIFI_GATEWAY_CONFIGS,
  buildWifiGatewayFromPayload,
  meshesURL,
} from '../wifi/WifiUtils';
import {createDevice} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
};

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    onCancel: () => void,
    onSave: WifiGateway => void,
  };

type State = {
  serialNumber: string,
  error: string,
};

class WifiDeviceDialog extends React.Component<Props, State> {
  state = {
    serialNumber: '',
    error: '',
  };

  render() {
    return (
      <Dialog open={true} onClose={this.props.onCancel}>
        <DialogTitle>Add Device</DialogTitle>
        <DialogContent>
          {this.state.error ? (
            <FormLabel error>{this.state.error}</FormLabel>
          ) : null}
          <FormGroup row>
            <TextField
              required
              className={this.props.classes.input}
              label="Serial Number"
              margin="normal"
              onChange={this.handleSerialNumberChange}
              value={this.state.serialNumber}
            />
          </FormGroup>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onCancel} skin="regular">
            Cancel
          </Button>
          <Button onClick={this.onSave}>Save</Button>
        </DialogActions>
      </Dialog>
    );
  }

  handleSerialNumberChange = ({target}) =>
    this.setState({serialNumber: target.value});

  onSave = async () => {
    const {match} = this.props;

    if (!this.state.serialNumber) {
      this.setState({error: 'Serial Number cannot be empty'});
      return;
    }

    const meshID = nullthrows(match.params.meshID);
    const {serialNumber} = this.state;
    const data = {
      hardware_id: DEFAULT_HW_ID_PREFIX + serialNumber,
      key: {key_type: 'ECHO'},
    };

    try {
      const groupURL = meshesURL(this.props.match) + '/' + meshID;
      const [deviceResult, groupResult] = await Promise.all([
        createDevice(
          serialNumber,
          data,
          'wifi', // type
          DEFAULT_WIFI_GATEWAY_CONFIGS,
          {mesh_id: meshID, info: serialNumber}, // wifi configs
          match,
        ),
        axios.get(groupURL),
      ]);

      // update devices list for group
      const devices = groupResult.data || [];
      devices.push(serialNumber);
      await axios.put(groupURL, devices);

      this.props.onSave(buildWifiGatewayFromPayload(deviceResult));
    } catch (e) {
      this.setState({error: e.response?.data?.message || e.message || e});
    }
  };
}

export default withStyles(styles)(withRouter(WifiDeviceDialog));
