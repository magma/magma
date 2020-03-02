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
import type {apn_list, subscriber} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
};

type EditingSubscriber = {
  imsiID: string,
  lteState: 'ACTIVE' | 'INACTIVE',
  authKey: string,
  authOpc: string,
  subProfile: string,
  apnList: apn_list,
};

type Props = ContextRouter &
  WithStyles<typeof styles> & {
    onClose: () => void,
    onSave: (subscriberID: string) => void,
    onSaveError: (reason: any) => void,
    editingSubscriber?: subscriber,
    subProfiles: Array<string>,
    apns: apn_list,
  };

type State = {
  error: string,
  editingSubscriber: EditingSubscriber,
};

class AddEditSubscriberDialog extends React.Component<Props, State> {
  state = {
    error: '',
    editingSubscriber: this.getEditingSubscriber(),
  };

  getEditingSubscriber(): EditingSubscriber {
    const {editingSubscriber} = this.props;
    if (!editingSubscriber) {
      return {
        imsiID: '',
        lteState: 'ACTIVE',
        authKey: '',
        authOpc: '',
        subProfile: 'default',
        apnList: [],
      };
    }

    const authKey = editingSubscriber.lte.auth_key
      ? base64ToHex(editingSubscriber.lte.auth_key)
      : '';

    const authOpc = editingSubscriber.lte.auth_opc
      ? base64ToHex(editingSubscriber.lte.auth_opc)
      : '';

    return {
      imsiID: editingSubscriber.id,
      lteState: editingSubscriber.lte.state,
      authKey,
      authOpc,
      subProfile: editingSubscriber.lte.sub_profile,
      apnList: editingSubscriber.active_apns || [],
    };
  }

  render() {
    const {classes} = this.props;
    const error = this.state.error ? (
      <FormLabel error>{this.state.error}</FormLabel>
    ) : null;

    return (
      <Dialog open={true} onClose={this.props.onClose}>
        <DialogTitle>
          {this.props.editingSubscriber ? 'Edit Subscriber' : 'Add Subscriber'}
        </DialogTitle>
        <DialogContent>
          {error}
          <TextField
            label="IMSI"
            className={classes.input}
            disabled={!!this.props.editingSubscriber}
            value={this.state.editingSubscriber.imsiID}
            onChange={this.imsiChanged}
          />
          <FormControl className={classes.input}>
            <InputLabel htmlFor="lteState">LTE Subscription State</InputLabel>
            <Select
              inputProps={{id: 'lteState'}}
              value={this.state.editingSubscriber.lteState}
              onChange={this.lteStateChanged}>
              <MenuItem value="ACTIVE">Active</MenuItem>
              <MenuItem value="INACTIVE">Inactive</MenuItem>
            </Select>
          </FormControl>
          <TextField
            label="LTE Auth Key"
            className={classes.input}
            value={this.state.editingSubscriber.authKey}
            onChange={this.authKeyChanged}
          />
          <TextField
            label="LTE Auth OPc"
            className={classes.input}
            value={this.state.editingSubscriber.authOpc}
            onChange={this.authOpcChanged}
          />
          <FormControl className={classes.input}>
            <InputLabel htmlFor="subProfile">Data Plan</InputLabel>
            <Select
              inputProps={{id: 'subProfile'}}
              value={this.state.editingSubscriber.subProfile}
              onChange={this.subProfileChanged}>
              {this.props.subProfiles.map(p => (
                <MenuItem value={p} key={p}>
                  {p}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          <FormControl className={classes.input}>
            <InputLabel htmlFor="apnList">Access Point Names</InputLabel>
            <Select
              inputProps={{id: 'apnList'}}
              value={this.state.editingSubscriber.apnList}
              multiple={true}
              onChange={this.apnListChanged}>
              {this.props.apns.map(apn => (
                <MenuItem value={apn} key={apn}>
                  {apn}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} skin="regular">
            Cancel
          </Button>
          <Button onClick={this.onSave}>Save</Button>
        </DialogActions>
      </Dialog>
    );
  }

  onSave = () => {
    if (
      !this.state.editingSubscriber.imsiID ||
      !this.state.editingSubscriber.authKey
    ) {
      this.setState({error: 'Please complete all fields'});
      return;
    }

    let {imsiID} = this.state.editingSubscriber;
    if (!imsiID.startsWith('IMSI')) {
      imsiID = `IMSI${imsiID}`;
    }

    const data = {
      id: imsiID,
      lte: {
        state: this.state.editingSubscriber.lteState,
        auth_algo: 'MILENAGE', // default auth algo
        auth_key: this.state.editingSubscriber.authKey,
        auth_opc: this.state.editingSubscriber.authOpc || undefined,
        sub_profile: this.state.editingSubscriber.subProfile,
      },
      active_apns: this.state.editingSubscriber.apnList,
    };
    if (data.lte.auth_key && isValidHex(data.lte.auth_key)) {
      data.lte.auth_key = hexToBase64(data.lte.auth_key);
    }
    if (data.lte.auth_opc && isValidHex(data.lte.auth_opc)) {
      data.lte.auth_opc = hexToBase64(data.lte.auth_opc);
    }
    if (this.props.editingSubscriber) {
      MagmaV1API.putLteByNetworkIdSubscribersBySubscriberId({
        networkId: nullthrows(this.props.match.params.networkId),
        subscriberId: data.id,
        subscriber: data,
      })
        .then(() => this.props.onSave(data.id))
        .catch(this.props.onSaveError);
    } else {
      MagmaV1API.postLteByNetworkIdSubscribers({
        networkId: this.props.match.params.networkId || '',
        subscriber: data,
      })
        .then(() => this.props.onSave(data.id))
        .catch(this.props.onSaveError);
    }
  };

  fieldChangedHandler = (
    field: 'imsiID' | 'authKey' | 'authOpc' | 'subProfile' | 'apnList',
  ) => event =>
    this.setState({
      editingSubscriber: {
        ...this.state.editingSubscriber,
        // $FlowFixMe Set state for each field
        [field]: event.target.value,
      },
    });

  imsiChanged = this.fieldChangedHandler('imsiID');
  lteStateChanged = event =>
    this.setState({
      editingSubscriber: {
        ...this.state.editingSubscriber,
        lteState: event.target.value === 'ACTIVE' ? 'ACTIVE' : 'INACTIVE',
      },
    });
  authKeyChanged = this.fieldChangedHandler('authKey');
  authOpcChanged = this.fieldChangedHandler('authOpc');
  subProfileChanged = this.fieldChangedHandler('subProfile');
  apnListChanged = this.fieldChangedHandler('apnList');
}

export default withStyles(styles)(withRouter(AddEditSubscriberDialog));
