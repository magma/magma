/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {ContextRouter} from 'react-router-dom';
import type {Gateway} from './GatewayUtils';
import type {WithStyles} from '@material-ui/core';

import React from 'react';
import axios from 'axios';
import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormField from './FormField';
import {GatewayStatus} from './GatewayUtils';
import Input from '@material-ui/core/Input';
import {MagmaAPIUrls} from '../common/MagmaAPI';
import moment from 'moment';

import {merge} from 'lodash-es';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = ContextRouter &
  WithAlert &
  WithStyles & {
    onClose: () => void,
    onSave: (gatewayID: string) => void,
    gateway: Gateway,
  };

type State = {
  name: string,
};

const styles = {
  input: {
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
};

class GatewaySummaryFields extends React.Component<Props, State> {
  state = {
    name: this.props.gateway.name,
  };

  render() {
    const {gateway} = this.props;
    return (
      <>
        <DialogContent>
          <FormField label="Name">
            <Input
              className={this.props.classes.input}
              value={this.state.name}
              onChange={this.nameChanged}
              placeholder="E.g. Gateway 1234"
            />
          </FormField>
          <FormField label="Hardware UUID">{gateway.hwid}</FormField>
          <FormField label="Gateway ID">{gateway.logicalID}</FormField>
          <FormField label="Last Checkin">
            {moment(parseInt(gateway.lastCheckin, 10)).fromNow()}
          </FormField>
          <FormField label="Version">{gateway.version}</FormField>
          <FormField label="VPN IP">{gateway.vpnIP}</FormField>
          <FormField label="RF Transmitter">
            <GatewayStatus
              isGrey={false}
              isActive={gateway.enodebRFTXEnabled}
            />
            {gateway.enodebRFTXEnabled ? '' : 'Not '}
            Allowed
            {'  '}
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.enodebRFTXOn}
            />
            {gateway.enodebRFTXOn ? '' : 'Not '}
            Connected
          </FormField>
          <FormField label="GPS synchronized">
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.gpsConnected}
            />
            {gateway.gpsConnected ? '' : 'Not '}
            Synced
          </FormField>
          <FormField label="MME">
            <GatewayStatus
              isGrey={gateway.isBackhaulDown}
              isActive={gateway.enodebConnected && gateway.mmeConnected}
            />
            {gateway.mmeConnected ? '' : 'Not '}
            Connected
          </FormField>
        </DialogContent>
        <DialogActions>
          <Button onClick={this.props.onClose} color="primary">
            Cancel
          </Button>
          <Button onClick={this.onSave} color="primary" variant="contained">
            Save
          </Button>
        </DialogActions>
      </>
    );
  }

  onSave = () => {
    const {match, gateway} = this.props;
    const id = gateway.logicalID;
    const url = MagmaAPIUrls.gateway(match, id);

    axios
      .get(url)
      .then(resp => {
        const data = merge(resp.data, {
          name: this.state.name,
        });
        axios
          .put(url, data)
          .then(() => this.props.onSave(id))
          .catch(error => this.props.alert(error.response.data.message));
      })
      .catch(error => this.props.alert(error.response.data.message));
  };

  nameChanged = ({target}: SyntheticInputEvent<*>) =>
    this.setState({name: target.value});
}

export default withStyles(styles)(withRouter(withAlert(GatewaySummaryFields)));
