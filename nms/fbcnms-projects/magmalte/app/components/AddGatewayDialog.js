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
import type {WithStyles} from '@material-ui/core';
import type {lte_gateway} from '../common/__generated__/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import MagmaV1API from '../common/MagmaV1API';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
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
    open: boolean,
    onClose: () => void,
    onSave: lte_gateway => void,
  };

type State = {
  error: string,
  name: string,
  description: string,
  hardware_id: string,
  gatewayID: string,
  challengeKey: string,
};

class AddGatewayDialog extends React.Component<Props, State> {
  state = {
    error: '',
    name: '',
    description: '',
    hardware_id: '',
    gatewayID: '',
    challengeKey: '',
  };

  render() {
    const {classes} = this.props;
    const error = this.state.error ? (
      <FormLabel error>{this.state.error}</FormLabel>
    ) : null;

    return (
      <Dialog open={this.props.open} onClose={this.props.onClose}>
        <DialogTitle>Add Gateway</DialogTitle>
        <DialogContent>
          {error}
          <TextField
            label="Gateway Name"
            className={classes.input}
            value={this.state.name}
            onChange={this.onNameChange}
            placeholder="Gateway 1"
          />
          <TextField
            label="Gateway Description"
            className={classes.input}
            value={this.state.description}
            onChange={this.onDescriptionChange}
            placeholder="Sample Gateway description"
          />
          <TextField
            label="Hardware UUID"
            className={classes.input}
            value={this.state.hardware_id}
            onChange={this.onHwidChange}
            placeholder="Eg. 4dfe212f-df33-4cd2-910c-41892a042fee"
          />
          <TextField
            label="Gateway ID"
            className={classes.input}
            value={this.state.gatewayID}
            onChange={this.onGatewayIDChange}
            placeholder="<country>_<org>_<location>_<sitenumber>"
          />
          <TextField
            label="Challenge Key"
            className={classes.input}
            value={this.state.challengeKey}
            onChange={this.onChallengeKeyChange}
            placeholder="A base64 bytestring of the key in DER format"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    );
  }

  onNameChange = ({target}) => this.setState({name: target.value});
  onDescriptionChange = ({target}) =>
    this.setState({description: target.value});
  onHwidChange = ({target}) => this.setState({hardware_id: target.value});
  onGatewayIDChange = ({target}) => this.setState({gatewayID: target.value});
  onChallengeKeyChange = ({target}) =>
    this.setState({challengeKey: target.value});

  onSave = async () => {
    const {
      name,
      description,
      hardware_id,
      gatewayID,
      challengeKey,
    } = this.state;
    if (!name || !description || !hardware_id || !gatewayID || !challengeKey) {
      this.setState({error: 'Please complete all fields'});
      return;
    }

    try {
      const networkId = nullthrows(this.props.match.params.networkId);
      await MagmaV1API.postLteByNetworkIdGateways({
        networkId,
        gateway: {
          id: gatewayID,
          name,
          description,
          cellular: {
            epc: {nat_enabled: false, ip_block: '192.168.0.1/32'},
            ran: {pci: 260, transmit_enabled: false},
            non_eps_service: undefined,
          },
          magmad: {
            autoupgrade_enabled: true,
            autoupgrade_poll_interval: 300,
            checkin_interval: 60,
            checkin_timeout: 10,
          },
          device: {
            hardware_id,
            key: {
              key: challengeKey,
              key_type: 'SOFTWARE_ECDSA_SHA256', // default key/challenge type
            },
          },
          connected_enodeb_serials: [],
          tier: 'default',
        },
      });
      const gateway = await MagmaV1API.getLteByNetworkIdGatewaysByGatewayId({
        networkId,
        gatewayId: gatewayID,
      });
      this.props.onSave(gateway);
    } catch (e) {
      this.setState({error: e?.response?.data?.message || e?.message || e});
    }
  };
}

export default withStyles(styles)(withRouter(AddGatewayDialog));
